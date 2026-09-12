package service

import (
	"ride-sharing/services/trip-service/internal/domain"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
)

type TripService struct {
	tripRepo     repository.TripRepository
	rideFareRepo repository.RideFareRepository
}

func NewTripService(tripRepo repository.TripRepository, rideFareRepo repository.RideFareRepository) *TripService {
	return &TripService{
		tripRepo:     tripRepo,
		rideFareRepo: rideFareRepo,
	}
}

func (s *TripService) Create(t *domain.Trip) error {
	rf, err := s.rideFareRepo.Create(&t.RideFare)
	if err != nil {
		return err
	}

	t.RideFare = *rf

	return s.tripRepo.Create(t)
}
