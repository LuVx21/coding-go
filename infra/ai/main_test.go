package ai

import (
	"fmt"
	"testing"
)

func Test_00(t *testing.T) {
	deepseek := NewServiceProvider("https://api.deepseek.com", "", "deepseek-chat", "deepseek-reasoner")
	siliconflow := NewServiceProvider("https://api.siliconflow.cn", "", "deepseek-ai/DeepSeek-V3", "deepseek-ai/DeepSeek-R1")
	volces := NewServiceProvider("https://ark.cn-beijing.volces.com/api/v3/", "", "deepseek-r1-250120", "deepseek-v3-250324", "deepseek-r1-distill-qwen-7b-250120", "deepseek-r1-distill-qwen-32b-250120")

	fmt.Println(deepseek, siliconflow, volces)
}
