package payment

import "time"

type Payment struct {
	Id             string    `json:"id" pg:"id"`
	PaymentCode    string    `json:"payment_code" pg:"payment_code" validate:"required"`
	Name           string    `json:"name" pg:"name" validate:"required"`
	Status         string    `json:"status" pg:"status"`
	ExpirationDate time.Time `json:"expiration_date" pg:"expiration_date"`
	CreatedAt      time.Time `json:"created_at" pg:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" pg:"updated_at"`
}

const (
	PAYMENT_CODE_STATUS_ACTIVE   = "ACTIVE"
	PAYMENT_CODE_STATUS_EXPIRED  = "EXPIRED"
	PAYMENT_CODE_STATUS_INACTIVE = "INACTIVE"
)
