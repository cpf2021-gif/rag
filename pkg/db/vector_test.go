package db

import (
	"context"
	"fmt"
	"rag/pkg/genai"
	"testing"
)

var (
	isCreate bool = true
	isInsert bool = true
)

func TestCreateClient(t *testing.T) {
	db, err := NewVectorDB()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	defer db.Close()
}

func TestCreateCollection(t *testing.T) {
	if !isCreate {
		db, err := NewVectorDB()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
		}
		defer db.Close()

		err = db.Check()
		if err != nil {
			t.Errorf("Error creating collection: %v", err)
		}
	}
}

func TestState(t *testing.T) {
	if isCreate {
		db, err := NewVectorDB()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
		}
		defer db.Close()

		state, err := db.State()
		if err != nil {
			t.Errorf("Error getting state: %v", err)
		}
		fmt.Println(state)
	}
}

func TestInsert(t *testing.T) {
	if isCreate && !isInsert {
		knowledgeBase := []string{
			"北风大陆是龙族的发源地，拥有神秘的寒冰之力。",
			"光辉之城是精灵族的家园，以高塔和森林闻名。",
			"灰岩山脉隐藏着矮人打造的无价之锤。",
			"时光图书馆存放着预言书籍，但只有智者能解读。",
			"星辰湖每到夜晚都会映射出宇宙的奥秘。",
			"魔法议会每百年召开一次，讨论世界的平衡。",
			"暗影沼泽中潜伏着不可名状的古老生物。",
			"红宝石帝国的王座镶嵌着传说中的七曜宝石。",
			"银月之刃是唯一能击败黑暗领主的武器。",
			"灵魂之树连接着生者与逝者的世界。",
		}

		e := genai.NewEmbedding()
		ctx := context.Background()

		resp, err := e.CreateEmbeddings(ctx, knowledgeBase)
		if err != nil {
			t.Errorf("embeddings error: %v\n", err)
		}

		ids := make([]int64, len(resp.Data))
		embeddings := make([][]float32, len(resp.Data))

		for i, data := range resp.Data {
			ids[i] = int64(data.Index) + 1
			embeddings[i] = data.Embedding
		}

		db, err := NewVectorDB()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
		}
		defer db.Close()

		err = db.Insert(ids, embeddings)
		if err != nil {
			t.Errorf("Error inserting embeddings: %v", err)
		}
	}
}

func TestSearch(t *testing.T) {
	if isCreate && isInsert {
		knowledgeBase := []string{
			"北风大陆是龙族的发源地，拥有神秘的寒冰之力。",
			"光辉之城是精灵族的家园，以高塔和森林闻名。",
			"灰岩山脉隐藏着矮人打造的无价之锤。",
			"时光图书馆存放着预言书籍，但只有智者能解读。",
			"星辰湖每到夜晚都会映射出宇宙的奥秘。",
			"魔法议会每百年召开一次，讨论世界的平衡。",
			"暗影沼泽中潜伏着不可名状的古老生物。",
			"红宝石帝国的王座镶嵌着传说中的七曜宝石。",
			"银月之刃是唯一能击败黑暗领主的武器。",
			"灵魂之树连接着生者与逝者的世界。",
		}

		text := []string{
			// "龙族的发源地是哪里",
			"时光图书馆",
		}

		e := genai.NewEmbedding()
		ctx := context.Background()

		resp, err := e.CreateEmbeddings(ctx, text)
		if err != nil {
			t.Errorf("embeddings error: %v\n", err)
		}
		vector := resp.Data[0].Embedding

		db, err := NewVectorDB()
		if err != nil {
			t.Errorf("Error creating client: %v", err)
		}
		defer db.Close()

		ids, err := db.Search(vector, 1)
		if err != nil {
			t.Errorf("Error searching embeddings: %v", err)
		}
		fmt.Println(ids)
		if len(ids) == 1 {
			fmt.Println(knowledgeBase[ids[0]-1])
		}
	}
}
