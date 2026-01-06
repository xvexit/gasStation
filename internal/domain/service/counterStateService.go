package service

import (
	"context"
	"fuelStation/internal/domain/entity"
	"fuelStation/internal/domain/interfaces"
	"time"
)

type CounterStateService struct {
	repo interfaces.CounterRepository
}

func NewCounterStateService(repo interfaces.CounterRepository) *CounterStateService {
	return &CounterStateService{
		repo: repo,
	}
}

// Обновление счетчика при заправке
func (s *CounterStateService) UpdateCounterDuringRefuel(ctx context.Context, newValue int) (entity.CounterState, error) {
	if err := s.validateUpdateCounter(ctx, newValue); err != nil {
		return entity.CounterState{}, err
	}

	updated, err := s.UpdateCounter(ctx, newValue)
	if err != nil {
		return entity.CounterState{}, err
	}

	return updated, nil
}

// Валидация при обновлении счетчика при заправке
func (s *CounterStateService) validateUpdateCounter(ctx context.Context, newValue int) error {

	current, err := s.repo.GetCurrent(ctx)
	if err != nil {
		return err
	}

	if current.CurrentValue > newValue {
		return ErrNewValueCanNotBeSmallerThanOld
	}

	return nil
}

// Обновление счетчика без валидации (пригодится если нужно будет выставить значение счетчика впервые либо после какого либо сбоя)
func (s *CounterStateService) UpdateCounter(ctx context.Context, newValue int) (entity.CounterState, error) {

	if newValue < 0 {
		return entity.CounterState{}, ErrCounterCanNotBeNegative
	}

	updated := entity.CounterState{
		Id: 0,
		CurrentValue: newValue,
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Save(ctx, updated); err != nil {
		return entity.CounterState{}, err
	}
	return updated, nil
}
