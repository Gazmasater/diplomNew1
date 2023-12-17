package systemservice

import (
	"context"

	"diplom.com/go-musthave-diploma-tpl/internal/dto"
	"diplom.com/go-musthave-diploma-tpl/internal/storage/postgres"
)

type service struct {
	db *postgres.Database
}

type Service interface {
	UserSetter(ctx context.Context, user, password string) error
	CredentialsGetter(ctx context.Context, user string) (string, string, error)
}

func NewService(db *postgres.Database) *service {
	service := service{
		db: db,
	}
	return &service
}

func (s *service) UserSetter(ctx context.Context, user, password, id string) error {
	err := s.db.UserSetter(ctx, user, password, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) CredentialsGetter(ctx context.Context, user string) (string, string, error) {
	hashPassword, id, err := s.db.CredentialsGetter(ctx, user)
	return hashPassword, id, err
}

func (s *service) OrderSetter(ctx context.Context, userID, orderNumber string) error {
	err := s.db.OrderSetter(ctx, userID, orderNumber)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetOrders(ctx context.Context, userID string) ([]dto.FullOrder, error) {
	orders, err := s.db.GetOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (s *service) GetUserBalance(ctx context.Context, userID string) (*dto.UserBalance, error) {
	balance, err := s.db.GetUserBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (s *service) ProcessOrder(ctx context.Context, order, userID string, sum float32) error {
	err := s.db.ProcessOrder(ctx, order, userID, sum)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetWithdrawals(ctx context.Context, userID string) ([]*dto.Billing, error) {
	balance, err := s.db.GetWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
