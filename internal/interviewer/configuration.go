package interviewer

import (
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/contracts"
	"job-interviewer/internal/interviewer/flow"
	"job-interviewer/internal/interviewer/gpt"
	interview2 "job-interviewer/internal/interviewer/service/interview"
	"job-interviewer/internal/interviewer/service/payments"
	"job-interviewer/internal/interviewer/service/subscription"
	"job-interviewer/internal/interviewer/storage/interview"
	"job-interviewer/internal/interviewer/storage/messages"
	"job-interviewer/internal/interviewer/storage/user"
	"job-interviewer/internal/interviewer/usecase/acceptanswer"
	"job-interviewer/internal/interviewer/usecase/finishinterview"
	"job-interviewer/internal/interviewer/usecase/getinterview"
	"job-interviewer/internal/interviewer/usecase/getnextquestion"
	"job-interviewer/internal/interviewer/usecase/startinterview"
	subscription2 "job-interviewer/internal/interviewer/usecase/subscription"
	user2 "job-interviewer/internal/interviewer/usecase/user"
	externalPayments "job-interviewer/pkg/payments"
	externalSubscription "job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
	"job-interviewer/pkg/variables"
)

type (
	ConfigurationUseCases struct {
		StartInterview  contracts.StartInterviewUseCase
		FinishInterview contracts.FinishInterviewUseCase
		GetNextQuestion contracts.GetNextQuestionUseCase
		AcceptAnswer    contracts.AcceptAnswerUseCase
		GetInterview    contracts.GetInterviewUsecase
		User            contracts.UserUseCase
		Subscription    contracts.SubscriptionUseCase
	}

	Configuration struct {
		UseCases *ConfigurationUseCases
	}
)

func NewConfiguration(
	db *sqlx.DB,
	transactionalTemplate transactional.Template,
	gptGateway gpt.Gateway,
	externalSubscriptionService externalSubscription.Service,
	externalPaymentsService externalPayments.Service,
	variables variables.Repository,
) *Configuration {
	interviewStorage := interview.NewStorage(db)
	userStorage := user.NewStorage(db)
	messagesStorage := messages.NewStorage(db)

	subscriptionService := subscription.NewService(externalSubscriptionService, messagesStorage)
	interviewService := interview2.NewInterviewService(gptGateway, interviewStorage, messagesStorage, transactionalTemplate, subscriptionService)
	interviewFlow := flow.NewDefaultInterviewFlow(interviewService)
	paymentsService, err := payments.NewDefaultService(
		externalPaymentsService,
		externalSubscriptionService,
		transactionalTemplate,
		variables,
	)
	if err != nil {
		panic(err)
	}

	useCases := &ConfigurationUseCases{
		StartInterview:  startinterview.NewUseCase(interviewService, interviewFlow, subscriptionService),
		FinishInterview: finishinterview.NewUseCase(interviewFlow, subscriptionService),
		GetNextQuestion: getnextquestion.NewUseCase(interviewFlow, subscriptionService),
		AcceptAnswer:    acceptanswer.NewUseCase(interviewFlow, subscriptionService),
		GetInterview:    getinterview.NewUseCase(interviewService),
		User:            user2.NewUseCase(userStorage, transactionalTemplate, externalSubscriptionService),
		Subscription:    subscription2.NewUseCase(paymentsService),
	}

	return &Configuration{UseCases: useCases}
}
