package startinterview

import (
	"job-interviewer/internal/interviewer/model"
)

var (
	levelToString = map[model.JobLevel]string{
		model.JobLevelJunior: "Junior",
		model.JobLevelMiddle: "Middle",
		model.JobLevelSenior: "Senior",
	}
)
