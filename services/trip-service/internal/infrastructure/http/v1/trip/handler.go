package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ride-sharing/services/trip-service/internal/service"
)

type TripHandler struct {
	svc *service.TripService
}

func NewTripHandlerV1(svc *service.TripService) *TripHandler {
	return &TripHandler{
		svc: svc,
	}
}

func (h *TripHandler) Create(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{
		"messages": "ok",
	})
}
