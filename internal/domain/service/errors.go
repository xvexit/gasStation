package service

import "errors"

var (
	ErrNotFoundOldPrice                      = errors.New("not found old price")
	ErrPriceCanNotBeNegative                 = errors.New("price can not be negative")
	ErrPriceTooHigh                          = errors.New("fuel price too high")
	ErrCounterCanNotBeNegative               = errors.New("counter can not be negative")
	ErrNewValueCanNotBeSmallerThanOld        = errors.New("new value can not be smaller than old value")
	ErrAmountCanNotBeNegative                = errors.New("amount paid can not be negative")
	ErrAmountTooHigh                         = errors.New("amount paid too high")
	ErrInvalidOperationStatus                = errors.New("this operation status must be CREATED or CONFIRMED")
	ErrCounterWasChangedDuringRefuelCreation = errors.New("counter was changed during refuel creation")
	ErrTryUseChangePrice                     = errors.New("you already have active price, try function change price")
	ErrNotFoundOper                          = errors.New("not found operations")
)
