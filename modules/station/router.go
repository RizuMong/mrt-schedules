package station

import (
	"github.com/RizuMong/mrt-schedules/common/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Initiate(router *gin.RouterGroup) {
	stationService := NewService()

	station := router.Group("/stations")

	station.GET("", func(c *gin.Context) {
		GetAllStations(c, stationService)
	})
}

func GetAllStations(c *gin.Context, service Service) {
	datas, err := service.GetAllStations()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.APIResponse{
		Success: true,
		Message: "Stations retrieved successfully",
		Data:    datas,
	})
}
