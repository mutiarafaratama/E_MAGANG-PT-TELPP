package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/telpp/emagang/internal/models"
)

type KnowledgeRepository struct {
	db *pgxpool.Pool
}

func NewKnowledgeRepository(db *pgxpool.Pool) *KnowledgeRepository {
	return &KnowledgeRepository{db: db}
}

func (r *KnowledgeRepository) GetAll(ctx context.Context) ([]models.ChatKnowledge, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, pertanyaan, kata_kunci, jawaban, kategori, is_active, created_at, updated_at
		 FROM chat_knowledge ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ChatKnowledge
	for rows.Next() {
		var k models.ChatKnowledge
		if err := rows.Scan(&k.ID, &k.Pertanyaan, &k.KataKunci, &k.Jawaban, &k.Kategori, &k.IsActive, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, k)
	}
	return list, nil
}

func (r *KnowledgeRepository) GetActive(ctx context.Context) ([]models.ChatKnowledge, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, pertanyaan, kata_kunci, jawaban, kategori, is_active, created_at, updated_at
		 FROM chat_knowledge WHERE is_active = true ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ChatKnowledge
	for rows.Next() {
		var k models.ChatKnowledge
		if err := rows.Scan(&k.ID, &k.Pertanyaan, &k.KataKunci, &k.Jawaban, &k.Kategori, &k.IsActive, &k.CreatedAt, &k.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, k)
	}
	return list, nil
}

func (r *KnowledgeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ChatKnowledge, error) {
	k := &models.ChatKnowledge{}
	err := r.db.QueryRow(ctx,
		`SELECT id, pertanyaan, kata_kunci, jawaban, kategori, is_active, created_at, updated_at
		 FROM chat_knowledge WHERE id=$1`, id).
		Scan(&k.ID, &k.Pertanyaan, &k.KataKunci, &k.Jawaban, &k.Kategori, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (r *KnowledgeRepository) Create(ctx context.Context, req models.ChatKnowledgeRequest) (*models.ChatKnowledge, error) {
	k := &models.ChatKnowledge{}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	kategori := req.Kategori
	if kategori == "" {
		kategori = "umum"
	}
	err := r.db.QueryRow(ctx,
		`INSERT INTO chat_knowledge (pertanyaan, kata_kunci, jawaban, kategori, is_active)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, pertanyaan, kata_kunci, jawaban, kategori, is_active, created_at, updated_at`,
		req.Pertanyaan, req.KataKunci, req.Jawaban, kategori, isActive).
		Scan(&k.ID, &k.Pertanyaan, &k.KataKunci, &k.Jawaban, &k.Kategori, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
	return k, err
}

func (r *KnowledgeRepository) Update(ctx context.Context, id uuid.UUID, req models.ChatKnowledgeRequest) (*models.ChatKnowledge, error) {
	k := &models.ChatKnowledge{}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	kategori := req.Kategori
	if kategori == "" {
		kategori = "umum"
	}
	err := r.db.QueryRow(ctx,
		`UPDATE chat_knowledge SET pertanyaan=$1, kata_kunci=$2, jawaban=$3, kategori=$4, is_active=$5, updated_at=$6
		 WHERE id=$7
		 RETURNING id, pertanyaan, kata_kunci, jawaban, kategori, is_active, created_at, updated_at`,
		req.Pertanyaan, req.KataKunci, req.Jawaban, kategori, isActive, time.Now(), id).
		Scan(&k.ID, &k.Pertanyaan, &k.KataKunci, &k.Jawaban, &k.Kategori, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
	return k, err
}

func (r *KnowledgeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM chat_knowledge WHERE id=$1`, id)
	return err
}

func (r *KnowledgeRepository) ToggleActive(ctx context.Context, id uuid.UUID) (*models.ChatKnowledge, error) {
	k := &models.ChatKnowledge{}
	err := r.db.QueryRow(ctx,
		`UPDATE chat_knowledge SET is_active = NOT is_active, updated_at=$1
		 WHERE id=$2
		 RETURNING id, pertanyaan, kata_kunci, jawaban, kategori, is_active, created_at, updated_at`,
		time.Now(), id).
		Scan(&k.ID, &k.Pertanyaan, &k.KataKunci, &k.Jawaban, &k.Kategori, &k.IsActive, &k.CreatedAt, &k.UpdatedAt)
	return k, err
}

// GetAllFAQAsKnowledge — ambil FAQ aktif dan konversi ke format knowledge untuk NLP
func (r *KnowledgeRepository) GetAllFAQAsKnowledge(ctx context.Context) ([]models.ChatKnowledge, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, pertanyaan, jawaban FROM faq WHERE is_active = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ChatKnowledge
	for rows.Next() {
		var k models.ChatKnowledge
		var pertanyaan, jawaban string
		if err := rows.Scan(&k.ID, &pertanyaan, &jawaban); err != nil {
			return nil, err
		}
		k.Pertanyaan = pertanyaan
		k.Jawaban = jawaban
		k.Kategori = "faq"
		k.IsActive = true
		// Buat kata kunci sederhana dari kata-kata dalam pertanyaan
		k.KataKunci = extractKeywords(pertanyaan)
		list = append(list, k)
	}
	return list, nil
}

// extractKeywords — ekstrak kata kunci dari teks (untuk FAQ yang tidak punya kata_kunci sendiri)
func extractKeywords(text string) []string {
	stopWords := map[string]bool{
		"yang": true, "dan": true, "di": true, "ke": true, "dari": true,
		"untuk": true, "dengan": true, "ini": true, "itu": true, "ada":  true,
		"tidak": true, "bisa": true, "apa": true, "bagaimana": true, "kapan": true,
		"berapa": true, "apakah": true, "adalah": true, "akan": true, "sudah": true,
		"atau": true, "saya": true, "kita": true, "kami": true, "mereka": true,
	}

	words := splitWords(text)
	var keywords []string
	seen := map[string]bool{}
	for _, w := range words {
		if len(w) > 2 && !stopWords[w] && !seen[w] {
			keywords = append(keywords, w)
			seen[w] = true
		}
	}
	return keywords
}

func splitWords(text string) []string {
	var words []string
	current := ""
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			current += string(r)
		} else {
			if len(current) > 0 {
				words = append(words, toLower(current))
				current = ""
			}
		}
	}
	if len(current) > 0 {
		words = append(words, toLower(current))
	}
	return words
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}
