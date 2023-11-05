package payments

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/pkg/payments/model"
	"job-interviewer/pkg/transactional"
)

type (
	sqlxNewPayment struct {
		ID          uuid.UUID `db:"id"`
		UserID      uuid.UUID `db:"user_id"`
		Amount      int64     `db:"amount_penny"`
		Type        string    `db:"type"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
	}

	sqlxUpdatePayment struct {
		ID          uuid.UUID         `db:"id"`
		Status      string            `db:"status"`
		ExternalID  *model.ExternalID `db:"external_id"`
		RedirectUrl *string           `db:"redirect_url"`
	}

	sqlxPayment struct {
		ID          uuid.UUID         `db:"id"`
		UserID      uuid.UUID         `db:"user_id"`
		Amount      int64             `db:"amount_penny"`
		Type        string            `db:"type"`
		Description string            `db:"description"`
		Status      string            `db:"status"`
		ExternalID  *model.ExternalID `db:"external_id"`
		RedirectUrl *string           `db:"redirect_url"`
	}

	DefaultRepository struct {
		db *sqlx.DB
	}
)

func (r *DefaultRepository) CreatePayment(ctx context.Context, tx transactional.Tx, in *model.Payment) error {
	query := `
		INSERT 
		INTO payment (id, user_id, amount_penny, type, description, status) 
		VALUES (:id, :user_id, :amount_penny, :type, :description, :status)
		ON CONFLICT (id) DO NOTHING
    `

	data := sqlxNewPayment{
		ID:          in.ID,
		UserID:      in.UserID,
		Amount:      int64(in.Amount),
		Type:        string(in.Type),
		Description: in.Description,
		Status:      string(model.StatusNew),
	}
	_, err := tx.NamedExecContext(ctx, query, data)
	return err
}

func (r *DefaultRepository) UpdatePayment(ctx context.Context, tx transactional.Tx, in *model.Payment) error {
	query := `
		UPDATE payment
		SET 
		    status=:status,
		    redirect_url=:redirect_url,
		    external_id=:external_id,
		    updated_at=now()
        WHERE id=:id 
    `
	data := sqlxUpdatePayment{
		ID:          in.ID,
		Status:      string(in.Status),
		RedirectUrl: in.RedirectURL,
		ExternalID:  in.ExternalID,
	}
	_, err := tx.NamedExecContext(ctx, query, data)
	return err
}

func (r *DefaultRepository) GetPaymentByExternalID(ctx context.Context, ID model.ExternalID) (*model.Payment, error) {
	const query = `
		SELECT p.id, p.user_id, p.external_id, p.amount_penny, p.type, p.description, p.redirect_url, p.status
		FROM payment as p
		WHERE p.external_id = $1
    `

	var results []sqlxPayment
	err := r.db.SelectContext(
		ctx,
		&results,
		query,
		ID,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, model.ErrPaymentNotFound
	}

	return convertToPayment(&results[0]), nil
}

func (r *DefaultRepository) GetActivePaymentByUserID(ctx context.Context, userID uuid.UUID) (*model.Payment, error) {
	const query = `
		SELECT p.id, p.user_id, p.external_id, p.amount_penny, p.type, p.description, p.redirect_url, p.status
		FROM payment as p
		WHERE p.user_id = $1 and (status = $2 or status = $3)
    `

	var results []sqlxPayment
	err := r.db.SelectContext(
		ctx,
		&results,
		query,
		userID,
		model.StatusPending,
		model.StatusNew,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, model.ErrPaymentNotFound
	}

	return convertToPayment(&results[0]), nil
}

func convertToPayment(in *sqlxPayment) *model.Payment {
	return &model.Payment{
		ID:          in.ID,
		ExternalID:  in.ExternalID,
		UserID:      in.UserID,
		Amount:      model.Penny(in.Amount),
		Type:        model.Type(in.Type),
		Description: in.Description,
		RedirectURL: in.RedirectUrl,
		Status:      model.Status(in.Status),
	}
}
