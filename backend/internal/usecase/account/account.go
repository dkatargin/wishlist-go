package account

import (
	"context"
	"errors"
	"wishlist-go/internal/domain"
)

type Service struct {
	repo domain.AccountRepository
}

func NewService(repo domain.AccountRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// EnsureExists гарантирует наличие аккаунта с данным telegramID: находит его
// или лениво создаёт. Используется в Telegram auth-middleware.
func (s *Service) EnsureExists(ctx context.Context, telegramID int64) error {
	if _, err := s.repo.GetAccountByID(telegramID); err == nil {
		return nil
	} else if !errors.Is(err, domain.ErrAccountNotFound) {
		return err
	}

	id := telegramID
	if _, err := s.repo.CreateAccount(ctx, &id); err != nil {
		return err
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.DeleteAccount(id)
}
