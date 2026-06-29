package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Favorite — сохранённый пользователем чужой (или свой) шаренный список.
type Favorite struct {
	ID           uuid.UUID
	AccountID    int64     // кто добавил
	WishlistCode uuid.UUID // ShareCode сохранённого списка
	CreatedAt    time.Time
}

type FavoriteRepository interface {
	// Add идемпотентен по unique(account, wishlist): повтор не ошибка.
	Add(ctx context.Context, f *Favorite) error
	Remove(ctx context.Context, accountID int64, wishlistCode uuid.UUID) error
	// ListByAccount возвращает сами списки (join к wishlist), а не голые id.
	ListByAccount(ctx context.Context, accountID int64) ([]*Wishlist, error)
}
