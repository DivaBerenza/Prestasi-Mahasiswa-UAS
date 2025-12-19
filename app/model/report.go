package model

type ReportStatistics struct {
	TotalStudents      int `json:"total_students"`
	TotalAchievements  int `json:"total_achievements"`
	SubmittedCount     int `json:"submitted_count"`
	VerifiedCount      int `json:"verified_count"`
	RejectedCount      int `json:"rejected_count"`
}

type ReportStudent struct {
	StudentID      string `json:"student_id"`
	Name           string `json:"name"`
	TotalAchievements int `json:"total_achievements"`
	SubmittedCount int `json:"submitted_count"`
	VerifiedCount  int `json:"verified_count"`
	RejectedCount  int `json:"rejected_count"`
}


