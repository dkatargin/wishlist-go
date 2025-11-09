package account

import (
	"context"
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

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repo.DeleteAccount(id)
}
