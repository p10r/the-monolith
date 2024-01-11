# Pedro

A little tool to notify you via Telegram that your beloved artists are playing a gig in your
area.

## Local Development

- user [.direnv](https://github.com/direnv/direnv)
- create an `.envrc` file in the project root
- in there, set all the variables listed in `cmd/main.go`
- run `direnv allow .`
- run `go run cmd/main.go`
- a new `local.db` will be created in the project root

## TODOs:

- use lefthook for pre-commit hooks, e.g. running tests
- golang ci lint
- parallel tests
- short/long tests
- move base package
- handle 404 when requesting events
- throw error if artist cant be found
- indicate if there's a space
- do vulnerability checks