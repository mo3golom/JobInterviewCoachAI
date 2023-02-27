package processinterview

import (
	"context"
	"errors"
	"fmt"
	"job-interviewer/internal/contracts"
	"job-interviewer/pkg/telegram"
	"job-interviewer/pkg/telegram/model"
)

type Handler struct {
	getNextQuestionUC contracts.GetNextQuestionUseCase
	acceptAnswerUC    contracts.AcceptAnswerUseCase
	finishInterviewUC contracts.FinishInterviewUseCase
}

func NewHandler(
	guc contracts.GetNextQuestionUseCase,
	fuc contracts.FinishInterviewUseCase,
	auc contracts.AcceptAnswerUseCase,
) *Handler {
	return &Handler{
		getNextQuestionUC: guc,
		finishInterviewUC: fuc,
		acceptAnswerUC:    auc,
	}
}

func (h *Handler) Handle(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	if request.Message != nil {
		err := h.acceptAnswer(ctx, request, sender)
		if err != nil {
			return err
		}
	}

	return h.getNextQuestion(ctx, request, sender)
}

func (h *Handler) getNextQuestion(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	response := model.NewResponse(request.Chat.ID)

	question, err := h.getNextQuestionUC.GetNextQuestion(ctx, request.User.OriginalID)
	if errors.Is(err, contracts.ErrNextQuestionEmpty) {
		return h.finishInterview(ctx, request, sender)
	}
	if errors.Is(err, contracts.ErrEmptyActiveInterview) {
		_, err = sender.Send(response.SetText(noActiveInterviewText))
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	_, err = sender.Send(
		response.SetText(fmt.Sprintf(questionText, question.Text)),
	)
	return err
}

func (h *Handler) finishInterview(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	err := h.finishInterviewUC.FinishInterview(ctx, request.User.OriginalID)
	if err != nil {
		return err
	}

	_, err = sender.Send(
		model.NewResponse(request.Chat.ID).
			SetText(finishText),
	)
	return err
}

func (h *Handler) acceptAnswer(ctx context.Context, request *model.Request, sender telegram.Sender) error {
	messageID, err := sender.Send(
		model.NewResponse(request.Chat.ID).SetText(processingAnswerText),
	)
	response := model.NewResponse(request.Chat.ID)

	out, err := h.acceptAnswerUC.AcceptAnswer(
		ctx,
		contracts.AcceptAnswerIn{
			Answer: request.Message.Text,
			UserID: request.User.OriginalID,
		},
	)
	if errors.Is(err, contracts.ErrNextQuestionEmpty) {
		return h.finishInterview(ctx, request, sender)
	}
	if errors.Is(err, contracts.ErrEmptyActiveInterview) {
		err = sender.Update(messageID, response.SetText(noActiveInterviewText))
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	err = sender.Update(messageID, response.SetText(
		fmt.Sprintf(AnswerAnswerText, out),
	))
	if err != nil {
		return err
	}

	return nil
}
