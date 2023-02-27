package usecase

import (
	"job-interviewer/internal/service/interview"
	"job-interviewer/internal/service/question"
)

type Interviewer struct {
	questionService  question.Service
	interviewService interview.Service
}

func NewInterviewer(q question.Service, i interview.Service) *Interviewer {
	return &Interviewer{
		questionService:  q,
		interviewService: i,
	}
}
