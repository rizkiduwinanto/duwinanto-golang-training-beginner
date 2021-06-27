package payment

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IService interface {
	Create(paymentCode *Payment) error
	Get(id string) (paymentCode *Payment, err error)
}

type Service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return Service{
		repo: r,
	}
}

func (s Service) Create(ctx context.Context, paymentCode *Payment) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	paymentCode.Id = id.String()

	now := time.Now().UTC()
	paymentCode.CreatedAt = now
	paymentCode.UpdatedAt = now

	expDate := now.AddDate(51, 0, 0)
	paymentCode.ExpirationDate = expDate

	paymentCode.Status = PAYMENT_CODE_STATUS_ACTIVE

	err = s.repo.Create(ctx, paymentCode)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Get(ctx context.Context, id string) (paymentCode *Payment, err error) {
	pc, err := s.repo.Get(ctx, id)
	if err != nil {
		return pc, err
	}

	return pc, nil
}
