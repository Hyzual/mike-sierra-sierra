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

# Create a non-root group and user and give the user permissions on /app
RUN addgroup -S mike && adduser -S mike -G mike \
  && chown mike:mike /app

USER mike

# Copy the compiled binary (not the sources)
COPY --from=builder --chown=mike:mike "/app/main" "/app/"
# Copy the assets
COPY --from=builder --chown=mike:mike "/app/assets" "/app/assets"

EXPOSE 8080

CMD ["./main"]

# Check that the homepage is displayed
HEALTHCHECK --interval=5m --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
