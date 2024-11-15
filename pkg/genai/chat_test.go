package genai

import (
	"context"
	"fmt"
	"rag/pkg/db"
	"testing"
)

func TestChat(t *testing.T) {
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

	chat := NewChat()	
	embedding := NewEmbedding()


	question := "什么地方隐藏了矮人打造的神器?"
	// question := "你是谁开发的?"
	// question := "龙族的发源地是哪里?"
	ctx := context.Background()
	text := []string{question}

	// 生成问题的embedding
	resp, err := embedding.CreateEmbeddings(ctx, text)
	if err != nil {
		t.Errorf("embeddings error: %v\n", err)
	}
	vector := resp.Data[0].Embedding

	// 搜索最相似的embedding
	db, err := db.NewVectorDB()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
	defer db.Close()

	ids, err := db.Search(vector, 1)
	if err != nil {
		t.Errorf("Error searching embeddings: %v", err)
	}

	// 生成答案
	res, err := chat.GenerateAnswer(question, knowledgeBase[ids[0]-1])
	if err != nil {
		t.Errorf("Error generating answer: %v", err)
	}
	fmt.Println(res)
}
