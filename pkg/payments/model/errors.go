package model

import "errors"

var (
	ErrPaymentWrongStatus = errors.New("payment has a wrong status for this action")
	ErrPaymentNotFound    = errors.New("payment not found")
)
