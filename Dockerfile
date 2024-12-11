FROM golang:1.23-alpine as builder

WORKDIR /app
COPY . /app

RUN CGO_ENABLED=0 go build -o zhurd ./cmd/api

FROM scratch

COPY --from=builder /app/zhurd .
COPY --from=builder /app/share/config.json.example ./config.json
CMD [ "./zhurd" ]
