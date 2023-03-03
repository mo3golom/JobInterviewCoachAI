package startinterview

import "job-interviewer/internal/model"

var (
	levelToString = map[model.JobLevel]string{
		model.JobLevelJunior: "Junior",
		model.JobLevelMiddle: "Middle",
		model.JobLevelSenior: "Senior",
	}
)

type (
	data struct {
	}
)
