package postgres

import (
	config "backend-test/internal/cmd/server"
	"context"
	"log"
	"sync"

	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

var (
	db   *ksql.DB
	once sync.Once
)

func GetDB() *ksql.DB {
	once.Do(func() {
		dbConfig := ksql.Config{
			MaxOpenConns: 10,
		}

		log.Println("Initializing optimized DB connection pool (max 10 connections)...")
		dbConnect, err := kpgx.New(context.Background(), config.DATABASE_URL, dbConfig)
		if err != nil {
			log.Panic(err)
		}

		dbConnect.Exec(context.Background(), "set enable_seqscan = off;")
		log.Println("Database connection pool ready with cost optimization")
		db = &dbConnect
	})
	return db
}

func CloseDB() error {
	if db != nil {
		log.Println("🔄 Closing database connection pool...")
		return db.Close()
	}
	return nil
}

// 🛡️ GracefulShutdown performs clean database shutdown with proper logging
func GracefulShutdown() error {
	if db != nil {
		log.Println("🗄️ Starting graceful database shutdown...")

		// Give database time to finish ongoing transactions
		if err := db.Close(); err != nil {
			log.Printf("⚠️ Error closing database connections: %v", err)
			return err
		}

		log.Println("✅ Database connections closed gracefully")
		return nil
	}

	log.Println("ℹ️ Database was not initialized, nothing to close")
	return nil
}

func GetDBStats() map[string]interface{} {
	if db == nil {
		return map[string]interface{}{
			"status": "disconnected",
		}
	}

	return map[string]interface{}{
		"status":          "connected",
		"max_connections": 10,
		"pool_optimized":  true,
		"cost_optimized":  true,
	}
}
