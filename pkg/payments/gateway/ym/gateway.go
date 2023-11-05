package ym

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"job-interviewer/pkg/payments/gateway"
	"job-interviewer/pkg/payments/model"
	"job-interviewer/pkg/structs"
	"net/http"
	"strconv"
)

const (
	currency    = "RUB"
	paymentsUrl = "https://api.yookassa.ru/v3/payments"

	confirmationTypeRedirect = "redirect"
)

var (
	statusMap = structs.NewBidirectionalMap[string, model.Status](
		[]structs.Pair[string, model.Status]{
			{Left: "pending", Right: model.StatusPending},
			{Left: "succeeded", Right: model.StatusPaid},
			{Left: "canceled", Right: model.StatusCanceled},
		},
	)
)

type (
	amount struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	}

	confirmation struct {
		Type      string `json:"type"`
		ReturnUrl string `json:"return_url"`
	}

	createPaymentRequest struct {
		Amount       amount       `json:"amount"`
		Capture      bool         `json:"capture"`
		Confirmation confirmation `json:"confirmation"`
		Description  string       `json:"description"`
	}

	createPaymentResponse struct {
		ID           string `json:"id"`
		Status       string `json:"status"`
		Confirmation struct {
			Type            string `json:"type"`
			ConfirmationUrl string `json:"confirmation_url"`
		} `json:"confirmation"`
		Test bool `json:"test"`
	}

	getPaymentResponse struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	Gateway struct {
		shopID     int64
		secretKey  string
		httpClient httpClient
	}
)

func NewYMGateway(
	shopID int64,
	secretKey string,
	httpClient httpClient,
) *Gateway {
	return &Gateway{
		shopID:     shopID,
		secretKey:  secretKey,
		httpClient: httpClient,
	}
}

func (g *Gateway) CreatePayment(ctx context.Context, in *gateway.CreatePaymentIn) (*gateway.CreatePaymentOut, error) {
	payload := createPaymentRequest{
		Amount: amount{
			Value:    float64(in.Amount),
			Currency: currency,
		},
		Capture: true,
		Confirmation: confirmation{
			Type:      confirmationTypeRedirect,
			ReturnUrl: "https://www.google.com/example",
		},
		Description: in.Description,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		paymentsUrl,
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Idempotence-Key", in.IDK.String())
	request.Header.Add(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			basicAuth(strconv.FormatInt(g.shopID, 10), g.secretKey),
		),
	)

	resp, err := g.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &createPaymentResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &gateway.CreatePaymentOut{
		ExternalID:  model.ExternalID(res.ID),
		RedirectURl: res.Confirmation.ConfirmationUrl,
	}, nil
}

func (g *Gateway) GetPaymentStatus(ctx context.Context, ID model.ExternalID) (*model.Status, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s", paymentsUrl, ID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add(
		"Authorization",
		fmt.Sprintf(
			"Basic %s",
			basicAuth(strconv.FormatInt(g.shopID, 10), g.secretKey),
		),
	)

	resp, err := g.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &getPaymentResponse{}
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	statusOut, ok := statusMap.GetRight(res.Status)
	if !ok {
		return nil, model.ErrPaymentWrongStatus
	}

	return &statusOut, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
