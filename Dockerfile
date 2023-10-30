# syntax=docker/dockerfile:1

## BUILD STAGE
##
ARG GO_VERSION=1.21.1
FROM golang:${GO_VERSION} AS build
WORKDIR /src

## Download dependencies.
  ## Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
RUN --mount=type=cache,target=/go/pkg/mod/ \
  ## Leverage bind mounts to go.sum and go.mod to avoid having to copy them into the container.
  --mount=type=bind,source=go.sum,target=go.sum \
  --mount=type=bind,source=go.mod,target=go.mod \
  go mod download -x

# Build the application.
RUN --mount=type=cache,target=/go/pkg/mod/ \
  --mount=type=bind,target=. \
  CGO_ENABLED=0 go build -o /bin/server .

## PACKAGING STAGE
##
FROM alpine:latest AS final

# Install any runtime dependencies that are needed to run your application.
RUN --mount=type=cache,target=/var/cache/apk \
  apk --update add \
    ca-certificates \
    tzdata \
    && \
    update-ca-certificates

# Create a non-privileged user that the app will run under.
ARG UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  appuser
USER appuser

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]