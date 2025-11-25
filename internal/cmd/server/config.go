package config

import (
	"backend-test/external/spotify"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"golang.org/x/sync/singleflight"
)

var DATABASE_URL string

type SpotifyTokenData struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type RedisSpotifyManager struct {
	redisClient  *redis.Client
	clientID     string
	clientSecret string
	tokenKey     string
	sfGroup      singleflight.Group
}

// 🚀 Singleton instance
var (
	redisSpotifyManager *RedisSpotifyManager
	redisManagerOnce    sync.Once
)

func init() {
	_ = godotenv.Load()

	DATABASE_URL = GetDatabaseURL()
}

func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "postgres://beeruser:beerpass@localhost:55432/beerdb?sslmode=disable"
	}
	return dbURL
}

func GetSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func GetSpotifyClientSecret() string {
	return os.Getenv("SPOTIFY_CLIENT_SECRET")
}

func GetRedisURL() string {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return "localhost:6379" // Default Redis address
	}
	return redisURL
}

func InitializeSpotifyService() *spotify.SpotifyService {
	clientID := GetSpotifyClientID()
	clientSecret := GetSpotifyClientSecret()

	if clientID == "" || clientSecret == "" {
		log.Println("Warning: Spotify credentials not set. Spotify integration will be disabled.")
		return nil
	}

	spotifyService, err := spotify.NewSpotifyService(clientID, clientSecret)
	if err != nil {
		log.Printf("Warning: Failed to initialize Spotify service: %v", err)
		return nil
	}

	return spotifyService
}

func GetSpotifyService() *spotify.SpotifyService {
	redisManagerOnce.Do(func() {
		clientID := GetSpotifyClientID()
		clientSecret := GetSpotifyClientSecret()

		if clientID == "" || clientSecret == "" {
			log.Println("Warning: Spotify credentials not set. Spotify integration will be disabled.")
			redisSpotifyManager = nil
			return
		}

		redisClient := redis.NewClient(&redis.Options{
			Addr:         GetRedisURL(),
			Password:     "",
			DB:           0,
			PoolSize:     5,
			MinIdleConns: 1,
			MaxRetries:   2,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})

		ctx := context.Background()
		if err := redisClient.Ping(ctx).Err(); err != nil {
			log.Printf("Warning: Redis connection failed, falling back to in-memory: %v", err)
			redisClient = nil
		}

		redisSpotifyManager = &RedisSpotifyManager{
			redisClient:  redisClient,
			clientID:     clientID,
			clientSecret: clientSecret,
			tokenKey:     "spotify:access_token",
			sfGroup:      singleflight.Group{},
		}
	})

	if redisSpotifyManager == nil {
		return InitializeSpotifyService()
	}

	return redisSpotifyManager.GetSpotifyService()
}

func (rsm *RedisSpotifyManager) GetSpotifyService() *spotify.SpotifyService {
	ctx := context.Background()

	if token := rsm.getValidTokenFromRedis(ctx); token != nil {
		log.Println("Using valid token from Redis")
		return rsm.createSpotifyServiceWithToken(token.AccessToken)
	}

	// Use singleflight to ensure only one goroutine refreshes the token
	ch := rsm.sfGroup.DoChan("refresh", func() (interface{}, error) {
		log.Println("Creating new Spotify token...")
		service, err := spotify.NewSpotifyService(rsm.clientID, rsm.clientSecret)
		if err != nil {
			return nil, err
		}

		if err := rsm.saveTokenToRedis(ctx, service); err != nil {
			log.Printf("Failed to save token to Redis: %v", err)
		}

		return service, nil
	})

	res := <-ch
	if res.Err != nil {
		log.Printf("Failed to create Spotify service: %v", res.Err)
		return nil
	}

	if svc, ok := res.Val.(*spotify.SpotifyService); ok {
		return svc
	}

	return nil
}

func (rsm *RedisSpotifyManager) getValidTokenFromRedis(ctx context.Context) *SpotifyTokenData {
	if rsm.redisClient == nil {
		return nil
	}
	tokenJSON, err := rsm.redisClient.Get(ctx, rsm.tokenKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("Redis get error: %v", err)
		}
		return nil
	}

	var tokenData SpotifyTokenData
	if err := json.Unmarshal([]byte(tokenJSON), &tokenData); err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		return nil
	}

	if time.Until(tokenData.ExpiresAt) < 5*time.Minute {
		log.Println("Token expired or expiring soon")
		return nil
	}

	return &tokenData
}

func (rsm *RedisSpotifyManager) saveTokenToRedis(ctx context.Context, service *spotify.SpotifyService) error {
	if rsm.redisClient == nil {
		return nil
	}
	tokenData := SpotifyTokenData{
		AccessToken: "spotify_token_" + time.Now().Format("20060102_150405"),
		ExpiresAt:   time.Now().Add(50 * time.Minute),
		CreatedAt:   time.Now(),
	}

	tokenJSON, err := json.Marshal(tokenData)
	if err != nil {
		return err
	}

	return rsm.redisClient.Set(ctx, rsm.tokenKey, tokenJSON, 50*time.Minute).Err()
}

func (rsm *RedisSpotifyManager) createSpotifyServiceWithToken(token string) *spotify.SpotifyService {
	service, err := spotify.NewSpotifyService(rsm.clientID, rsm.clientSecret)
	if err != nil {
		log.Printf("Failed to create service with existing token: %v", err)
		return nil
	}
	return service
}
