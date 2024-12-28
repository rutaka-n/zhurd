FROM golang:1.23-alpine as builder

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM scratch

COPY --from=builder go/bin/goose bin/goose

ENTRYPOINT ["bin/goose"]
CMD ["up"]
