package models

// LearningResource merepresentasikan satu sumber belajar di dalam sebuah modul.
type LearningResource struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	ResourceType  string  `json:"resource_type"` // 'article', 'video', 'project'
	URL           *string `json:"url,omitempty"`   // Pointer agar bisa NULL, omitempty agar tidak muncul di JSON jika nil
	Content       *string `json:"content,omitempty"`
	ResourceOrder int     `json:"resource_order"`
}

// LearningModule merepresentasikan satu modul di dalam sebuah roadmap.
type LearningModule struct {
	ID          int                `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	ModuleOrder int                `json:"module_order"`
	Resources   []LearningResource `json:"resources"`
}

// LearningRoadmap adalah struktur data lengkap untuk sebuah jalur belajar.
type LearningRoadmap struct {
	ID     int              `json:"id"`
	Name   string           `json:"name"`
	Modules []LearningModule `json:"modules"`
}