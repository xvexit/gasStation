package interfaces

import (
	"context"
	"fuelStation/internal/domain/entity"
	"time"
)

type FuelPriceRepository interface {
	// Получить текущую активную цену
	GetActive(ctx context.Context) (entity.FuelPrice, error)

	// Создать/изменить цену (деактивирует предыдущую)
	ChangePrice(ctx context.Context, price entity.FuelPrice) error

	// Универсальный поиск цен
	Find(ctx context.Context, filter FuelPriceFilter) ([]entity.FuelPrice, error)

	// По ID (для редких случаев)
	GetByID(ctx context.Context, id string) (entity.FuelPrice, error)

	// Деактивировать все активные цены и активировать новую
	ActivateNewPrice(ctx context.Context, price entity.FuelPrice) error
}

type FuelPriceFilter struct {
	IsActive *bool
	DateFrom *time.Time
	DateTo   *time.Time
	Limit    *int
	Offset   *int
}

type CounterRepository interface {
	// Текущее состояние счётчика (всегда одна запись)
	GetCurrent(ctx context.Context) (entity.CounterState, error)

	// Сохранить новое состояние
	Save(ctx context.Context, state entity.CounterState) error
}

type RefuelOperationsRepo interface {
	Create(ctx context.Context, operation entity.RefuelOperation) error

	// Получить по ID
	GetByID(ctx context.Context, id string) (entity.RefuelOperation, error)

	// Универсальный поиск операций
	Find(ctx context.Context, filter RefuelFilter) ([]entity.RefuelOperation, error)

	// Обновить статус операции
	UpdateStatus(ctx context.Context, id string, status string, reason *string) error
}

type RefuelFilter struct {
	DeviceID *string
	DateFrom *time.Time
	DateTo   *time.Time
	Status   *string
	Limit    *int
	Offset   *int
}

type LogRepository interface {
	// Сохранить лог
	Create(ctx context.Context, log *entity.LogRecord) error

	// Универсальный поиск логов
	Find(ctx context.Context, filter LogFilter) ([]entity.LogRecord, error)

	// Удалить старые логи (для очистки)
	DeleteOld(ctx context.Context, before time.Time) error
}

type LogFilter struct {
	DeviceID  *string
	Level     *string // INFO, ERROR, WARNING
	EventType *string // RefuelCreated, PriceChanged
	DateFrom  *time.Time
	DateTo    *time.Time
	Limit     *int
	Offset    *int
}
