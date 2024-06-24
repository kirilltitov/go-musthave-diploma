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
	Order      string   `json:"order"`
	Status     string   `json:"status"`
	AccrualRaw *float64 `json:"accrual"`
}

func (r CalculationResult) Accrual() *decimal.Decimal {
	return decimalPtrFromFloat(r.AccrualRaw)
}

func decimalPtrFromFloat(f *float64) *decimal.Decimal {
	if f == nil {
		return nil
	}

	result := decimal.NewFromFloat(*f)

	return &result
}
