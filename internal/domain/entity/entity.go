package entity

import (
	"time"

	"github.com/google/uuid"
)

// Текущая цена топлива
type FuelPrice struct {
	ID            string    `json:"id"`
	PricePerLiter float64   `json:"price_per_liter"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	IsActive      bool      `json:"is_active"`
}

// Текущее состояние счётчика колонки/фургона
type CounterState struct {
	CurrentValue int       // текущее показание счётчика в литрах
	UpdatedAt    time.Time // время последнего обновления
}

// Операция заправки
type RefuelOperation struct {
	ID               string     // уникальный идентификатор операции
	AmountPaid       float64    // сумма денег, внесённая клиентом
	CalculatedLiters float64    // литры, рассчитанные по цене
	PricePerLiter    float64    // цена за литр на момент операции (копия)
	CounterBefore    int        // показание счётчика до заправки
	CounterAfter     int        // показание счётчика после заправки
	Status           string     // статус операции: Created, Confirmed, Cancelled и т.п.
	CreatedAt        time.Time  // дата и время создания операции
	CancelledAt      *time.Time // дата и время отмены (опционально)
}

// Логи для просмотра в приложении
type LogRecord struct {
	ID        string    // уникальный ID лога
	DeviceID  string    // ID устройства, с которого пришёл лог
	Level     string    // уровень: INFO, WARNING, ERROR и т.п.
	EventType string    // тип события, например, RefuelCreated, PriceChanged
	Message   string    // краткое описание события
	Meta      string    // дополнительные данные в JSON (может быть пустым)
	CreatedAt time.Time // время события
}

func GenerateIDFuelPrice() string {
	return uuid.New().String()
}

func GenerateIDRefuelOperation() string {
	return uuid.New().String()
}
