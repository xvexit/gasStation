package usecase

import (
	"context"
	"fuelStation/internal/domain/entity"
	"fuelStation/internal/domain/interfaces"
	"fuelStation/internal/domain/service"
	"time"
)

type UseCase struct {
	refuelService *service.RefuelOperationService
	refuelRepo    interfaces.RefuelOperationsRepo

	priceService *service.FuelPriceService
	priceRepo    interfaces.FuelPriceRepository

	counterService *service.CounterStateService
	counterRepo    interfaces.CounterRepository
}

func NewUsecase(
	refuelService *service.RefuelOperationService,
	refuelRepo interfaces.RefuelOperationsRepo,
	priceService *service.FuelPriceService,
	priceRepo interfaces.FuelPriceRepository,
	counterService *service.CounterStateService,
	counterRepo interfaces.CounterRepository,

) *UseCase {
	return &UseCase{
		refuelService: refuelService,
		refuelRepo:    refuelRepo,

		priceService: priceService,
		priceRepo:    priceRepo,

		counterService: counterService,
		counterRepo:    counterRepo,
	}
}

// Создание заправки
func (u *UseCase) CreateRefuel(ctx context.Context, amountPaid float64, counterBeforeRefill int) (entity.RefuelOperation, error) {
	return u.refuelService.CreateRefuel(ctx, amountPaid, counterBeforeRefill)
}

// Подтверждение заправки
func (u *UseCase) ConfirmRefuel(ctx context.Context, id string) (entity.RefuelOperation, error) {
	return u.refuelService.ConfirmRefuel(ctx, id)
}

// Отмена заправки
func (u *UseCase) CancelRefuel(ctx context.Context, id, reason string) (entity.RefuelOperation, error) {
	return u.refuelService.CancelRefuel(ctx, id, reason)
}

// Проверка наличия незавершенных заправок
func (u *UseCase) HasPendingOperations(ctx context.Context) (bool, error) {
	return u.refuelService.HasPendingOperations(ctx)
}

// Показать незавершенные заправки
func (u *UseCase) GetPendingOperations(ctx context.Context) ([]entity.RefuelOperation, error) {
	return u.refuelService.GetPendingOperations(ctx)
}

// Получить историю заправок за период
func (u *UseCase) GetRefuelHistory(ctx context.Context, from, to time.Time, status string)([]entity.RefuelOperation, error){
	filter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo: &to,
		Status: &status,
	}	
	return u.refuelRepo.Find(ctx, filter)
}

// Получить заработанные деньги за период
func (u *UseCase) GetTotalRevenue(ctx context.Context, from, to time.Time) (float64, error) {
	return u.refuelService.GetTotalRevenue(ctx, from, to)
}

// Получить потраченные литры за период
func (u *UseCase) GetTotalLiters(ctx context.Context, from, to time.Time) (float64, error) {
	return u.refuelService.GetTotalLiters(ctx, from, to)
}

// Получить статистику за период
func (u *UseCase) GetStatistics(ctx context.Context, from, to time.Time) (service.RefuelStatistics, error) {
	return u.refuelService.GetStatistics(ctx, from, to)
}

// Получить операцио по id
func (u *UseCase) GetRefuelById(ctx context.Context, id string) (entity.RefuelOperation, error) {
	return u.refuelRepo.GetByID(ctx, id)
}

// Получить все заправки за период
func (u *UseCase) GetAllRefuel(ctx context.Context, from, to time.Time) ([]entity.RefuelOperation, error) {
	filter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo:   &to,
	}
	return u.refuelRepo.Find(ctx, filter)
}

// Получить активную цену
func (u *UseCase) GetPricePerLiter(ctx context.Context) (entity.FuelPrice, error) {
	return u.priceRepo.GetActive(ctx)
}

// Поменять цену за литр на новую
func (u *UseCase) ChangePricePerLiter(ctx context.Context, newPrice float64) (entity.FuelPrice, error) {
	return u.priceService.ChangePrice(ctx, newPrice)
}

// Установить цену за литр впервые
func (u *UseCase) InitPricePerLiter(ctx context.Context, newPrice float64) (entity.FuelPrice, error) {
	return u.priceService.InitPrice(ctx, newPrice)
}

// Найти изменения цен за промежуток
func (u *UseCase) GetPriceHistory(ctx context.Context, from, to time.Time) ([]entity.FuelPrice, error) {
	filter := interfaces.FuelPriceFilter{
		DateFrom: &from,
		DateTo:   &to,
	}
	return u.priceRepo.Find(ctx, filter)
}

// Получить текущее значение счетчика
func (u *UseCase) GetCurrent(ctx context.Context) (entity.CounterState, error) {
	return u.counterRepo.GetCurrent(ctx)
}

// Смена значения счетчика перед заправкой (с проверкой на большесть значения)
func (u *UseCase) UpdateCounterDuringRefuel(ctx context.Context, val int)(entity.CounterState, error){
	return u.counterService.UpdateCounterDuringRefuel(ctx, val)
}

// Смена значения счетчика без валидации
func (u *UseCase) UpdateCounter(ctx context.Context, val int)(entity.CounterState, error){
	return u.counterService.UpdateCounter(ctx, val)
}

