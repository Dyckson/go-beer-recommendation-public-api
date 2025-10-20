package main

import (
	_ "backend-test/internal/cmd/server"
	"backend-test/internal/http/handler"
	"backend-test/internal/http/router"
	postgres "backend-test/internal/storage/database"
	"os"
)

func main() {
	defer postgres.CloseDB()

	r := router.NewRouter()
	handler.HandleRequests(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "1111"
	}

	r.Run(":" + port)
}
