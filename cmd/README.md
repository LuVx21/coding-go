

### ai: 简单的命令行大语言模型cli客户端

自定义配置服务商及模型, 配置文件优先级依次为:
* `~/config/ai/config/models.json`
* `~/data/ai/config/models.json`
* `models.json`

格式为:
```json
[
  {
    "baseUrl": "https://api.deepseek.com",
    "apiKey": "DEEPSEEK_KEY",
    "modelIds": [
      "deepseek-chat",
      "deepseek-reasoner"
    ]
  }
]
```
* baseUrl:
  * 以`#`结尾, 使用原始地址
  * 以`/`结尾, 忽略v1版本自动添加`/chat/completions`,
  * 其他, 自动添加`/v1/chat/completions`
* apiKey: 配置为真实token的环境变量名
