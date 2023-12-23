package paysubscription

import (
	"context"
	interviewContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers/command"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/keyboard"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	languageStorage language.Storage
	subscriptionUC  interviewContracts.SubscriptionUseCase
}

func NewHandler(
	subscriptionUC interviewContracts.SubscriptionUseCase,
) *Handler {
	return &Handler{
		subscriptionUC:  subscriptionUC,
		languageStorage: configLanguage(),
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	userLang := language.Russian
	result, err := h.subscriptionUC.CreatePayment(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	payInlineKeyboard, err := keyboard.BuildInlineKeyboardList(
		keyboard.BuildInlineKeyboardIn{
			Buttons: []keyboard.InlineButton{
				{
					Value: h.languageStorage.GetText(userLang, textKeyPay),
					Data:  []string{result.RedirectURL},
					Type:  keyboard.ButtonUrl,
				},
				{
					Value: h.languageStorage.GetText(userLang, textKeyCheckPayment),
					Data:  []string{command.CheckPaymentCommand},
					Type:  keyboard.ButtonData,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse().
			SetText(h.languageStorage.GetText(userLang, textKeyPayRedirectURL)).
			SetInlineKeyboardMarkup(payInlineKeyboard),
	)

	return err
}

func (h *Handler) Command() string {
	return command.PaySubscriptionCommand
}

func (h *Handler) Aliases() []string {
	return nil
}
