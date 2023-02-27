package transactional

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type DefaultTemplate struct {
	db *sqlx.DB
}

func NewTemplate(db *sqlx.DB) *DefaultTemplate {
	return &DefaultTemplate{db: db}
}

func (t *DefaultTemplate) Execute(ctx context.Context, callback func(tx Tx) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	err = callback(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
