package accrual

import (
	"github.com/shopspring/decimal"

	"github.com/kirilltitov/go-musthave-diploma/internal/storage"
)

type Accrual interface {
	CalculateAmount(order storage.Order) (*CalculationResult, error)
}

const (
	StatusRegistered = "REGISTERED"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type CalculationResult struct {
	Order   string           `json:"order"`
	Status  string           `json:"status"`
	Accrual *decimal.Decimal `json:"accrual"`
}
