package models

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	ID    uint            `json:"id"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price,omitempty"`
}
