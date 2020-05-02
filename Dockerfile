# Builder image
FROM golang:1.14.2-alpine3.11 as builder

WORKDIR /app

# Copy only module and sum file to keep the layer in cache
# It will only change when the dependencies change
COPY go.mod ./
RUN go mod download

# Copy the rest of the source files
COPY . .

RUN go build -o main .

# Runtime image
FROM alpine:3.11.6

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy only the compiled binary
COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]

# Check that the homepage is displayed
HEALTHCHECK --interval=5m --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
