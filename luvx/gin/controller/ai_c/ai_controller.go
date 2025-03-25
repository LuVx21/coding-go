package ai_c

import (
	"encoding/json"
	"fmt"
	"luvx/gin/controller/ai_c/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luvx21/coding-go/coding-common/cast_x"
	"github.com/luvx21/coding-go/coding-common/common_x"
	"github.com/luvx21/coding-go/infra/ai"
)

// HandleChatCompletion 处理聊天补全请求
func HandleChatCompletion(c *gin.Context) {
	interval := cast_x.ToInt64(c.Query("interval"))

	var request ai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 模拟处理逻辑
	if request.Stream {
		// 流式响应处理
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		for i := range len(responses) {
			time.Sleep(time.Duration(common_x.IfThen(interval != 0, interval, 100)) * time.Millisecond)
			response := generateStreamResponse(request, i)
			data, _ := json.Marshal(response)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		}

		fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
		c.Writer.Flush()
	} else {
		// 普通响应处理
		response := generateNormalResponse(request)
		c.JSON(http.StatusOK, response)
	}
}

// generateNormalResponse 生成普通响应
func generateNormalResponse(request ai.ChatCompletionRequest) models.ChatCompletionResponse {
	return models.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", rand.Intn(1000000)),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   request.Model,
		Choices: []models.Choice{
			{
				Index: 0,
				Message: ai.Message{
					Role:    "assistant",
					Content: getAIResponse(request.Messages),
				},
				FinishReason: "stop",
			},
		},
		Usage: models.Usage{
			PromptTokens:     calculateTokens(request.Messages),
			CompletionTokens: 50,
			TotalTokens:      calculateTokens(request.Messages) + 50,
		},
	}
}

var responses = []string{
	"你好！我是 **DeepSeek Chat**，由深度求索（DeepSeek）公司开发的智能AI助手。我可以帮助你解答各种问题，包括学习、工作、编程、写作、翻译、生活建议等。\n",
	"\n",
	"### **我的特点：**\n",
	"✅ **免费使用**：目前无需付费，随时可以向我提问！\n",
	"✅ **知识丰富**：我的知识截止到 **2024年7月**，能提供较新的信息。\n",
	"✅ **超长上下文**：支持 **128K** 上下文，可以处理超长文档和复杂对话。\n",
	"✅ **文件阅读**：可以上传 **PDF、Word、Excel、PPT、TXT** 等文件，帮你提取和分析内容。\n",
	"✅ **编程助手**：支持多种编程语言，能帮你调试代码、优化算法、解释技术概念。\n",
	"✅ **写作与创意**：可以帮你写文章、改写文案、生成故事、头脑风暴创意点子。\n",
	"\n",
	"### **我不能做什么？**\n",
	"❌ 无法实时联网（但你可以手动上传最新资料让我分析）。\n",
	"❌ 不能生成图片、音频或视频（但可以帮你写脚本或描述）。\n",
	"\n",
	"无论是学习、工作还是日常问题，都可以来找我聊聊！😊 你今天想了解什么呢？\n",
	"\n",
	"| a    | b    | c    |\n",
	"| :--- | :--- | :--- |\n",
	"| aa   | bb   | cc   |\n",
}

// generateStreamResponse 生成流格式的响应
func generateStreamResponse(request ai.ChatCompletionRequest, idx int) models.ChatCompletionResponse {

	return models.ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", rand.Intn(1000000)),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   request.Model,
		Choices: []models.Choice{
			{
				Index: 0,
				Message: ai.Message{
					Role:    "assistant",
					Content: responses[idx],
				},
				FinishReason: "stop",
			},
		},
	}
}

// getAIResponse 生成AI回答
func getAIResponse(messages []ai.Message) string {
	if len(messages) == 0 {
		return "Hello! How can I help you today?"
	}
	lastMessage := messages[len(messages)-1]
	return fmt.Sprintf("I received your message saying: '%s'. This is a simulated response from the testing API.", lastMessage.Content)
}

// calculateTokens 简单计算token数
func calculateTokens(messages []ai.Message) int {
	total := 0
	for _, msg := range messages {
		total += len(msg.Content)/4 + 1 // 简单估算
	}
	return total
}

// CORSMiddleware 处理跨域请求
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
