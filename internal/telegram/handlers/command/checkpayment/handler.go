package checkpayment

import (
	"context"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/payments"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	languageStorage language.Storage
	paymentsService payments.Service
}

func NewHandler(
	paymentsService payments.Service,
) *Handler {
	return &Handler{
		paymentsService: paymentsService,
		languageStorage: configLanguage(),
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := language.Russian
	result, err := h.paymentsService.CheckPendingPaymentWithResult(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}
	if !result.Paid {
		return nil
	}

	_, err = sender.Send(
		model.NewResponse().
			SetText(h.languageStorage.GetText(userLang, textKeyPaymentSuccess)),
	)
	return err
}

func (h *Handler) Command() string {
	return command.CheckPaymentCommand
}

func (h *Handler) Aliases() []string {
	return nil
}
