package model

const (
	JobLevelJunior JobLevel = "junior"
	JobLevelMiddle JobLevel = "middle"
	JobLevelSenior JobLevel = "senior"
)

type (
	JobLevel string

	JobInfo struct {
		Position string
		Level    JobLevel
	}
)
