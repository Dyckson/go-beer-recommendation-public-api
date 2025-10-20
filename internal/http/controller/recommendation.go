package controller

import (
	"backend-test/internal/domain"
	"backend-test/internal/service"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RecommendationController struct {
	RecommendationService service.RecommendationServiceInterface
	ValidationService     service.ValidationServiceInterface
}

func NewRecommendationController(recommendationService service.RecommendationServiceInterface, validationService service.ValidationServiceInterface) *RecommendationController {
	return &RecommendationController{
		RecommendationService: recommendationService,
		ValidationService:     validationService,
	}
}

func (rc *RecommendationController) SuggestSpotifyPlaylist(c *gin.Context) {
	var request domain.TemperatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("controller=RecommendationController func=SuggestSpotifyPlaylist err=%v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body format",
		})
		return
	}

	if err := rc.ValidationService.ValidateTemperatureInput(request.Temperature); err != nil {
		log.Printf("controller=RecommendationController func=SuggestSpotifyPlaylist temperature=%.1f err=%v", request.Temperature, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	recommendation, err := rc.RecommendationService.GetRecommendationForTemperature(request.Temperature)
	if err != nil {
		log.Printf("controller=RecommendationController func=SuggestSpotifyPlaylist temperature=%.1f err=%v", request.Temperature, err)

		var status int
		var message string

		errorMsg := err.Error()
		switch {
		case strings.Contains(errorMsg, "no playlist found"):
			status = http.StatusNotFound
			message = errorMsg
		case strings.Contains(errorMsg, "spotify service unavailable"):
			status = http.StatusServiceUnavailable
			message = "Spotify service is temporarily unavailable"
		case strings.Contains(errorMsg, "failed to find best beer style"):
			status = http.StatusInternalServerError
			message = "Unable to determine suitable beer style"
		default:
			status = http.StatusInternalServerError
			message = "Internal server error"
		}

		c.JSON(status, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, recommendation)
}
