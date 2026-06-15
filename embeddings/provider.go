package embeddings

type Provider interface {
	Embed(text string) ([]float32, error)
}
