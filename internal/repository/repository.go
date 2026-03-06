package repository

import (
	"log/slog"
	"time"

	"test-junior-go/internal/model"

	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub *model.Subscription) error
	GetByID(id uint) (*model.Subscription, error)

	GetAll(
		limit int,
		offset int,
		userID string,
		serviceName string,
	) ([]model.Subscription, error)

	Update(sub *model.Subscription) error
	Delete(id uint) error

	GetTotal(
		userID string,
		serviceName string,
		startDate time.Time,
		endDate time.Time,
	) (int64, error)
}

type subscriptionRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewSubscriptionRepository(db *gorm.DB, logger *slog.Logger) SubscriptionRepository {
	return &subscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *subscriptionRepository) Create(sub *model.Subscription) error {
	r.logger.Info("creating subscription",
		"service_name", sub.ServiceName,
		"user_id", sub.UserID.String(),
		"price", sub.Price,
	)

	if err := r.db.Create(sub).Error; err != nil {
		r.logger.Error("failed to create subscription", "error", err,
			"service_name", sub.ServiceName,
			"user_id", sub.UserID.String(),
		)
		return err
	}

	r.logger.Info("subscription created", "id", sub.ID)
	return nil
}

func (r *subscriptionRepository) GetByID(id uint) (*model.Subscription, error) {
	var sub model.Subscription
	r.logger.Info("fetching subscription by id", "id", id)

	err := r.db.First(&sub, id).Error
	if err != nil {
		r.logger.Error("failed to get subscription", "id", id, "error", err)
		return nil, err
	}

	r.logger.Info("subscription fetched", "id", sub.ID, "service_name", sub.ServiceName)
	return &sub, nil
}

func (r *subscriptionRepository) GetAll(
	limit int,
	offset int,
	userID string,
	serviceName string,
) ([]model.Subscription, error) {

	var subs []model.Subscription

	query := r.db.Model(&model.Subscription{})

	r.logger.Info("fetching subscriptions list",
		"limit", limit,
		"offset", offset,
		"user_id", userID,
		"service_name", serviceName,
	)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&subs).Error; err != nil {
		r.logger.Error("failed to get subscriptions", "error", err)
		return nil, err
	}

	r.logger.Info("subscriptions fetched", "count", len(subs))
	return subs, nil
}

func (r *subscriptionRepository) Update(sub *model.Subscription) error {
	r.logger.Info("updating subscription", "id", sub.ID)

	if err := r.db.Save(sub).Error; err != nil {
		r.logger.Error("failed to update subscription", "id", sub.ID, "error", err)
		return err
	}

	r.logger.Info("subscription updated", "id", sub.ID)
	return nil
}

func (r *subscriptionRepository) Delete(id uint) error {

	r.logger.Info("deleting subscription", "id", id)

	res := r.db.Delete(&model.Subscription{}, id)

	if res.Error != nil {
		r.logger.Error("failed to delete subscription", "id", id, "error", res.Error)
		return res.Error
	}

	if res.RowsAffected == 0 {
		r.logger.Warn("subscription not found for delete", "id", id)
		return gorm.ErrRecordNotFound
	}

	r.logger.Info("subscription deleted", "id", id)
	return nil
}

func (r *subscriptionRepository) GetTotal(
	userID string,
	serviceName string,
	startDate time.Time,
	endDate time.Time,
) (int64, error) {

	var total int64

	r.logger.Info("calculating total subscription cost",
		"user_id", userID,
		"service_name", serviceName,
		"start_date", startDate,
		"end_date", endDate,
	)

	query := r.db.Model(&model.Subscription{}).
		Select(`
COALESCE(SUM(
	price * (
		DATE_PART('year', AGE(
			LEAST(COALESCE(end_date, ?), ?),
			GREATEST(start_date, ?)
		)) * 12 +
		DATE_PART('month', AGE(
			LEAST(COALESCE(end_date, ?), ?),
			GREATEST(start_date, ?)
		)) + 1
	)
),0)
`, endDate, endDate, startDate, endDate, endDate, startDate).
		Where("user_id = ?", userID).
		Where("start_date <= ?", endDate).
		Where("(end_date IS NULL OR end_date >= ?)", startDate)

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if err := query.Scan(&total).Error; err != nil {
		r.logger.Error("failed to calculate total subscription cost",
			"error", err,
			"user_id", userID,
		)
		return 0, err
	}

	r.logger.Info("total subscription cost calculated", "total", total)
	return total, nil
}
