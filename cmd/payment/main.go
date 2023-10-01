package main

import (
	"context"
	"github.com/google/uuid"
	"job-interviewer/pkg/payments"
	"job-interviewer/pkg/payments/gateway/ym"
	"net/http"
)

func main() {
	ctx := context.Background()
	gateway := ym.NewYMGateway(
		255279,
		"test_NJ9ngSgKIDntPfBDKhvK5t4UUIlqWbCFg02e7i3Wb8Q",
		&http.Client{},
	)

	gateway.CreatePayment(
		ctx,
		&payments.GatewayCreatePaymentIn{
			IDK:         uuid.New(),
			Description: "test pay",
			Amount:      10,
		},
	)

	return
}
