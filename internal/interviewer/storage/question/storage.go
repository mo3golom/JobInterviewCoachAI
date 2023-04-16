package question

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"job-interviewer/internal/interviewer/model"
	"job-interviewer/pkg/transactional"
)

type sqlxQuestion struct {
	ID          uuid.UUID `db:"id"`
	Text        string    `db:"text"`
	JobPosition string    `db:"job_position"`
	JobLevel    string    `db:"job_level"`
	IQID        int64     `db:"iq_id"`
}

type sqlxInterviewQuestion struct {
	InterviewID uuid.UUID                     `db:"interview_id"`
	QuestionID  uuid.UUID                     `db:"question_id"`
	Answer      *string                       `db:"answer"`
	GptComment  *string                       `db:"gpt_comment"`
	Status      model.InterviewQuestionStatus `db:"status"`
}

type DefaultStorage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *DefaultStorage {
	return &DefaultStorage{db: db}
}

func (s *DefaultStorage) CreateQuestion(ctx context.Context, tx transactional.Tx, in *model.Question) error {
	query := `
		INSERT 
		INTO question (id, text, job_position, job_level) 
		VALUES (:id, :text, :job_position, :job_level)
		ON CONFLICT DO NOTHING 
    `

	questions := sqlxQuestion{
		ID:          in.ID,
		Text:        in.Text,
		JobPosition: in.JobInfo.Position,
		JobLevel:    "unknown",
	}
	_, err := tx.NamedExecContext(ctx, query, questions)
	return err
}

func (s *DefaultStorage) AttachQuestionToInterview(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID, question *model.Question) error {
	query := `
		INSERT 
		INTO interview_question (interview_id, question_id, status) 
		VALUES (:interview_id, :question_id, :status)
		ON CONFLICT DO NOTHING 
    `

	in := sqlxInterviewQuestion{
		InterviewID: interviewID,
		QuestionID:  question.ID,
		Status:      model.InterviewQuestionStatusActive,
	}
	_, err := tx.NamedExecContext(ctx, query, in)
	return err
}

func (s *DefaultStorage) FindActiveQuestionByInterviewID(ctx context.Context, tx transactional.Tx, interviewID uuid.UUID) (*model.Question, error) {
	query := `
		SELECT q.id, q.text, q.job_position, q.job_level, iq.id as iq_id
		FROM question as q
		JOIN interview_question as iq on q.id = iq.question_id
		WHERE iq.interview_id = $1 AND (iq.status = $2)
		LIMIT 1
    `

	var results []sqlxQuestion
	err := tx.SelectContext(
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

func (s *DefaultStorage) FindAnswersCommentsByInterviewID(ctx context.Context, interviewID uuid.UUID) ([]string, error) {
	query := `
		SELECT iq.gpt_comment
		FROM  interview_question as iq
		WHERE iq.interview_id = $1 and (gpt_comment is not null or gpt_comment != '')
		LIMIT 1
    `

	var results []string
	err := s.db.SelectContext(
		ctx,
		&results,
		query,
		interviewID,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (s *DefaultStorage) SetQuestionAnswered(ctx context.Context, tx transactional.Tx, in SetQuestionAnsweredIn) error {
	query := `
		UPDATE interview_question 
        SET status = :status, answer=:answer, gpt_comment=:gpt_comment, updated_at=now()
        WHERE interview_id = :interview_id AND question_id = :question_id
    `

	answer := in.Answer
	gptComment := in.GptComment
	_, err := tx.NamedExecContext(
		ctx,
		query,
		sqlxInterviewQuestion{
			InterviewID: in.InterviewID,
			QuestionID:  in.QuestionID,
			Answer:      &answer,
			GptComment:  &gptComment,
			Status:      model.InterviewQuestionStatusAnswered,
		},
	)
	return err
}

func (s *DefaultStorage) UpdateInterviewQuestionStatus(ctx context.Context, tx transactional.Tx, in UpdateInterviewQuestionStatusIn) error {
	query := `
		UPDATE interview_question 
        SET status = $1, updated_at=now()
        WHERE interview_id = $2 AND status=$3 
    `

	_, err := tx.ExecContext(
		ctx,
		query,
		in.Target,
		in.InterviewID,
		in.Current,
	)
	return err
}

func convertQuestion(in *sqlxQuestion) *model.Question {
	return &model.Question{
		ID:   in.ID,
		Text: in.Text,
		JobInfo: model.JobInfo{
			Position: in.JobPosition,
		},
	}
}
