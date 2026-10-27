## Local Development

- install [direnv](https://github.com/direnv/direnv)
- install [Lefthook](https://github.com/evilmartians/lefthook)
- install [nilaway](https://github.com/uber-go/nilaway?tab=readme-ov-file#standalone-checker)
- run `lefthook install` to set up pre-commit hooks
- copy `.envrc-example` to `.envrc`, adjusting the values
- run `direnv allow .`
- run `go run cmd/main.go`

## Deployment

```
docker build -t monolith .
# to inspect:
docker run -it monolith

docker run monolith

fly launch
fly deploy
fly logs
```

## TODOs

- [move away from Docker](https://fly.io/docs/languages-and-frameworks/golang/)
