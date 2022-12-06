FROM golang:alpine AS builder
WORKDIR /build
ADD go.mod .
COPY . .
RUN go build -o tg-feedback-bot main.go
FROM alpine
WORKDIR /build
COPY --from=builder /build/tg-feedback-bot /build/tg-feedback-bot
RUN ls -all
ENTRYPOINT [ "/build/tg-feedback-bot" ]