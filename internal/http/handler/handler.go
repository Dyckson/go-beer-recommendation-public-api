package handler

import (
	"backend-test/internal/http/controller"
	"backend-test/internal/service"
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
}
