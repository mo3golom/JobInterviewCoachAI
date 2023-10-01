package payments

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		ID          uuid.UUID `db:"id"`
		Status      string    `db:"status"`
		ExternalID  string    `db:"external_id"`
		RedirectUrl string    `db:"redirect_url"`
	}

	sqlxPayment struct {
		ID          uuid.UUID `db:"id"`
		UserID      uuid.UUID `db:"user_id"`
		Amount      int64     `db:"amount_penny"`
		Type        string    `db:"type"`
		Description string    `db:"description"`
		Status      string    `db:"status"`
		ExternalID  string    `db:"external_id"`
		RedirectUrl string    `db:"redirect_url"`
	}

	DefaultRepository struct {
		db *sqlx.DB
	}
)

func (r *DefaultRepository) CreatePayment(ctx context.Context, in *NewPayment) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT 
		INTO payment (id, user_id, amount_penny, type, description, status) 
		VALUES (:id, :user_id, :amount_penny, :type, :description, :status)
		ON CONFLICT DO NOTHING 
    `

	data := sqlxNewPayment{
		ID:          in.ID,
		Amount:      int64(in.Amount),
		Type:        string(in.Type),
		Description: in.Description,
		Status:      string(StatusNew),
	}
	_, err = tx.NamedExecContext(ctx, query, data)
	return err
}

func (r *DefaultRepository) UpdatePayment(ctx context.Context, in *Payment) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

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
		ExternalID:  string(in.ExternalID),
	}
	_, err = tx.NamedExecContext(ctx, query, data)
	return err
}

func (r *DefaultRepository) GetPaymentByExternalID(ctx context.Context, ID ExternalID) (*Payment, error) {
	const query = `
		SELECT p.id, p.user_id, p.external_id, p.amount, p.type, p.description, p.redirect_url, p.status
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
		return nil, ErrPaymentNotFound
	}

	return convertToPayment(&results[0]), nil
}

func (r *DefaultRepository) GetPendingPaymentByUserID(ctx context.Context, userID uuid.UUID) (*Payment, error) {
	const query = `
		SELECT p.id, p.user_id, p.external_id, p.amount, p.type, p.description, p.redirect_url, p.status
		FROM payment as p
		WHERE p.user_id = $1 and status = $2
    `

	var results []sqlxPayment
	err := r.db.SelectContext(
		ctx,
		&results,
		query,
		userID,
		StatusPending,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrPaymentNotFound
	}

	return convertToPayment(&results[0]), nil
}

func convertToPayment(in *sqlxPayment) *Payment {
	return &Payment{
		ID:          in.ID,
		ExternalID:  ExternalID(in.ExternalID),
		UserID:      in.UserID,
		Amount:      Penny(in.Amount),
		Type:        Type(in.Type),
		Description: in.Description,
		RedirectURL: in.RedirectUrl,
		Status:      Status(in.Status),
	}
}
