package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"test-junior-go/internal/dto"
	suberrors "test-junior-go/internal/errors"
	"test-junior-go/internal/model"
	"test-junior-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *slog.Logger
}

func NewSubscriptionHandler(
	service service.SubscriptionService,
	logger *slog.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

// CreateSubscription godoc
// @Summary Create subscription
// @Description Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {

	var req dto.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user_id",
		})
		return
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_date format, expected MM-YYYY",
		})
		return
	}

	var endDate *time.Time
	if req.EndDate != "" {
		t, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid end_date format, expected MM-YYYY",
			})
			return
		}
		endDate = &t
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      uid,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := h.service.Create(sub); err != nil {

		if errors.Is(err, suberrors.ErrInvalidPrice) ||
			errors.Is(err, suberrors.ErrEmptyServiceName) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		h.logger.Error("failed to create subscription",
			"error", err,
			"user_id", req.UserID,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// GetSubscriptionByID godoc
// @Summary Get subscription
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	sub, err := h.service.GetByID(uint(id))
	if err != nil {

		if errors.Is(err, suberrors.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "subscription not found",
			})
			return
		}

		h.logger.Error("failed to get subscription",
			"id", id,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// GetSubscriptions godoc
// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} model.Subscription
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {

	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	subs, err := h.service.GetAll(limit, offset, userID, serviceName)
	if err != nil {

		h.logger.Error("failed to get subscriptions", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// DeleteSubscription godoc
// @Summary Delete subscription
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {

	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	err = h.service.Delete(uint(id))
	if err != nil {

		if errors.Is(err, suberrors.ErrSubscriptionNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "subscription not found",
			})
			return
		}

		h.logger.Error("failed to delete subscription",
			"id", id,
			"error", err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	h.logger.Info("subscription deleted", "id", id)

	c.Status(http.StatusNoContent)
}

// GetTotalCost godoc
// @Summary Calculate total subscription cost
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User ID"
// @Param service_name query string false "Service name"
// @Param start query string true "Start period (MM-YYYY)"
// @Param end query string true "End period (MM-YYYY)"
// @Success 200 {object} map[string]int64
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotal(c *gin.Context) {

	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	start := c.Query("start")
	end := c.Query("end")

	startDate, err := time.Parse("01-2006", start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start format",
		})
		return
	}

	endDate, err := time.Parse("01-2006", end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end format",
		})
		return
	}

	total, err := h.service.GetTotal(userID, serviceName, startDate, endDate)
	if err != nil {

		h.logger.Error("failed to calculate total", "error", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
	})
}
