package errors

import "errors"

var (
	ErrInvalidPrice     = errors.New("invalid price: must be greater than zero")
	ErrEmptyServiceName = errors.New("service_name is required")

	ErrSubscriptionNotFound = errors.New("subscription not found")

	ErrInvalidStartDate = errors.New("invalid start date format, expected MM-YYYY")
	ErrInvalidEndDate   = errors.New("invalid end date format, expected MM-YYYY")
)
