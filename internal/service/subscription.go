package service

import (
	"errors"
	"log/slog"
	"time"

	suberrors "test-junior-go/internal/errors"
	"test-junior-go/internal/model"
	"test-junior-go/internal/repository"
)

type SubscriptionService interface {
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

type subscriptionService struct {
	repo   repository.SubscriptionRepository
	logger *slog.Logger
}

func NewSubscriptionService(
	repo repository.SubscriptionRepository,
	logger *slog.Logger,
) SubscriptionService {
	return &subscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *subscriptionService) Create(sub *model.Subscription) error {

	if sub.Price <= 0 {
		s.logger.Warn("validation failed: invalid price",
			"price", sub.Price,
			"user_id", sub.UserID,
		)
		return suberrors.ErrInvalidPrice
	}

	if sub.ServiceName == "" {
		s.logger.Warn("validation failed: empty service_name",
			"user_id", sub.UserID,
		)
		return suberrors.ErrEmptyServiceName
	}

	if err := s.repo.Create(sub); err != nil {
		s.logger.Error("failed to create subscription",
			"error", err,
			"user_id", sub.UserID,
		)
		return err
	}

	return nil
}

func (s *subscriptionService) GetByID(id uint) (*model.Subscription, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, suberrors.ErrSubscriptionNotFound) {
			s.logger.Warn("subscription not found", "id", id)
			return nil, err
		}

		s.logger.Error("failed to get subscription", "id", id, "error", err)
		return nil, err
	}

	return sub, nil
}
func (s *subscriptionService) GetAll(
	limit int,
	offset int,
	userID string,
	serviceName string,
) ([]model.Subscription, error) {

	subs, err := s.repo.GetAll(limit, offset, userID, serviceName)
	if err != nil {
		s.logger.Error("failed to get subscriptions",
			"limit", limit,
			"offset", offset,
			"user_id", userID,
			"service_name", serviceName,
			"error", err,
		)
		return nil, err
	}

	return subs, nil
}

func (s *subscriptionService) Update(sub *model.Subscription) error {

	if sub.Price <= 0 {
		s.logger.Warn("validation failed: invalid price",
			"price", sub.Price,
			"id", sub.ID,
		)
		return suberrors.ErrInvalidPrice
	}

	if err := s.repo.Update(sub); err != nil {
		s.logger.Error("failed to update subscription",
			"id", sub.ID,
			"error", err,
		)
		return err
	}

	return nil
}

func (s *subscriptionService) Delete(id uint) error {
	err := s.repo.Delete(id)
	if err != nil {

		if errors.Is(err, suberrors.ErrSubscriptionNotFound) {
			s.logger.Warn("subscription not found for delete", "id", id)
			return err
		}

		s.logger.Error("failed to delete subscription",
			"id", id,
			"error", err,
		)
		return err
	}

	return nil
}

func (s *subscriptionService) GetTotal(
    userID string,
    serviceName string,
    startDate time.Time,
    endDate time.Time,
) (int64, error) {

    total, err := s.repo.GetTotal(userID, serviceName, startDate, endDate)
    if err != nil {
        s.logger.Error("failed to calculate total subscription cost",
            "user_id", userID,
            "service_name", serviceName,
            "error", err,
        )
        return 0, err
    }

    return total, nil
}

