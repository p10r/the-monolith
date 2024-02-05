FROM golang:1.21.0-bullseye as builder

COPY . /workdir
WORKDIR /workdir

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

RUN go build -ldflags "-s -w" -trimpath ./cmd/main.go

FROM gcr.io/distroless/base-debian11:nonroot
COPY --from=builder /workdir/main /bin/app

CMD ["/bin/app"]