package db

import (
	"context"

	milvusdb "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type VectorDB struct {
	db milvusdb.Client
}

const (
	dim                 = 512
	idCol, embeddingCol = "ID", "Embedding"
)

func NewVectorDB() (*VectorDB, error) {
	c, err := milvusdb.NewClient(context.Background(), milvusdb.Config{
		Address: "localhost:19530",
	})
	return &VectorDB{c}, err
}

func (v *VectorDB) Close() {
	_ = v.db.Close()
}

func (v *VectorDB) Check() error {
	collectionName := "doc"
	hasCollection, err := v.db.HasCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}
	if hasCollection {
		return nil
	}

	schema := entity.NewSchema().WithName(collectionName).
		WithField(entity.NewField().WithName(idCol).WithDataType(entity.FieldTypeInt64).WithIsPrimaryKey(true).WithIsAutoID(false)).
		WithField(entity.NewField().WithName(embeddingCol).WithDataType(entity.FieldTypeFloatVector).WithDim(dim))

	err = v.db.CreateCollection(context.Background(), schema, entity.DefaultShardNumber)
	if err != nil {
		return err
	}

	idx, err := entity.NewIndexIvfFlat(entity.L2, 2)
	if err != nil {
		return err
	}
	err = v.db.CreateIndex(context.Background(), collectionName, "Embedding", idx, false)
	return err
}

func (v *VectorDB) Insert(ids []int64, embeddings [][]float32) error {
	idCols := entity.NewColumnInt64(idCol, ids)
	embedingCols := entity.NewColumnFloatVector(embeddingCol, dim, embeddings)

	_, err := v.db.Insert(context.Background(), "doc", "", idCols, embedingCols)
	if err != nil {
		return err
	}

	err = v.db.Flush(context.Background(), "doc", false)
	return err
}

func (v *VectorDB) Search(query []float32, topk int) ([]int64, error) {
	err := v.db.LoadCollection(context.Background(), "doc", false)
	if err != nil {
		return nil, err
	}

	vector := entity.FloatVector(query)
	sp, _ := entity.NewIndexFlatSearchParam()
	sr, err := v.db.Search(context.Background(), "doc", []string{}, "", []string{"ID"}, []entity.Vector{vector}, "Embedding", entity.L2, topk, sp)
	if err != nil {
		return nil, err
	}

	ids := []int64{}
	for _, r := range sr {
		var idColumn *entity.ColumnInt64
		for _, field := range r.Fields {
			if field.Name() == "ID" {
				c, ok := field.(*entity.ColumnInt64)
				if ok {
					idColumn = c
				}
			}
		}
		if idColumn == nil {
			continue
		}
		for i := 0; i < r.ResultCount; i++ {
			id, err := idColumn.ValueByIdx(i)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func (v *VectorDB) State() (entity.LoadState, error) {
	state, err := v.db.GetLoadState(context.Background(), "doc", []string{})
	if err != nil {
		return entity.LoadStateNotExist, err
	}
	return state, nil
}
