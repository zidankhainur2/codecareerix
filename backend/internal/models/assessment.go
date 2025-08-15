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