package embeddings

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) Generate(text string) ([]float32, error) {
	return s.provider.Embed(text)
}
