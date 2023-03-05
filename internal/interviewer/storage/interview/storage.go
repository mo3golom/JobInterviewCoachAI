package interview

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	model2 "job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

type sqlxInterview struct {
	ID            uuid.UUID `db:"id"`
	UserID        uuid.UUID `db:"user_id"`
	Status        string    `db:"status"`
	QuestionCount int64     `db:"question_count"`
	JobPosition   string    `db:"job_position"`
	JobLevel      string    `db:"job_level"`
}

type DefaultStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) CreateInterview(ctx context.Context, tx transactional.Tx, interview *model2.Interview) error {
	query := `
		INSERT 
		INTO interview (id, user_id, status, job_position, job_level, question_count) 
		VALUES (:id, :user_id, :status, :job_position, :job_level, :question_count)
		ON CONFLICT DO NOTHING 
    `

	in := sqlxInterview{
		ID:            interview.ID,
		UserID:        interview.UserID,
		Status:        string(interview.Status),
		QuestionCount: interview.QuestionsCount,
		JobPosition:   interview.JobInfo.Position,
		JobLevel:      string(interview.JobInfo.Level),
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}

func (s *DefaultStorage) UpdateInterview(ctx context.Context, tx transactional.Tx, interview *model2.Interview) error {
	query := `
		UPDATE interview
		SET  
		    status=:status,
		    updated_at=now()
        WHERE id=:id 
    `

	in := sqlxInterview{
		ID:     interview.ID,
		Status: string(interview.Status),
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}

func (s *DefaultStorage) FindActiveInterviewByUserID(ctx context.Context, tx transactional.Tx, userID uuid.UUID) (*model2.Interview, error) {
	query := `
		SELECT i.id, i.user_id, i.status, i.job_position, i.job_level, i.question_count
		FROM interview as i
		WHERE i.user_id = $1 and i.status = $2
    `

	var results []sqlxInterview
	err := tx.SelectContext(
		ctx,
		&results,
		query,
		userID,
		model2.InterviewStatusStarted,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrEmptyInterviewResult
	}

	return convertInterview(&results[0]), nil
}

func convertInterview(in *sqlxInterview) *model2.Interview {
	return &model2.Interview{
		ID:     in.ID,
		UserID: in.UserID,
		Status: model2.InterviewStatus(in.Status),
		JobInfo: model2.JobInfo{
			Position: in.JobPosition,
			Level:    model2.JobLevel(in.JobLevel),
		},
		QuestionsCount: in.QuestionCount,
	}
}
