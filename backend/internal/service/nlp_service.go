package service

import (
        "context"
        "strings"

        "github.com/telpp/emagang/internal/models"
        "github.com/telpp/emagang/internal/repository"
)

// NLPService — keyword scoring engine tanpa LLM
// Cocokkan pesan user dengan knowledge base + FAQ menggunakan skor overlap kata kunci
type NLPService struct {
        knowledgeRepo *repository.KnowledgeRepository
}

func NewNLPService(knowledgeRepo *repository.KnowledgeRepository) *NLPService {
        return &NLPService{knowledgeRepo: knowledgeRepo}
}

const nlpThreshold = 0.30 // Skor minimum untuk dianggap "terjawab"

// Query — proses pesan user, kembalikan jawaban atau sinyal tidak terjawab
func (s *NLPService) Query(ctx context.Context, pesan string) models.NLPQueryResponse {
        pesanLower := strings.ToLower(pesan)
        pesanTokens := tokenize(pesanLower)

        // Ambil semua knowledge aktif + FAQ
        knowledgeList, _ := s.knowledgeRepo.GetActive(ctx)
        faqList, _ := s.knowledgeRepo.GetAllFAQAsKnowledge(ctx)

        allKnowledge := append(knowledgeList, faqList...)

        bestSkor := 0.0
        bestJawaban := ""
        bestSumber := ""

        for _, k := range allKnowledge {
                skor := hitungSkor(pesanTokens, pesanLower, k)
                if skor > bestSkor {
                        bestSkor = skor
                        bestJawaban = k.Jawaban
                        if k.Kategori == "faq" {
                                bestSumber = "faq"
                        } else {
                                bestSumber = "knowledge"
                        }
                }
        }

        if bestSkor >= nlpThreshold {
                return models.NLPQueryResponse{
                        Terjawab: true,
                        Jawaban:  bestJawaban,
                        Skor:     bestSkor,
                        Sumber:   bestSumber,
                }
        }

        return models.NLPQueryResponse{
                Terjawab: false,
                Jawaban:  "",
                Skor:     bestSkor,
                Sumber:   "",
        }
}

// hitungSkor — hitung overlap antara pesan user dengan kata kunci knowledge
func hitungSkor(pesanTokens []string, pesanLower string, k models.ChatKnowledge) float64 {
        if len(k.KataKunci) == 0 {
                return 0
        }

        matchCount := 0
        for _, kw := range k.KataKunci {
                kwLower := strings.ToLower(kw)
                // Cek apakah kata kunci ada di pesan (substring atau token match)
                if strings.Contains(pesanLower, kwLower) {
                        matchCount++
                        continue
                }
                // Cek token per token
                kwTokens := tokenize(kwLower)
                for _, kwTok := range kwTokens {
                        if len(kwTok) > 2 && containsToken(pesanTokens, kwTok) {
                                matchCount++
                                break
                        }
                }
        }

        // Skor dari sisi KB: berapa % kata kunci KB yang cocok dengan pesan
        skKB := float64(matchCount) / float64(len(k.KataKunci))

        // Skor dari sisi user: berapa % token user yang ditemukan dalam daftar kata kunci KB
        // Ini mencegah entry dengan banyak keyword mendapat skor rendah padahal user-nya cocok
        allKWLower := strings.ToLower(strings.Join(k.KataKunci, " "))
        userMatchCount := 0
        for _, tok := range pesanTokens {
                if len(tok) > 2 && strings.Contains(allKWLower, tok) {
                        userMatchCount++
                }
        }
        skUser := 0.0
        if len(pesanTokens) > 0 {
                skUser = float64(userMatchCount) / float64(len(pesanTokens))
        }

        // Gunakan nilai tertinggi dari kedua perspektif
        skor := skKB
        if skUser > skor {
                skor = skUser
        }

        // Bonus: jika pertanyaan knowledge mirip dengan pesan user (substring match)
        pertanyaanLower := strings.ToLower(k.Pertanyaan)
        if strings.Contains(pesanLower, pertanyaanLower) || strings.Contains(pertanyaanLower, pesanLower) {
                skor = min(skor+0.3, 1.0)
        }

        return skor
}

func tokenize(text string) []string {
        stopWords := map[string]bool{
                "yang": true, "dan": true, "di": true, "ke": true, "dari": true,
                "untuk": true, "dengan": true, "ini": true, "itu": true, "ada": true,
                "tidak": true, "bisa": true, "apa": true, "bagaimana": true, "kapan": true,
                "berapa": true, "apakah": true, "adalah": true, "akan": true, "sudah": true,
                "atau": true, "saya": true, "kita": true, "kami": true, "mereka": true,
                "ya": true, "yah": true, "dong": true, "deh": true, "nih": true,
                "gimana": true, "gmn": true, "mau": true, "tolong": true,
                "bantu": true, "mohon": true,
        }

        words := splitText(text)
        var tokens []string
        for _, w := range words {
                if len(w) > 2 && !stopWords[w] {
                        tokens = append(tokens, w)
                }
        }
        return tokens
}

func splitText(text string) []string {
        var words []string
        current := strings.Builder{}
        for _, r := range text {
                if isAlphaNum(r) {
                        current.WriteRune(r)
                } else {
                        if current.Len() > 0 {
                                words = append(words, current.String())
                                current.Reset()
                        }
                }
        }
        if current.Len() > 0 {
                words = append(words, current.String())
        }
        return words
}

func isAlphaNum(r rune) bool {
        return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func containsToken(tokens []string, target string) bool {
        for _, t := range tokens {
                if t == target {
                        return true
                }
        }
        return false
}

func min(a, b float64) float64 {
        if a < b {
                return a
        }
        return b
}
