package question

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/model"
	"job-interviewer/pkg/transactional"
)

type sqlxQuestion struct {
	ID          uuid.UUID      `db:"id"`
	Text        string         `db:"text"`
	JobPosition string         `db:"job_position"`
	JobLevel    model.JobLevel `db:"job_level"`
	IQID        int64          `db:"iq_id"`
}

type sqlxInterviewQuestion struct {
	InterviewID uuid.UUID                     `db:"interview_id"`
	QuestionID  uuid.UUID                     `db:"question_id"`
	Status      model.InterviewQuestionStatus `db:"status"`
}

type DefaultStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) CreateQuestions(ctx context.Context, tx transactional.Tx, in []model.Question) error {
	query := `
		INSERT 
		INTO question (id, text, job_position, job_level) 
		VALUES (:id, :text, :job_position, :job_level)
		ON CONFLICT DO NOTHING 
    `

	questions := make([]sqlxQuestion, 0, len(in))
	for _, q := range in {
		questions = append(
			questions,
			sqlxQuestion{
				ID:          q.ID,
				Text:        q.Text,
				JobPosition: q.JobInfo.Position,
				JobLevel:    q.JobInfo.Level,
			},
		)
	}

	_, err := tx.NamedExecContext(ctx, query, questions)
	return err
}

func (s *DefaultStorage) AttachQuestionsToInterview(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, questions []model.Question) error {
	query := `
		INSERT 
		INTO interview_question (interview_id, question_id, status) 
		VALUES (:interview_id, :question_id, :status)
		ON CONFLICT DO NOTHING 
    `

	in := make([]sqlxInterviewQuestion, 0, len(questions))
	for _, q := range questions {
		in = append(
			in,
			sqlxInterviewQuestion{
				InterviewID: interviewID,
				QuestionID:  q.ID,
				Status:      model.InterviewQuestionStatusCreated,
			},
		)
	}

	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}

func (s *DefaultStorage) FindNextQuestion(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID) (*model.Question, error) {
	query := `
		SELECT q.id, q.text, q.job_position, q.job_level, iq.id as iq_id
		FROM question as q
		JOIN interview_question as iq on q.id = iq.question_id
		WHERE iq.interview_id = $1 AND (iq.status = $2 or iq.status = $3)
		LIMIT 1
    `

	var results []sqlxQuestion
	err := s.db.SelectContext(
		ctx,
		&results,
		query,
		interviewID,
		model.InterviewQuestionStatusCreated,
		model.InterviewQuestionStatusActive,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrEmptyQuestionResult
	}

	query = `
        UPDATE interview_question SET status=$1
		WHERE id=$2
    `
	_, err = tx.ExecContext(ctx, query, model.InterviewQuestionStatusActive, results[0].IQID)
	if err != nil {
		return nil, err
	}

	return convertQuestion(&results[0]), nil
}

func (s *DefaultStorage) FindActiveQuestionByInterviewID(ctx context.Context, interviewID uuid.UUID) (*model.Question, error) {
	query := `
		SELECT q.id, q.text, q.job_position, q.job_level, iq.id as iq_id
		FROM question as q
		JOIN interview_question as iq on q.id = iq.question_id
		WHERE iq.interview_id = $1 AND (iq.status = $2)
		LIMIT 1
    `

	var results []sqlxQuestion
	err := s.db.SelectContext(
		ctx,
		&results,
		query,
		interviewID,
		model.InterviewQuestionStatusActive,
	)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrEmptyQuestionResult
	}

	return convertQuestion(&results[0]), nil
}

func (s *DefaultStorage) SetQuestionAnswered(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, questionID uuid.UUID) error {
	query := `
		UPDATE interview_question 
        SET status = :status
        WHERE interview_id = :interview_id AND question_id = :question_id
    `

	_, err := tx.NamedExecContext(
		ctx,
		query,
		sqlxInterviewQuestion{
			InterviewID: interviewID,
			QuestionID:  questionID,
			Status:      model.InterviewQuestionStatusAnswered,
		},
	)
	return err
}

func convertQuestion(in *sqlxQuestion) *model.Question {
	return &model.Question{
		ID:   in.ID,
		Text: in.Text,
		JobInfo: model.JobInfo{
			Position: in.JobPosition,
			Level:    model.JobLevel(in.JobLevel),
		},
	}
}
