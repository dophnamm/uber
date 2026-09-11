package routes

import (
	"github.com/gin-gonic/gin"

	tripV1 "ride-sharing/services/trip-service/internal/infrastructure/http/v1/trip"
)

func NewTripRoutesV1(r *gin.Engine, h *tripV1.TripHandler) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/trip", h.Create)
	}
}
