# go-rag api
rag api 项目

## 运行
在当前目录下
```bash
go run main.go
```

## retriever
```bash
curl --request POST \
  --url http://localhost:8000/v1/retriever \
  --header 'Content-Type: application/json' \
  --data '{
    "question":"未来计划",
    "top_k":5,
    "score":0.2
}'
```
