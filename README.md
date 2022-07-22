# gtp-backend


```bash
curl -H "Content-Type: application/json" -d '{ "query": "mutation { createTodo(input: { text: \"item1\", userId: \"user1\" }) { id } }" }' http://localhost:8080/query
```