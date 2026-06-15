package embeddings

import "fmt"

func GenerateChunkEmbeddings(
	documentID string,
	chunks []string,
	service *Service,
) ([]ChunkEmbedding, error) {

	fmt.Println("GENERATOR STARTED")
	fmt.Println("Total Chunks:", len(chunks))

	var results []ChunkEmbedding

	for i, chunk := range chunks {

		fmt.Printf("Processing Chunk %d\n", i)

		vector, err := service.Generate(chunk)

		if err != nil {
			fmt.Printf("Chunk %d Failed: %v\n", i, err)
			return nil, err
		}

		fmt.Printf("Chunk %d Success\n", i)

		results = append(results, ChunkEmbedding{
			DocumentID: documentID,
			ChunkID:    fmt.Sprintf("chunk_%d", i),
			Text:       chunk,
			Vector:     vector,
		})
	}

	fmt.Println("GENERATOR FINISHED")

	return results, nil
}
