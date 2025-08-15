package models

// AssessmentOption merepresentasikan sebuah pilihan jawaban untuk sebuah pertanyaan.
type AssessmentOption struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"-"` // Kita tidak perlu kirim ini ke client
	OptionText string `json:"option_text"`
}

// AssessmentQuestion merepresentasikan sebuah pertanyaan asesmen beserta pilihan jawabannya.
type AssessmentQuestion struct {
	ID           int                `json:"id"`
	QuestionText string             `json:"question_text"`
	QuestionType string             `json:"question_type"`
	Options      []AssessmentOption `json:"options"`
}

type UserAnswer struct {
	QuestionID int `json:"question_id" binding:"required"`
	OptionID   int `json:"option_id" binding:"required"`
}

type SubmitAssessmentInput struct {
	Answers []UserAnswer `json:"answers" binding:"required,dive"`
}

type CareerRecommendation struct {
	CareerPathID      int    `json:"career_path_id"`
	CareerPathName    string `json:"career_path_name"`
	CareerDescription string `json:"career_description"`
	MatchScore        int    `json:"match_score"`
}