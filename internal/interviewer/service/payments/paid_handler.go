package payments

import (
	"context"
	"github.com/google/uuid"
	externalSubscription "job-interviewer/pkg/subscription"
	"job-interviewer/pkg/transactional"
)

type defaultPaidHandler struct {
	externalSubscriptionService externalSubscription.Service
	transactionalTemplate       transactional.Template
}

func (h *defaultPaidHandler) Handle(ctx context.Context, userID uuid.UUID) error {
	return h.transactionalTemplate.Execute(ctx, func(tx transactional.Tx) error {
		return h.externalSubscriptionService.ActivateSubscription(
			ctx,
			tx,
			userID,
			externalSubscription.PlanMonth,
		)
	})
}
