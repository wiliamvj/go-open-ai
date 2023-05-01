# Go OpenAI

### Project with integration to OpenAI's ChatGPT
### Using http and gRPC protocol

Technologies used
- Golang
- MYSQL
- Docker
- SQLC

To count the tokens the [tiktoken-go](https://github.com/j178/tiktoken-go).

To start docker
```
docker-compose up -d
```

Generate sql
```
sqlc generate
```

Run Migrate
```
make migrate
```
