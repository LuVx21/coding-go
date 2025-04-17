```bash
grpcurl --plaintext -d '{"name": "张三"}' localhost:18888 proto.Greeter/SayHello
```