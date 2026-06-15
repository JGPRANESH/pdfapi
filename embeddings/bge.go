package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BGEProvider struct {
	BaseURL string
}

func (b *BGEProvider) Embed(text string) ([]float32, error) {

	fmt.Println("===================================")
	fmt.Println("BGEProvider.Embed CALLED")
	fmt.Println("BaseURL =", b.BaseURL)
	fmt.Println("Final URL =", b.BaseURL+"/embed")
	fmt.Println("===================================")

	req := EmbedRequest{
		Text: text,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		b.BaseURL+"/embed",
		"application/json",
		bytes.NewBuffer(body),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"embedding service returned %d",
			resp.StatusCode,
		)
	}

	var result EmbedResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	fmt.Println("Embedding Length:", len(result.Embedding))

	return result.Embedding, nil
}
