# Builder image for Golang
FROM golang:1.14.10-alpine3.11 as go-builder
# Install dependencies. gcc and musl-dev are needed for go-sqlite3 (cgo)
RUN apk --no-cache add git gcc musl-dev

WORKDIR /app

# Copy only module and sum file to keep the layer in cache
# It will only change when the dependencies change
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Copy the rest of the source files
COPY . .

# Build the go webserver and cli binaries and tar the webserver and the folders needed at runtime together.
# It avoids having many COPY layers in the runtime image
RUN go build -o . ./cmd/webserver ./cmd/cli \
  && tar -cf built.tar webserver LICENSE templates/ database/

# -----
# Builder image for frontend assets
FROM alpine:3.12.1 as front-builder

# Install nodejs and npm
RUN apk --no-cache add npm

WORKDIR /app

# Copy the source files
COPY . .

# Build the frontend assets and tar the assets together.
# It avoids having to mkdir the assets folder in the runtime image
RUN npm install --no-audit && npm run build \
  && tar -cf assets.tar assets/

# -----
# Runtime image
FROM alpine:3.12.1

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Database volume, TLS cert and key volume, music volume
VOLUME ["/app/database/file/", "/app/secrets/", "/music/"]

# Create a non-root group and user and give the user permissions on /app
RUN addgroup -S mike && adduser -S mike -G mike \
  && chown -R mike:mike /app

# Copy the CLI app so that it is in the PATH
COPY --from=go-builder ["/app/cli", "/usr/local/bin/mike"]

# Change to non-root user
USER mike

# Copy the compiled binary and templates (not the sources)
COPY --from=go-builder --chown=mike:mike ["/app/built.tar",  "./"]
# Copy the minified assets (not the sources)
COPY --from=front-builder --chown=mike:mike ["/app/assets.tar", "./"]
# Extract the compiled binary and assets tarballs. It avoids having many COPY layers
RUN tar -xf ./built.tar && rm ./built.tar \
  && tar -xf ./assets.tar && rm ./assets.tar

EXPOSE 8080 8443

CMD ["./webserver"]

# Check that the homepage is displayed
HEALTHCHECK --interval=5m --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/ || exit 1
