package payment

import (
	"context"

	"github.com/go-pg/pg/v10"
)

type IPaymentCodeRepository interface {
	Create(ctx context.Context, payment Payment) error
	Get(ctx context.Context, id string) (*Payment, error)
}

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) Create(ctx context.Context, payment *Payment) error {
	_, err := r.db.Model(payment).Insert()
	return err
}

func (r Repository) Get(ctx context.Context, id string) (*Payment, error) {
	payment := new(Payment)
	err := r.db.Model(payment).Where("id = ?", id).Select()
	return payment, err
}
