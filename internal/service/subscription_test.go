package service

import (
	"log/slog"
	"testing"
	"time"

	"test-junior-go/internal/model"

	"github.com/stretchr/testify/assert"
)

type mockRepo struct {
	createCalled bool
}

func (m *mockRepo) Create(sub *model.Subscription) error {
	m.createCalled = true
	return nil
}

func (m *mockRepo) GetByID(id uint) (*model.Subscription, error) {
	return nil, nil
}

func (m *mockRepo) GetAll(limit int, offset int, userID string, serviceName string) ([]model.Subscription, error) {
	return nil, nil
}

func (m *mockRepo) Update(sub *model.Subscription) error {
	return nil
}

func (m *mockRepo) Delete(id uint) error {
	return nil
}

func (m *mockRepo) GetTotal(userID string, serviceName string, startDate time.Time, endDate time.Time) (int64, error) {
	return 1000, nil
}

func TestCreateSubscription(t *testing.T) {

	repo := &mockRepo{}
	logger := slog.Default()

	svc := NewSubscriptionService(repo, logger)

	sub := &model.Subscription{
		ServiceName: "Netflix",
		Price:       500,
	}

	err := svc.Create(sub)

	assert.NoError(t, err)
	assert.True(t, repo.createCalled)
}

func TestCreateSubscription_InvalidPrice(t *testing.T) {

	repo := &mockRepo{}
	logger := slog.Default()

	svc := NewSubscriptionService(repo, logger)

	sub := &model.Subscription{
		ServiceName: "Netflix",
		Price:       0,
	}

	err := svc.Create(sub)

	assert.Error(t, err)
}

func TestGetTotal(t *testing.T) {

	repo := &mockRepo{}
	logger := slog.Default()

	svc := NewSubscriptionService(repo, logger)

	startDate, _ := time.Parse("01-2006", "01-2025")
	endDate, _ := time.Parse("01-2006", "12-2025")

	total, err := svc.GetTotal(
		"user",
		"",
		startDate,
		endDate,
	)

	assert.NoError(t, err)
	assert.Equal(t, int64(1000), total)
}
