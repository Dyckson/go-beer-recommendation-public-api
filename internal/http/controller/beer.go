package controller

import (
	"backend-test/internal/domain"
	"backend-test/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BeerController struct {
	BeerService       service.BeerServiceInterface
	ValidationService service.ValidationServiceInterface
	UpdateService     service.UpdateServiceInterface
}

func NewBeerController(beerService service.BeerServiceInterface, validationService service.ValidationServiceInterface, updateService service.UpdateServiceInterface) *BeerController {
	return &BeerController{
		BeerService:       beerService,
		ValidationService: validationService,
		UpdateService:     updateService,
	}
}

func (bc *BeerController) ListAllBeerStyles(c *gin.Context) {
	beerStyles, err := bc.BeerService.ListAllBeerStyles()
	if err != nil {
		log.Printf("controller=BeerController func=ListAllBeerStyles err=%v", err)

		status := http.StatusInternalServerError
		message := "internal error"

		if bc.ValidationService.IsNoRowsError(err) {
			status = http.StatusNotFound
			message = "no beer styles found"
		}

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"beerStyles": beerStyles,
	})
}

func (bc *BeerController) CreateBeerStyle(c *gin.Context) {
	var rawData map[string]interface{}
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "cannot read request body",
		})
		return
	}

	if err := json.Unmarshal(body, &rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid JSON format",
		})
		return
	}

	if _, exists := rawData["name"]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "name is required",
		})
		return
	}

	if _, exists := rawData["temp_min"]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "temp_min is required",
		})
		return
	}

	if _, exists := rawData["temp_max"]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "temp_max is required",
		})
		return
	}

	var inputStyle domain.BeerStyle
	if err := json.Unmarshal(body, &inputStyle); err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid field types",
		})
		return
	}

	if err := bc.ValidationService.ValidateUniqueNameForCreate(inputStyle.Name); err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)

		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to validate beer style name",
		})
		return
	}

	if err := bc.ValidationService.ValidateTemperatureRange(inputStyle); err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	newBeerStyle, err := bc.BeerService.CreateBeerStyle(inputStyle)
	if err != nil {
		log.Printf("controller=BeerController func=CreateBeerStyle name=%s err=%v", inputStyle.Name, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create beer style",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": newBeerStyle,
	})
}

func (bc *BeerController) UpdateBeerStyle(c *gin.Context) {
	beerUUID := c.Param("beerUUID")
	if beerUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "beerUUID is required",
		})
		return
	}

	var updateRequest domain.BeerStyleUpdateRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	currentBeerStyle, err := bc.BeerService.GetBeerStyleByUUID(beerUUID)
	if err != nil {
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
		status := http.StatusInternalServerError
		message := "internal error"

		if bc.ValidationService.IsNoRowsError(err) {
			status = http.StatusNotFound
			message = "beer style not found"
		}

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	if updateRequest.Name != nil && *updateRequest.Name != "" && *updateRequest.Name != currentBeerStyle.Name {
		if err := bc.ValidationService.ValidateUniqueNameForUpdate(*updateRequest.Name, currentBeerStyle.UUID); err != nil {
			log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)

			if strings.Contains(err.Error(), "already exists") {
				c.JSON(http.StatusConflict, gin.H{
					"message": err.Error(),
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to validate beer style name",
			})
			return
		}
	}

	changed := bc.UpdateService.ApplyBeerStyleUpdates(&currentBeerStyle, updateRequest)

	if changed {
		if err := bc.ValidationService.ValidateTemperatureRange(currentBeerStyle); err != nil {
			log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
	}

	if !changed {
		c.JSON(http.StatusOK, gin.H{
			"message": "No changes detected.",
			"data":    currentBeerStyle,
		})
		return
	}

	updatedBeerStyle, err := bc.BeerService.UpdateBeerStyle(currentBeerStyle)
	if err != nil {
		log.Printf("controller=BeerController func=UpdateBeerStyle beerUUID=%s err=%v", beerUUID, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "failed to update beer style",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beer style updated.",
		"data":    updatedBeerStyle,
	})
}

func (bc *BeerController) DeleteBeerStyle(c *gin.Context) {
	beerUUID := c.Param("beerUUID")
	if beerUUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "beerUUID is required",
		})
		return
	}

	_, err := bc.BeerService.GetBeerStyleByUUID(beerUUID)
	if err != nil {
		log.Printf("controller=BeerController func=DeleteBeerStyle beerUUID=%s err=%v", beerUUID, err)
		status := http.StatusInternalServerError
		message := "internal error"

		if bc.ValidationService.IsNoRowsError(err) {
			status = http.StatusNotFound
			message = "beer style not found"
		}

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	err = bc.BeerService.DeleteBeerStyle(beerUUID)
	if err != nil {
		log.Printf("controller=BeerController func=DeleteBeerStyle beerUUID=%s err=%v", beerUUID, err)
		status := http.StatusInternalServerError
		message := "internal error"

		c.AbortWithStatusJSON(status, gin.H{
			"message": message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beer style deleted successfully",
	})
}
