package service

import (
	"context"
	"errors"
	"fuelStation/internal/domain/entity"
	"fuelStation/internal/domain/interfaces"
	"time"
)

type FuelPriceService struct {
	repo interfaces.FuelPriceRepository
}

func NewFuelPriceService(repo interfaces.FuelPriceRepository) *FuelPriceService {
	return &FuelPriceService{
		repo: repo,
	}
}

// Получить активную цену
func (f *FuelPriceService) GetActive(ctx context.Context)(entity.FuelPrice, error){
	return f.repo.GetActive(ctx)
}

// Уставновка цены впервые (аккуратно, может возникнуть две активные цены)
func (f *FuelPriceService) InitPrice(ctx context.Context, newPriceRub float64) (entity.FuelPrice, error) {

	if err := f.validatePrice(newPriceRub); err != nil {
		return entity.FuelPrice{}, err
	}

	_, err := f.repo.GetActive(ctx)
	if err == nil {
		return entity.FuelPrice{}, ErrTryUseChangePrice
	}

	if !errors.Is(err, ErrNotFoundOldPrice) {
		return entity.FuelPrice{}, err
	}

	fp := f.setPrice(newPriceRub)

	if err := f.repo.ChangePrice(ctx, fp); err != nil {
		return entity.FuelPrice{}, err
	}

	return fp, nil
}

// ChangePrice изменяет цену топлива, деактивируя предыдущую активную цену.
func (f *FuelPriceService) ChangePrice(ctx context.Context, newPriceRub float64) (entity.FuelPrice, error) {

	if err := f.validatePrice(newPriceRub); err != nil {
		return entity.FuelPrice{}, err
	}

	_, err := f.repo.GetActive(ctx)
	if err != nil {
		return entity.FuelPrice{}, err
	}

	newPrice := f.setPrice(newPriceRub) // добавление новой цены

	if err := f.repo.ActivateNewPrice(ctx, newPrice); err != nil {
		return entity.FuelPrice{}, err
	}

	return newPrice, nil
}

func (f *FuelPriceService) setPrice(newPriceRub float64) entity.FuelPrice {
	return entity.FuelPrice{
		ID:            entity.GenerateIDFuelPrice(),
		PricePerLiter: newPriceRub,
		CreatedAt:     time.Now(),
		IsActive:      true,
	}
}

func (f *FuelPriceService) validatePrice(price float64) error {
	if price <= 0 {
		return ErrPriceCanNotBeNegative
	}
	if price > 10000 {
		return ErrPriceTooHigh
	}
	return nil
}
