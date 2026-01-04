package service

import (
	"context"
	"fuelStation/internal/domain/entity"
	"fuelStation/internal/domain/interfaces"
	"time"
)

const (
	RefuelStatusCreated   = "Created"
	RefuelStatusConfirmed = "Confirmed"
	RefuelStatusCancelled = "Cancelled"
	LitersPerCounterUnit  = 10
)

type RefuelOperationService struct {
	refuelRepo     interfaces.RefuelOperationsRepo
	priceService   *FuelPriceService
	counterService *CounterStateService
}

func NewRefuelOperationService(refuelRepo interfaces.RefuelOperationsRepo, priceService *FuelPriceService, counterService *CounterStateService) *RefuelOperationService {
	return &RefuelOperationService{
		refuelRepo:     refuelRepo,
		priceService:   priceService,
		counterService: counterService,
	}
}

// Операция заправки
func (s *RefuelOperationService) CreateRefuel(ctx context.Context, amountPaid float64, counterBeforeRefill int) (entity.RefuelOperation, error) {

	//Валидация внесенных денег
	if err := s.validateAmount(amountPaid); err != nil {
		return entity.RefuelOperation{}, err
	}

	//Получение активной цены за литр
	priceObj, err := s.priceService.GetActive(ctx)
	if err != nil {
		return entity.RefuelOperation{}, err
	}

	//Получение значения счетчика
	currentCounter, err := s.counterService.repo.GetCurrent(ctx)
	if err != nil {
		return entity.RefuelOperation{}, err
	}

	// Проверка на идентичность введенного счетчика и имеющегося
	if currentCounter.CurrentValue != counterBeforeRefill {
		// Обновление счетчика введенным значением
		if _, err := s.counterService.UpdateCounter(ctx, counterBeforeRefill); err != nil {
			return entity.RefuelOperation{}, err
		}
	}

	// Подсчет литров
	liters := amountPaid / priceObj.PricePerLiter

	// Подсчет состояния счетчика после (счетчик содержит десятые части литра без точки)
	counterAfter := counterBeforeRefill + int(liters*LitersPerCounterUnit)

	operation := entity.RefuelOperation{
		ID:               entity.GenerateIDRefuelOperation(),
		AmountPaid:       amountPaid,
		CalculatedLiters: liters,
		PricePerLiter:    priceObj.PricePerLiter,
		CounterBefore:    counterBeforeRefill,
		CounterAfter:     counterAfter,
		Status:           RefuelStatusCreated,
		CreatedAt:        time.Now(),
	}

	// Создание новой записи
	if err := s.refuelRepo.Create(ctx, operation); err != nil {
		return entity.RefuelOperation{}, err
	}

	return operation, nil
}

func (s *RefuelOperationService) validateAmount(amount float64) error {

	if amount <= 0 {
		return ErrAmountCanNotBeNegative
	}

	if amount > 100000 {
		return ErrAmountTooHigh
	}

	return nil
}

// Подтверждает операцию и обновляет счётчик
func (s *RefuelOperationService) ConfirmRefuel(ctx context.Context, id string) (entity.RefuelOperation, error) {

	// Находит операцию по айди
	operation, err := s.refuelRepo.GetByID(ctx, id)
	if err != nil {
		return entity.RefuelOperation{}, err
	}

	// Проверка на соответствие статуса
	if operation.Status != RefuelStatusCreated {
		return entity.RefuelOperation{}, ErrInvalidOperationStatus
	}

	//  Проверяем что текущий счётчик не изменился
	current, err := s.counterService.repo.GetCurrent(ctx)
	if err != nil {
		return entity.RefuelOperation{}, err
	}

	// Если текущий != counterBefore (из создания), что-то не так
	if current.CurrentValue != operation.CounterBefore {
		return entity.RefuelOperation{}, ErrCounterWasChangedDuringRefuelCreation
	}

	// Обновление счетчика
	if _, err := s.counterService.UpdateCounterDuringRefuel(ctx, operation.CounterAfter); err != nil {
		return entity.RefuelOperation{}, err
	}

	// Обновление статуса
	if err := s.refuelRepo.UpdateStatus(ctx, id, RefuelStatusConfirmed, nil); err != nil {
		return entity.RefuelOperation{}, err
	}

	// Для вывода
	operation.Status = RefuelStatusConfirmed

	return operation, nil
}

// Отменяет операцию с указанием причины
func (s *RefuelOperationService) CancelRefuel(ctx context.Context, id string, reason string) (entity.RefuelOperation, error) {

	// Поиск операции
	operation, err := s.refuelRepo.GetByID(ctx, id)
	if err != nil {
		return entity.RefuelOperation{}, err
	}

	// Проверка на соответствие статуса
	if operation.Status != RefuelStatusConfirmed && operation.Status != RefuelStatusCreated {
		return entity.RefuelOperation{}, ErrInvalidOperationStatus
	}

	// Скрутить счетчик обратно
	if operation.Status == RefuelStatusConfirmed {
		counter, err := s.counterService.repo.GetCurrent(ctx)
		if err != nil {
			return entity.RefuelOperation{}, err
		}
		if _, errr := s.counterService.UpdateCounter(ctx, counter.CurrentValue-int(operation.CalculatedLiters*LitersPerCounterUnit)); errr != nil {
			return entity.RefuelOperation{}, errr
		}
	}

	// Обновление статуса
	if err := s.refuelRepo.UpdateStatus(ctx, id, RefuelStatusCancelled, &reason); err != nil {
		return entity.RefuelOperation{}, err
	}

	operation.Status = RefuelStatusCancelled
	p := time.Now()
	operation.CancelledAt = &p

	return operation, nil
}

// Проверка наличия незавершенных операций
func (s *RefuelOperationService) HasPendingOperations(ctx context.Context) (bool, error) {
	operations, err := s.GetPendingOperations(ctx)
	if err != nil {
		return false, err
	}

	return len(operations) > 0, nil
}

// Получение всех операций в статусе Created (незавершенных)
func (s *RefuelOperationService) GetPendingOperations(ctx context.Context) ([]entity.RefuelOperation, error) {

	status := RefuelStatusCreated
	filter := interfaces.RefuelFilter{
		Status: &status,
	}

	return s.refuelRepo.Find(ctx, filter)
}

// Получить заработанные деньги за период
func (s *RefuelOperationService) GetTotalRevenue(ctx context.Context, from, to time.Time) (float64, error) {
	var total float64
	status := RefuelStatusConfirmed

	filter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo:   &to,
		Status:   &status,
	}

	op, err := s.refuelRepo.Find(ctx, filter)
	if err != nil {
		return 0, err
	}

	for _, operation := range op {
		total += operation.AmountPaid
	}

	return total, nil
}

// Получение статистики за промежуток времени
func (s *RefuelOperationService) GetStatistics(
	ctx context.Context,
	from, to time.Time,
) (RefuelStatistics, error) {

	status := RefuelStatusConfirmed
	filter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo:   &to,
		Status:   &status,
	}

	op, err := s.refuelRepo.Find(ctx, filter)
	if err != nil {
		return RefuelStatistics{}, err
	}

	stats := RefuelStatistics{
		StartDate:       from,
		EndDate:         to,
		TotalOperations: 0,
		ConfirmedCount:  0,
		TotalRevenue:    0,
		TotalLiters:     0,
		AverageLiters:   0,
		AverageAmount:   0,
	}

	var (
		revenue, liters, averageLit, averageAmount float64
		confirmed                                  int64
	)

	for _, oper := range op {
		revenue += oper.AmountPaid
		liters += oper.CalculatedLiters
		confirmed++
	}

	if confirmed > 0 {
		averageLit = liters / float64(confirmed)
		averageAmount = revenue / float64(confirmed)
	}

	cStatus := RefuelStatusCancelled
	cFilter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo:   &to,
		Status:   &cStatus,
	}

	cancelledOp, err := s.refuelRepo.Find(ctx, cFilter)
	if err != nil {
		return RefuelStatistics{}, err
	}

	crStatus := RefuelStatusCreated
	crFilter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo:   &to,
		Status:   &crStatus,
	}

	createdOp, err := s.refuelRepo.Find(ctx, crFilter)
	if err != nil {
		return RefuelStatistics{}, err
	}

	stats.TotalOperations = confirmed + int64(len(cancelledOp)) + int64(len(createdOp))
	stats.TotalRevenue = revenue
	stats.TotalLiters = liters
	stats.AverageLiters = averageLit
	stats.AverageAmount = averageAmount
	stats.ConfirmedCount = confirmed
	stats.CancelledCount = int64(len(cancelledOp))
	stats.PendingCount = int64(len(createdOp))

	return stats, nil
}

// GetAveragePricePerLiter()

// Получить потраченные литры за промежуток
func (s *RefuelOperationService) GetTotalLiters(ctx context.Context, from, to time.Time) (float64, error) {
	var total float64
	status := RefuelStatusConfirmed

	filter := interfaces.RefuelFilter{
		DateFrom: &from,
		DateTo:   &to,
		Status:   &status,
	}

	operations, err := s.refuelRepo.Find(ctx, filter)
	if err != nil {
		return 0, err
	}

	for _, op := range operations {
		total += op.CalculatedLiters
	}

	return total, nil
}

type RefuelStatistics struct {
	TotalOperations int64     // всего операций
	TotalRevenue    float64   // рубли
	TotalLiters     float64   // литры
	AverageLiters   float64   // средний размер заправки
	AverageAmount   float64   // средняя сумма
	ConfirmedCount  int64     // подтверждённых
	CancelledCount  int64     // отменено
	PendingCount    int64     // ожидают подтверждения
	StartDate       time.Time // начало периода
	EndDate         time.Time // конец периода
}
