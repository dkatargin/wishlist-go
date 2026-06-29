package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Публичные/гостевые ответы не должны содержать идентичность владельца (OwnerID/account_id).
func TestPublicResponses_DoNotLeakOwner(t *testing.T) {
	name := "gift"
	shared := SharedWishlistResponse{
		ShareCode: uuid.New(),
		Name:      "list",
		Items:     []SharedWishItem{{ID: 1, Name: &name, MarketCurrency: "RUB"}},
	}
	detail := WishlistDetailResponse{ShareCode: uuid.New(), Name: "list", IsOwner: true}

	for _, v := range []any{shared, detail} {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		s := string(b)
		// account_id = OwnerID у Wishlist, owner_id = OwnerID у WishItem (is_owner — легитимный флаг)
		for _, leak := range []string{"account_id", "owner_id"} {
			if strings.Contains(s, leak) {
				t.Fatalf("ответ содержит %q (утечка владельца): %s", leak, s)
			}
		}
	}
}
