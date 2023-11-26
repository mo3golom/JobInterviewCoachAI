package acceptanswer

import (
	"context"
	"errors"
	"fmt"
	interviewerContracts "job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/telegram/handlers"
	"job-interviewer/internal/telegram/service"
	"job-interviewer/pkg/language"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	acceptAnswerUC  interviewerContracts.AcceptAnswerUseCase
	service         service.Service
	languageStorage language.Storage
}

func NewHandler(
	auc interviewerContracts.AcceptAnswerUseCase,
	s service.Service,
) *Handler {
	return &Handler{
		acceptAnswerUC:  auc,
		service:         s,
		languageStorage: configLanguage(),
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if request.Voice != nil {
		_, err := sender.Send(
			model.NewResponse().SetText(
				h.languageStorage.GetText(language.Russian, textKeyVoiceMessageIsUnsupported),
			),
		)

		return err
	}

	if request.Message == nil || request.Message.Text == "" {
		return nil
	}

	answerMessageID, err := sender.Send(
		model.NewResponse().SetText(
			fmt.Sprintf(
				"%s %s",
				handlers.RobotPrefix,
				h.languageStorage.GetText(language.English, textKeyProcessingAnswer),
			),
		),
	)
	if err != nil {
		return err
	}
	err = h.acceptAnswerUC.AcceptAnswer(
		ctx,
		interviewerContracts.AcceptAnswerIn{
			Answer: request.Message.Text,
			UserID: request.User.OriginalID,
		},
	)
	if errors.Is(err, interviewerContracts.ErrEmptyActiveInterview) {
		return nil
	}
	if err != nil {
		return err
	}

	return h.service.GetNextQuestion(ctx, request, sender, answerMessageID)
}
