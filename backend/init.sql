-- modul 1: manajemen pengguna dan otentikasi
-- Tabel untuk menyimpan data utama pengguna
CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Menggunakan UUID lebih baik untuk microservices
  "full_name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) UNIQUE NOT NULL,
  "password_hash" TEXT NOT NULL,
  "profile_picture_url" TEXT,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Tabel untuk peran pengguna, jika nanti dibutuhkan (misal: admin)
CREATE TABLE "roles" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(50) UNIQUE NOT NULL -- Contoh: 'user', 'admin'
);

CREATE TABLE "user_roles" (
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "role_id" INTEGER REFERENCES "roles"("id") ON DELETE CASCADE,
  PRIMARY KEY ("user_id", "role_id")
);

--modul 2: asesmen dan rekomendasi karir
-- Tabel untuk menyimpan semua jalur karier yang tersedia
CREATE TABLE "career_paths" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) UNIQUE NOT NULL, -- Contoh: 'Frontend Developer', 'Data Analyst'
  "description" TEXT NOT NULL,
  "salary_estimation" VARCHAR(255) -- Deskripsi potensi gaji
);

-- Tabel untuk menyimpan pertanyaan-pertanyaan asesmen
CREATE TABLE "assessment_questions" (
  "id" SERIAL PRIMARY KEY,
  "question_text" TEXT NOT NULL,
  "question_type" VARCHAR(50) NOT NULL -- Contoh: 'multiple_choice', 'logical'
);

-- Tabel untuk menyimpan pilihan jawaban dari setiap pertanyaan
CREATE TABLE "assessment_options" (
  "id" SERIAL PRIMARY KEY,
  "question_id" INTEGER REFERENCES "assessment_questions"("id") ON DELETE CASCADE,
  "option_text" TEXT NOT NULL,
  "weight" INTEGER NOT NULL -- Bobot nilai untuk setiap jawaban, digunakan untuk kalkulasi
);

-- Tabel untuk mencatat sesi asesmen yang diambil oleh pengguna
CREATE TABLE "user_assessments" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "status" VARCHAR(50) NOT NULL DEFAULT 'started', -- 'started', 'completed'
  "started_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "completed_at" TIMESTAMPTZ
);

-- Tabel untuk menyimpan setiap jawaban yang dipilih oleh pengguna
CREATE TABLE "user_assessment_answers" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_assessment_id" UUID REFERENCES "user_assessments"("id") ON DELETE CASCADE,
  "question_id" INTEGER REFERENCES "assessment_questions"("id"),
  "option_id" INTEGER REFERENCES "assessment_options"("id")
);

-- Tabel untuk menyimpan hasil akhir rekomendasi untuk pengguna
CREATE TABLE "user_career_recommendations" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_assessment_id" UUID UNIQUE REFERENCES "user_assessments"("id") ON DELETE CASCADE,
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "recommendations" JSONB NOT NULL -- Menyimpan hasil lengkap (misal: [{path_id: 1, score: 95}, {path_id: 2, score: 80}]) [cite: 44]
);

--Modul 3: Kurikulum & Progres Belajar
-- Tabel untuk roadmap belajar yang terikat pada setiap jalur karier
CREATE TABLE "learning_roadmaps" (
  "id" SERIAL PRIMARY KEY,
  "career_path_id" INTEGER UNIQUE REFERENCES "career_paths"("id") ON DELETE CASCADE,
  "name" VARCHAR(255) NOT NULL
);

-- Tabel untuk modul-modul dalam sebuah roadmap
CREATE TABLE "learning_modules" (
  "id" SERIAL PRIMARY KEY,
  "roadmap_id" INTEGER REFERENCES "learning_roadmaps"("id") ON DELETE CASCADE,
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "module_order" INTEGER NOT NULL
);

-- Tabel untuk sumber belajar eksternal/internal di setiap modul
CREATE TABLE "learning_resources" (
  "id" SERIAL PRIMARY KEY,
  "module_id" INTEGER REFERENCES "learning_modules"("id") ON DELETE CASCADE,
  "title" VARCHAR(255) NOT NULL,
  "resource_type" VARCHAR(50) NOT NULL, -- 'article', 'video', 'project'
  "url" TEXT, -- Untuk sumber eksternal
  "content" TEXT, -- Untuk konten internal
  "resource_order" INTEGER NOT NULL
);

-- Tabel untuk melacak progres penyelesaian resource oleh pengguna
CREATE TABLE "user_progress" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "resource_id" INTEGER REFERENCES "learning_resources"("id") ON DELETE CASCADE,
  "status" VARCHAR(50) NOT NULL DEFAULT 'not_started', -- 'in_progress', 'completed'
  "completed_at" TIMESTAMPTZ,
  UNIQUE ("user_id", "resource_id")
);

-- Tabel untuk gamifikasi (badges) 
CREATE TABLE "badges" (
  "id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "icon_url" TEXT
);

-- Tabel penghubung antara pengguna dan badge yang diraih
CREATE TABLE "user_badges" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "badge_id" INTEGER REFERENCES "badges"("id") ON DELETE CASCADE,
  "earned_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE ("user_id", "badge_id")
);

--Modul 4: Portofolio & Proyek
-- Tabel untuk menyimpan detail proyek yang dikerjakan pengguna
CREATE TABLE "user_projects" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "resource_id" INTEGER REFERENCES "learning_resources"("id"), -- Proyek bisa terikat ke resource belajar
  "title" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "project_url" TEXT, -- Link ke Github, Vercel, dll.
  "cover_image_url" TEXT,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

--Modul 5: Lowongan Kerja & Pencocokan
-- Tabel untuk menyimpan data lowongan kerja yang di-scrape
CREATE TABLE "job_postings" (
  "id" BIGSERIAL PRIMARY KEY,
  "title" VARCHAR(255) NOT NULL,
  "company_name" VARCHAR(255),
  "location" VARCHAR(255),
  "description_raw" TEXT, -- Deskripsi HTML/mentah dari sumber
  "description_clean" TEXT, -- Deskripsi bersih setelah parsing
  "source_url" TEXT UNIQUE NOT NULL, -- URL asli lowongan
  "job_type" VARCHAR(100), -- 'internship', 'junior', 'full-time' 
  "skills_required_nlp" JSONB, -- Hasil ekstraksi skill dari NLP [cite: 52]
  "scraped_at" TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- (Opsional) Tabel untuk melacak kecocokan antara user dan lowongan
CREATE TABLE "user_job_matches" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" UUID REFERENCES "users"("id") ON DELETE CASCADE,
  "job_id" BIGINT REFERENCES "job_postings"("id") ON DELETE CASCADE,
  "match_score" INTEGER, -- Skor dari 0-100
  "is_recommended" BOOLEAN DEFAULT false,
  "is_clicked" BOOLEAN DEFAULT false, -- Untuk melacak KPI [cite: 71]
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE ("user_id", "job_id")
);