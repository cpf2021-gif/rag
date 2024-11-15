package genai

import (
	"context"
	"fmt"
	"os"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

const ragTemplateStr = `
我将向你提出一个问题，并提供一些额外的上下文信息。  
请将这些上下文信息视为事实和正确的内容，它们属于内部文档的一部分。  

如果问题与上下文相关，请结合上下文信息来回答问题。  
如果问题与上下文无关，则正常回答问题, 当你无法回答问题时，请回答"缺乏信息, 无法回答.".

例如，假设上下文中没有关于热带花卉的信息；  
那么如果我问你关于热带花卉的问题，请根据你已知的信息回答，而不要引用上下文内容。  

再比如，如果上下文中提到了矿物学，而我问你相关问题，  
请结合上下文信息和你的常识来提供答案。

问题：  
%s  

上下文：  
%s`

type Chat struct {
	client *arkruntime.Client
}

func NewChat() *Chat {
	return &Chat{
		client: arkruntime.NewClientWithApiKey(
        os.Getenv("ARK_API_KEY"),
        arkruntime.WithBaseUrl("https://ark.cn-beijing.volces.com/api/v3"),
        arkruntime.WithRegion("cn-beijing"),
    	),
	}
}

func (c *Chat) GenerateAnswer(question, doc string) (string, error) {
	ctx := context.Background()
	req := model.ChatCompletionRequest{
		Model: "ep-20240814144928-mjbzp",
		Messages: []*model.ChatCompletionMessage{
		   {
			  Role: model.ChatMessageRoleSystem,
			  Content: &model.ChatCompletionMessageContent{
				 StringValue: volcengine.String("你是智慧助手,是由xypf开发的."),
			  },
		   },
		   {
			  Role: model.ChatMessageRoleUser,
			  Content: &model.ChatCompletionMessageContent{
				 StringValue: volcengine.String(fmt.Sprintf(ragTemplateStr, question, doc)),
			  },
		   },
		},
	 }
 
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}
	return *resp.Choices[0].Message.Content.StringValue, nil
}
