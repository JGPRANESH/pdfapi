package embeddings

type EmbedRequest struct {
	Text string `json:"text"`
}

type EmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

type ChunkEmbedding struct {
	DocumentID string    `json:"document_id"`
	ChunkID    string    `json:"chunk_id"`
	Text       string    `json:"text"`
	Vector     []float32 `json:"vector"`
}
