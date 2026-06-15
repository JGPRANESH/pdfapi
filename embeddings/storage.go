package embeddings

import (
	"encoding/json"
	"os"
)

func LoadEmbeddings(path string) ([]ChunkEmbedding, error) {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []ChunkEmbedding{}, nil
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var embeddings []ChunkEmbedding

	err = json.Unmarshal(data, &embeddings)

	if err != nil {
		return nil, err
	}

	return embeddings, nil
}

func SaveEmbeddings(
	path string,
	embeddings []ChunkEmbedding,
) error {

	data, err := json.MarshalIndent(
		embeddings,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	return os.WriteFile(
		path,
		data,
		0644,
	)
}
