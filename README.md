# Pedro

A little tool to notify you via Telegram that your beloved artists are playing a gig in your
area.

## Local Development

- install [.direnv](https://github.com/direnv/direnv)
- install [Lefthook](https://github.com/evilmartians/lefthook)
- run `lefthook install` to set up pre-commit hooks
- copy `.envrc-example` to `.envrc`, adjusting the values
- run `direnv allow .`
- run `go run cmd/main.go`
- a new `local.db` will be created in `local/`

## TODOs:

- use lefthook for pre-commit hooks, e.g. running tests
- golang ci lint
- parallel tests
- short/long tests

- handle 404 when requesting events
- throw error if artist cant be found
- indicate if there's a space
- do vulnerability checks
- Use JSON functionality of sqlite
- Give user info if they don't follow anyone yet