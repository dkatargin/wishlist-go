package account

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"
)

// mockAccountRepo — ручной мок domain.AccountRepository для unit-тестов usecase.
type mockAccountRepo struct {
	getFn    func(id int64) (*domain.Account, error)
	createFn func(ctx context.Context, telegramID *int64) (int64, error)
	deleteFn func(id int64) error
	created  []int64
}

func (m *mockAccountRepo) CreateAccount(ctx context.Context, telegramID *int64) (int64, error) {
	if m.createFn != nil {
		return m.createFn(ctx, telegramID)
	}
	m.created = append(m.created, *telegramID)
	return *telegramID, nil
}

func (m *mockAccountRepo) GetAccountByID(id int64) (*domain.Account, error) {
	return m.getFn(id)
}

func (m *mockAccountRepo) DeleteAccount(id int64) error {
	if m.deleteFn != nil {
		return m.deleteFn(id)
	}
	return nil
}

func TestEnsureExists_CreatesAccountWhenMissing(t *testing.T) {
	repo := &mockAccountRepo{
		getFn: func(id int64) (*domain.Account, error) {
			return nil, domain.ErrAccountNotFound
		},
	}
	svc := NewService(repo)

	if err := svc.EnsureExists(context.Background(), 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.created) != 1 || repo.created[0] != 42 {
		t.Fatalf("expected account 42 to be created, got created=%v", repo.created)
	}
}

func TestEnsureExists_DoesNotCreateWhenPresent(t *testing.T) {
	repo := &mockAccountRepo{
		getFn: func(id int64) (*domain.Account, error) {
			return &domain.Account{ID: id}, nil
		},
		createFn: func(ctx context.Context, telegramID *int64) (int64, error) {
			t.Fatalf("CreateAccount must not be called when the account already exists")
			return 0, nil
		},
	}
	svc := NewService(repo)

	if err := svc.EnsureExists(context.Background(), 7); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureExists_PropagatesUnexpectedLookupError(t *testing.T) {
	boom := errors.New("db is down")
	repo := &mockAccountRepo{
		getFn: func(id int64) (*domain.Account, error) {
			return nil, boom
		},
		createFn: func(ctx context.Context, telegramID *int64) (int64, error) {
			t.Fatalf("CreateAccount must not be called when lookup fails unexpectedly")
			return 0, nil
		},
	}
	svc := NewService(repo)

	if err := svc.EnsureExists(context.Background(), 99); !errors.Is(err, boom) {
		t.Fatalf("expected lookup error to propagate, got %v", err)
	}
}
