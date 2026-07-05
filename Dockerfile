# Stage 1: Create a minimal runtime image
FROM alpine:latest

# Import variables
ARG BUILD_TIMESTAMP
ARG VERSION
ARG APP_BINARY

# Labels for date, timestamp, and version
LABEL org.opencontainers.image.created=${BUILD_TIMESTAMP}
LABEL org.opencontainers.image.version=${VERSION}

# Set opt workdir
WORKDIR /opt/
COPY ${APP_BINARY} .

USER 1000
CMD ["./${APP_BINARY}"]
