package interview

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

type sqlxInterview struct {
	ID            uuid.UUID `db:"id"`
	UserID        uuid.UUID `db:"user_id"`
	Status        string    `db:"status"`
	QuestionCount int64     `db:"question_count"`
	JobPosition   string    `db:"job_position"`
	JobLevel      string    `db:"job_level"`
	State         string    `db:"state"`
}

type DefaultStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) CreateInterview(ctx context.Context, tx transactional.Tx, interview *model.Interview) error {
	const query = `
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
		JobPosition:   string(interview.JobInfo.Position),
		JobLevel:      "unknown",
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}

func (s *DefaultStorage) UpdateInterview(ctx context.Context, tx transactional.Tx, interview *model.Interview) error {
	const query = `
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

func (s *DefaultStorage) UpdateInterviewState(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, state model.InterviewState) error {
	const query = `
		UPDATE interview
		SET  
		    state=:state,
		    updated_at=now()
        WHERE id=:id 
    `

	in := sqlxInterview{
		ID:    interviewID,
		State: string(state),
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}

func (s *DefaultStorage) FindActiveInterviewByUserID(ctx context.Context, tx transactional.Tx, userID uuid.UUID) (*model.Interview, error) {
	const query = `
		SELECT i.id, i.user_id, i.status, i.job_position, i.job_level, i.question_count, i.state
		FROM interview as i
		WHERE i.user_id = $1 and i.status = $2
    `

	var results []sqlxInterview
	err := tx.SelectContext(
		ctx,
		&results,
		query,
		userID,
		model.InterviewStatusStarted,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrEmptyInterviewResult
	}

	return convertInterview(&results[0]), nil
}

func convertInterview(in *sqlxInterview) *model.Interview {
	return &model.Interview{
		ID:     in.ID,
		UserID: in.UserID,
		Status: model.InterviewStatus(in.Status),
		JobInfo: model.JobInfo{
			Position: internal.Position(in.JobPosition),
		},
		QuestionsCount: in.QuestionCount,
		State:          model.InterviewState(in.State),
	}
}
