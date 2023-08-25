package messages

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
	"time"
)

type (
	sqlxMessage struct {
		InterviewID uuid.UUID `db:"interview_id"`
		Content     string    `db:"content"`
		Role        string    `db:"role"`
		CreatedAt   time.Time `db:"created_at"`
	}

	DefaultStorage struct {
		db *sqlx.DB
	}
)

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (d *DefaultStorage) CreateMessage(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, in *model.Message) error {
	query := `
		INSERT 
		INTO interview_messages (interview_id, content, role) 
		VALUES (:interview_id, :content, :role)
		ON CONFLICT DO NOTHING 
    `

	message := sqlxMessage{
		InterviewID: interviewID,
		Content:     in.Content,
		Role:        string(in.Role),
	}
	_, err := tx.NamedExecContext(ctx, query, message)
	return err
}

func (d *DefaultStorage) GetMessagesByInterviewID(ctx context.Context, interviewID uuid.UUID) ([]model.Message, error) {
	query := `
		SELECT content, role, created_at
		FROM  interview_messages
		WHERE interview_id = $1
    `

	var results []sqlxMessage

	err := d.db.SelectContext(
		ctx,
		&results,
		query,
		interviewID,
	)
	if err != nil {
		return nil, err
	}

	out := make([]model.Message, 0, len(results))
	for _, res := range results {
		out = append(out, model.Message{
			Role:      model.Role(res.Role),
			Content:   res.Content,
			CreatedAt: res.CreatedAt,
		})
	}
	return out, nil
}
