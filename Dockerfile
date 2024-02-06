FROM alpine:edge AS builder
RUN apk add --no-cache --update go gcc g++

#ENV GOPATH /workdir
COPY . /workdir
WORKDIR /workdir

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-s -w" -trimpath ./cmd/main.go

FROM alpine:latest

# this is here to be able to download the DB
RUN apk --no-cache add \
    openssh-client \
    coreutils

COPY --from=builder /workdir/main /bin/app

CMD ["/bin/app"]