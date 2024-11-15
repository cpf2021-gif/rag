package genai

import (
	"context"
	"fmt"
	"testing"

)

func TestEmbeddingTest(t *testing.T) {
	e := NewEmbedding()
	ctx := context.Background()

    fmt.Println("----- embeddings request -----")
	Text := []string{
		"花椰菜又称菜花、花菜，是一种常见的蔬菜。",
	}

    resp, err := e.CreateEmbeddings(ctx, Text)
    if err != nil {
       t.Errorf("embeddings error: %v\n", err)
    }

	fmt.Println(len(resp.Data[0].Embedding))
}
