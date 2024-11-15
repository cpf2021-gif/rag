package genai

import (
	"context"
	"math"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

type Embedding struct {
    client *arkruntime.Client
}

func NewEmbedding() *Embedding {
    return &Embedding{
        client: arkruntime.NewClientWithApiKey(
            os.Getenv("ARK_API_KEY"),
            arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
            arkruntime.WithRegion("cn-beijing"),
        ),
    }
}

func (e *Embedding) CreateEmbeddings(ctx context.Context, texts []string) (model.EmbeddingResponse, error) {
    req := model.EmbeddingRequestStrings{
        Input: texts,
        Model: "ep-20241115154400-d96cj",
        Dimensions: 512,
    }

    resp, err := e.client.CreateEmbeddings(ctx, req)
    if err != nil {
        return model.EmbeddingResponse{}, err
    } else {
        for i := range resp.Data {
            resp.Data[i].Embedding = normalizeVector(resp.Data[i].Embedding, 512)
        }
    }

    return resp, nil
}

func normalizeVector(vector []float32, dim int) []float32 {
    vector = vector[:dim]

    var sum float32
    for _, v := range vector {
       sum += v * v
    }
    sum = float32(math.Sqrt(float64(sum)))

    if sum == 0 {
       return vector
    }

    var newVector = make([]float32, dim)
    for i, v := range vector {
        newVector[i] = v / sum
    }
    return newVector
}
