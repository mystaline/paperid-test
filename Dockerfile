# builder stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-w -s" -o /cmd/server .

# final stage using lightweight image
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /cmd/server /server

EXPOSE 8080
ENTRYPOINT ["/server"]