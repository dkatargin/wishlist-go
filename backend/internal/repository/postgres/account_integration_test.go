package postgres

import (
	"context"
	"errors"
	"testing"
	"wishlist-go/internal/domain"
)

func TestAccountRepo_DeleteAccount_RemovesRow(t *testing.T) {
	db := newTestDB(t)
	repo := NewAccountRepository(db)
	ctx := context.Background()

	var tgID int64 = 12345
	if _, err := repo.CreateAccount(ctx, &tgID); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}

	if err := repo.DeleteAccount(tgID); err != nil {
		t.Fatalf("DeleteAccount: %v", err)
	}

	if _, err := repo.GetAccountByID(tgID); !errors.Is(err, domain.ErrAccountNotFound) {
		t.Fatalf("ожидался ErrAccountNotFound после удаления, получили %v", err)
	}
}
