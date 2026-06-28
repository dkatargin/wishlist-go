package postgres

import (
	"context"
	"errors"
	"time"
	"wishlist-go/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reservationModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	WishID      int64     `gorm:"not null;uniqueIndex:uq_reservation_reserver_wish"`
	ReserverID  int64     `gorm:"not null;uniqueIndex:uq_reservation_reserver_wish"`
	Quantity    int       `gorm:"not null"`
	IsAnonymous bool      `gorm:"not null;default:false"`
	IsDone      bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// OnDelete:CASCADE — при удалении желания его резервы исчезают,
	// иначе FK RESTRICT не дал бы владельцу удалить зарезервированное желание.
	Wish wishItemModel `json:"-" gorm:"foreignKey:WishID;constraint:OnDelete:CASCADE"`
}

func (reservationModel) TableName() string { return "reservations" }

type reservationRepo struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) domain.ReservationRepository {
	return &reservationRepo{db: db}
}

func (r *reservationRepo) Create(ctx context.Context, res *domain.Reservation) error {
	m := reservationDomainToModel(res)
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		// Обычно дубль ловит usecase пред-проверкой; но в гонке INSERT упирается в
		// unique(reserver, wish) — транслируем в доменную ошибку, чтобы получить 409, а не 500.
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrReservationExists
		}
		return err
	}
	res.ID = m.ID
	res.CreatedAt = m.CreatedAt
	res.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *reservationRepo) Delete(ctx context.Context, id uuid.UUID, reserverID int64) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND reserver_id = ?", id, reserverID).
		Delete(&reservationModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrReservationNotFound
	}
	return nil
}

func (r *reservationRepo) MarkDone(ctx context.Context, id uuid.UUID, reserverID int64) error {
	result := r.db.WithContext(ctx).Model(&reservationModel{}).
		Where("id = ? AND reserver_id = ?", id, reserverID).
		Update("is_done", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrReservationNotFound
	}
	return nil
}

func (r *reservationRepo) GetByReserverAndWish(ctx context.Context, reserverID, wishID int64) (*domain.Reservation, error) {
	var m reservationModel
	err := r.db.WithContext(ctx).
		Where("reserver_id = ? AND wish_id = ?", reserverID, wishID).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrReservationNotFound
		}
		return nil, err
	}
	return reservationModelToDomain(&m), nil
}

func (r *reservationRepo) SumQuantityByWish(ctx context.Context, wishID int64) (int, error) {
	var sum int64
	err := r.db.WithContext(ctx).Model(&reservationModel{}).
		Where("wish_id = ?", wishID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&sum).Error
	return int(sum), err
}

func (r *reservationRepo) ListByWishlist(ctx context.Context, shareCode uuid.UUID) ([]*domain.Reservation, error) {
	var models []reservationModel
	err := r.db.WithContext(ctx).
		Joins("JOIN wish_items ON wish_items.id = reservations.wish_id").
		Where("wish_items.wish_list_code = ?", shareCode.String()).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return reservationsToDomain(models), nil
}

func (r *reservationRepo) ListByReserver(ctx context.Context, reserverID int64) ([]*domain.Reservation, error) {
	var models []reservationModel
	err := r.db.WithContext(ctx).
		Where("reserver_id = ?", reserverID).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	return reservationsToDomain(models), nil
}

func reservationDomainToModel(r *domain.Reservation) *reservationModel {
	return &reservationModel{
		ID:          r.ID,
		WishID:      r.WishID,
		ReserverID:  r.ReserverID,
		Quantity:    r.Quantity,
		IsAnonymous: r.IsAnonymous,
		IsDone:      r.IsDone,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func reservationModelToDomain(m *reservationModel) *domain.Reservation {
	return &domain.Reservation{
		ID:          m.ID,
		WishID:      m.WishID,
		ReserverID:  m.ReserverID,
		Quantity:    m.Quantity,
		IsAnonymous: m.IsAnonymous,
		IsDone:      m.IsDone,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func reservationsToDomain(models []reservationModel) []*domain.Reservation {
	out := make([]*domain.Reservation, 0, len(models))
	for i := range models {
		out = append(out, reservationModelToDomain(&models[i]))
	}
	return out
}
