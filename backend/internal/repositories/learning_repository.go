package repositories

import (
	"database/sql"

	"github.com/zidankhainur2/codecareerix/backend/internal/models"
)

type LearningRepository struct {
	db *sql.DB
}

func NewLearningRepository(db *sql.DB) *LearningRepository {
	return &LearningRepository{db: db}
}

// GetRoadmapByCareerPathID mengambil data roadmap lengkap.
func (r *LearningRepository) GetRoadmapByCareerPathID(careerPathID int) (*models.LearningRoadmap, error) {
	var roadmap models.LearningRoadmap

	// 1. Ambil data roadmap utama
	queryRoadmap := `SELECT id, name FROM learning_roadmaps WHERE career_path_id = $1`
	err := r.db.QueryRow(queryRoadmap, careerPathID).Scan(&roadmap.ID, &roadmap.Name)
	if err != nil {
		return nil, err // Akan mengembalikan sql.ErrNoRows jika tidak ditemukan
	}

	// 2. Ambil semua modul untuk roadmap ini
	queryModules := `
		SELECT id, title, description, module_order 
		FROM learning_modules 
		WHERE roadmap_id = $1 
		ORDER BY module_order`
	rowsModules, err := r.db.Query(queryModules, roadmap.ID)
	if err != nil {
		return nil, err
	}
	defer rowsModules.Close()

	modulesMap := make(map[int]*models.LearningModule)
	for rowsModules.Next() {
		var mod models.LearningModule
		if err := rowsModules.Scan(&mod.ID, &mod.Title, &mod.Description, &mod.ModuleOrder); err != nil {
			return nil, err
		}
		mod.Resources = []models.LearningResource{} // Inisialisasi slice
		modulesMap[mod.ID] = &mod
		roadmap.Modules = append(roadmap.Modules, mod) // Langsung tambahkan ke slice untuk menjaga urutan
	}

	// 3. Ambil semua resource untuk roadmap ini dalam satu query
	queryResources := `
		SELECT r.id, r.module_id, r.title, r.resource_type, r.url, r.content, r.resource_order
		FROM learning_resources r
		JOIN learning_modules m ON r.module_id = m.id
		WHERE m.roadmap_id = $1
		ORDER BY r.resource_order`
	rowsResources, err := r.db.Query(queryResources, roadmap.ID)
	if err != nil {
		return nil, err
	}
	defer rowsResources.Close()

	// 4. Susun resource ke dalam modul yang sesuai
	for rowsResources.Next() {
		var res models.LearningResource
		var moduleID int
		if err := rowsResources.Scan(&res.ID, &moduleID, &res.Title, &res.ResourceType, &res.URL, &res.Content, &res.ResourceOrder); err != nil {
			return nil, err
		}
		if module, found := modulesMap[moduleID]; found {
			module.Resources = append(module.Resources, res)
		}
	}
	
	// Karena kita menggunakan pointer di map, roadmap.Modules sudah terisi dengan resources
	// Kita perlu update slice roadmap.Modules dengan data dari map untuk memastikan data resources-nya masuk
	for i := range roadmap.Modules {
		if mod, found := modulesMap[roadmap.Modules[i].ID]; found {
			roadmap.Modules[i] = *mod
		}
	}


	return &roadmap, nil
}