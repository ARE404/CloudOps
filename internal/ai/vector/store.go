package vector

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"sync"

	"github.com/GoSimplicity/CloudOps/internal/model"
	"go.uber.org/zap"
)

// Document represents a searchable document in the vector store
type Document struct {
	ID        int
	Title     string
	Content   string
	Embedding []float64
}

// Store is an in-memory cosine similarity vector store
type Store struct {
	mu     sync.RWMutex
	docs   []*Document
	logger *zap.Logger
}

func NewStore(logger *zap.Logger) *Store {
	return &Store{logger: logger}
}

// Add inserts or replaces a document
func (s *Store) Add(doc *Document) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.docs {
		if d.ID == doc.ID {
			s.docs[i] = doc
			return
		}
	}
	s.docs = append(s.docs, doc)
}

// Remove deletes a document by ID
func (s *Store) Remove(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, d := range s.docs {
		if d.ID == id {
			s.docs = append(s.docs[:i], s.docs[i+1:]...)
			return
		}
	}
}

// Search returns the top-k most similar documents by cosine similarity
func (s *Store) Search(_ context.Context, query []float64, topK int) []*model.AiKnowledgeChunk {
	s.mu.RLock()
	defer s.mu.RUnlock()

	type scored struct {
		doc   *Document
		score float64
	}

	results := make([]scored, 0, len(s.docs))
	for _, doc := range s.docs {
		sim := cosineSimilarity(query, doc.Embedding)
		results = append(results, scored{doc: doc, score: sim})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if topK > len(results) {
		topK = len(results)
	}

	chunks := make([]*model.AiKnowledgeChunk, 0, topK)
	for _, r := range results[:topK] {
		chunks = append(chunks, &model.AiKnowledgeChunk{
			Title:   r.doc.Title,
			Content: r.doc.Content,
			Score:   r.score,
		})
	}
	return chunks
}

// LoadFromChunks initializes the store from persisted knowledge chunks
func (s *Store) LoadFromChunks(chunks []*model.AiKnowledgeChunk) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs = make([]*Document, 0, len(chunks))
	for _, c := range chunks {
		var emb []float64
		if err := json.Unmarshal([]byte(c.Embedding), &emb); err != nil {
			continue
		}
		s.docs = append(s.docs, &Document{
			ID:        c.ID,
			Title:     c.Title,
			Content:   c.Content,
			Embedding: emb,
		})
	}
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
