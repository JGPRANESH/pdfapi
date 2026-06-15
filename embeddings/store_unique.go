package embeddings

import "fmt"

func StoreUniqueEmbeddings(
	path string,
	newEmbeddings []ChunkEmbedding,
	threshold float64,
) (int, error) {

	existingEmbeddings, err := LoadEmbeddings(path)

	if err != nil {
		return 0, err
	}

	var uniqueChunks []ChunkEmbedding

	for _, newChunk := range newEmbeddings {

		bestScore := 0.0

		for _, oldChunk := range existingEmbeddings {

			score := CosineSimilarity(
				newChunk.Vector,
				oldChunk.Vector,
			)

			if score > bestScore {
				bestScore = score
			}
		}

		fmt.Printf(
			"Chunk %s Similarity: %.4f\n",
			newChunk.ChunkID,
			bestScore,
		)

		if bestScore < threshold {
			uniqueChunks = append(
				uniqueChunks,
				newChunk,
			)
		}
	}

	allEmbeddings := append(
		existingEmbeddings,
		uniqueChunks...,
	)

	err = SaveEmbeddings(
		path,
		allEmbeddings,
	)

	if err != nil {
		return 0, err
	}

	return len(uniqueChunks), nil
}
