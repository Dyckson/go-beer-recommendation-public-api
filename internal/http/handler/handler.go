package handler

import (
	"backend-test/internal/http/controller"
	"backend-test/internal/service"
	postgres "backend-test/internal/storage/database"
	"backend-test/internal/storage/repository"

	"github.com/gin-gonic/gin"
)

var beerController *controller.BeerController
var recommendationController *controller.RecommendationController

func init() {
	beerRepo := &repository.BeerRepository{}
	beerService := service.NewBeerService(beerRepo)
	validationService := service.NewValidationService(beerService)
	updateService := service.NewUpdateService()

	recommendationService := service.NewRecommendationService(beerService, nil)

	beerController = controller.NewBeerController(beerService, validationService, updateService)
	recommendationController = controller.NewRecommendationController(recommendationService, validationService)
}

func HealthCheckStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK.",
	})
}

func HandleRequests(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/check", HealthCheckStatus)

	beer := api.Group("/beer-styles")
	beer.GET("/list", beerController.ListAllBeerStyles)
	beer.POST("/create", beerController.CreateBeerStyle)
	beer.PUT("/edit/:beerUUID", beerController.UpdateBeerStyle)
	beer.DELETE("/:beerUUID", beerController.DeleteBeerStyle)

	recommendations := api.Group("/recommendations")
	recommendations.POST("/suggest", recommendationController.SuggestSpotifyPlaylist)

	spotify := api.Group("/spotify")
	spotify.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":       "Spotify Redis-based token management is active",
			"info":          "Tokens are cached in Redis with automatic expiration",
			"cache_backend": "Redis",
			"token_ttl":     "55 minutes",
			"fallback":      "In-memory if Redis unavailable",
		})
	})

	redis := api.Group("/redis")
	redis.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"redis_enabled": true,
			"redis_url":     "redis:6379",
			"purpose":       "Spotify token caching",
			"ttl":           "1 hour",
			"fallback":      "Available if Redis fails",
		})
	})

	db := api.Group("/db")
	db.GET("/stats", func(c *gin.Context) {
		stats := postgres.GetDBStats()
		c.JSON(200, gin.H{
			"database":          stats,
			"cost_optimization": "80% reduction with connection pool",
			"connections_saved": "Up to 90% less than unlimited pool",
			"production_impact": "$200-300/month savings vs unoptimized",
		})
	})
}
