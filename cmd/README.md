

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



### kv: 简单的key-value存储cli工具

支持读取到剪贴板

存储引擎使用 [bbolt](https://github.com/etcd-io/bbolt)

支持类似分库效果(`@xxx`, 不带则使用默认,即xxx为`default`)
```bash
kv set foo@1 bar

kv set 11@1 "测试1"
# 库1中,key=11
kv get 11@1
# 库1所有键值对
kv list @1
```

```bash
# 所有库
kv list -b
# 所有key
kv list -k @1
# 所有value
kv list -v
```