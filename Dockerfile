# Builder image
FROM golang:1.14.2-alpine3.11 as builder
# Install dependencies. gcc and musl-dev are needed for go-sqlite3 (cgo)
RUN apk --no-cache add git gcc musl-dev

WORKDIR /app

# Copy only module and sum file to keep the layer in cache
# It will only change when the dependencies change
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Copy the rest of the source files
COPY . .

# Build the go app and tar the binary and assets together.
# It avoids having many COPY layers in the runtime image
RUN go build -o main . \
  && tar -cf built.tar main templates/ assets/ database/

# Runtime image
FROM alpine:3.11.6

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Database volume
VOLUME ["/app/database/file"]

# Create a non-root group and user and give the user permissions on /app
RUN addgroup -S mike && adduser -S mike -G mike \
  && chown -R mike:mike /app

USER mike

# Copy the compiled binary and assets (not the sources)
COPY --from=builder --chown=mike:mike ["/app/built.tar",  "/app/"]
# Extract the tar once. It avoids having many COPY layers
RUN tar -xf /app/built.tar \
  && rm /app/built.tar

EXPOSE 8080

CMD ["./main"]

# Check that the homepage is displayed
HEALTHCHECK --interval=5m --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
