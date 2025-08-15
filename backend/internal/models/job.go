package models

import "time"

// JobPosting merepresentasikan data di tabel 'job_postings'.
type JobPosting struct {
	ID                int64     `json:"id"`
	Title             string    `json:"title"`
	CompanyName       *string   `json:"company_name,omitempty"`
	Location          *string   `json:"location,omitempty"`
	DescriptionClean  *string   `json:"description,omitempty"`
	SourceURL         string    `json:"source_url"`
	JobType           *string   `json:"job_type,omitempty"`
	ScrapedAt         time.Time `json:"scraped_at"`
}