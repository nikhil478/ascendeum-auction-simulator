# Stage 1: Build
FROM golang:1.25.0-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .

# Create outputs directory in builder stage
RUN mkdir /app/outputs

RUN go build -o simulator ./cmd/simulator

# Stage 2: Distroless
FROM gcr.io/distroless/base:nonroot

# Copy binary and outputs directory from builder
COPY --from=builder /app/simulator /simulator
COPY --from=builder /app/outputs /outputs

WORKDIR /

USER nonroot:nonroot

ENTRYPOINT ["/simulator"]
CMD ["-auctions", "40", "-bidders", "100", "-timeout", "5s", "-out", "/outputs"]
