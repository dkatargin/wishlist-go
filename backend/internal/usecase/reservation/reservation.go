// Package reservation реализует бизнес-правила резервирования подарков:
// кто может зарезервировать, в каком количестве и что видно гостю списка.
package reservation

import (
	"context"
	"errors"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
)

// Service инкапсулирует правила доступа к резервам.
type Service struct {
	resRepo  domain.ReservationRepository
	wishRepo domain.WishItemRepository
}

func NewService(resRepo domain.ReservationRepository, wishRepo domain.WishItemRepository) *Service {
	return &Service{resRepo: resRepo, wishRepo: wishRepo}
}

// WishReservationView — агрегат резервов по одному желанию для гостевого вида.
// Анонимные дарители учитываются в сумме, но не попадают в Reservers.
type WishReservationView struct {
	ReservedQuantity int
	Reservers        []int64
}

// Reserve резервирует quantity единиц желания wishID в списке shareCode от имени reserverID.
// Правила: желание существует; нельзя резервировать на своём списке; quantity ≥ 1;
// суммарный резерв не превышает market_quantity (если он задан); один даритель — один резерв на желание.
func (s *Service) Reserve(ctx context.Context, reserverID int64, shareCode uuid.UUID, wishID int64, quantity int, isAnonymous bool) (*domain.Reservation, error) {
	if quantity < 1 {
		quantity = 1
	}

	wish, err := s.wishRepo.GetWishItemByID(wishID, shareCode)
	if err != nil {
		return nil, err
	}
	if wish.OwnerID == reserverID {
		return nil, domain.ErrSelfReservation
	}

	// один даритель не может зарезервировать одно желание дважды
	if _, err := s.resRepo.GetByReserverAndWish(ctx, reserverID, wishID); err == nil {
		return nil, domain.ErrReservationExists
	} else if !errors.Is(err, domain.ErrReservationNotFound) {
		return nil, err
	}

	// Лимит по количеству действует, только если у желания задан market_quantity (> 0).
	// Известное ограничение: проверка суммы и вставка не атомарны — при одновременном
	// резерве двумя разными дарителями сумма может превысить market_quantity. Вред мал
	// (мягкий лимит), на текущем масштабе принято; при необходимости — транзакция с
	// блокировкой строки желания.
	if wish.MarketQuantity != nil && *wish.MarketQuantity > 0 {
		sum, err := s.resRepo.SumQuantityByWish(ctx, wishID)
		if err != nil {
			return nil, err
		}
		if sum+quantity > *wish.MarketQuantity {
			return nil, domain.ErrReserveQuantityExceeded
		}
	}

	res := &domain.Reservation{
		WishID:      wishID,
		ReserverID:  reserverID,
		Quantity:    quantity,
		IsAnonymous: isAnonymous,
	}
	if err := s.resRepo.Create(ctx, res); err != nil {
		return nil, err
	}
	return res, nil
}

// CancelByWish снимает резерв дарителя на конкретном желании.
// Резерв уникален по (даритель, желание), поэтому ищем его и удаляем со скоупом по дарителю.
func (s *Service) CancelByWish(ctx context.Context, reserverID int64, wishID int64) error {
	res, err := s.resRepo.GetByReserverAndWish(ctx, reserverID, wishID)
	if err != nil {
		return err // ErrReservationNotFound пробрасывается как есть
	}
	return s.resRepo.Delete(ctx, res.ID, reserverID)
}

// MarkPurchased помечает резерв купленным; скоуп по дарителю на стороне репозитория.
func (s *Service) MarkPurchased(ctx context.Context, reserverID int64, reservationID uuid.UUID) error {
	return s.resRepo.MarkDone(ctx, reservationID, reserverID)
}

// ListMine возвращает резервы, сделанные дарителем reserverID.
func (s *Service) ListMine(ctx context.Context, reserverID int64) ([]*domain.Reservation, error) {
	return s.resRepo.ListByReserver(ctx, reserverID)
}

// AggregateForSharedList агрегирует резервы по всем желаниям списка shareCode.
// Для каждого желания: суммарное зарезервированное количество и id не-анонимных дарителей.
func (s *Service) AggregateForSharedList(ctx context.Context, shareCode uuid.UUID) (map[int64]WishReservationView, error) {
	all, err := s.resRepo.ListByWishlist(ctx, shareCode)
	if err != nil {
		return nil, err
	}
	out := make(map[int64]WishReservationView, len(all))
	for _, r := range all {
		v := out[r.WishID]
		v.ReservedQuantity += r.Quantity
		if !r.IsAnonymous {
			v.Reservers = append(v.Reservers, r.ReserverID)
		}
		out[r.WishID] = v
	}
	return out, nil
}
