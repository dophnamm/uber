package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ride-sharing/services/trip-service/internal/domain"
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
	var req domain.Trip

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Create(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, req)
}
