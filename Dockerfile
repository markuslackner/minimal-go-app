# Start from the official Golang image
FROM golang:1.19.4

# Set the working directory inside the container
WORKDIR /app

RUN mkdir /app/static

# Dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build the Go binary
COPY main.go .
RUN GOOS=linux go build -ldflags '-linkmode=external' -o /app/minimal-go-app

# static resources
COPY static static
ARG BUILD_TIMESTAMP
ARG VERSION
ARG COMMIT_HASH
ARG APP_NAME
ARG WORKFLOW_LINK
ARG REPOSITORY
ARG REPOSITORY_URL

RUN sed -i -e "s/BUILD_TIMESTAMP/$(date '+%Y-%m-%d-%H:%M:%S+%Z')/g" /app/static/openapi.yaml
RUN sed -i -e "s/VERSION/${VERSION}/g" /app/static/openapi.yaml
RUN sed -i -e "s/COMMIT_HASH/${COMMIT_HASH}/g" /app/static/openapi.yaml
RUN sed -i -e "s/APP_NAME/${APP_NAME}/g" /app/static/openapi.yaml
RUN sed -i -e "s/WORKFLOW_LINK/$(echo ${WORKFLOW_LINK} | sed -e 's/\//\\\//g')/g" /app/static/openapi.yaml
RUN sed -i -e "s/REPOSITORY_URL/$(echo ${REPOSITORY_URL} | sed -e 's/^git\(.*\)\.git$/https\1/g' | sed -e 's/\//\\\//g')/g" /app/static/openapi.yaml
RUN sed -i -e "s/REPOSITORY/$(echo ${REPOSITORY} | sed -e 's/\//\\\//g')/g" /app/static/openapi.yaml
RUN echo "templated openapi.yaml:\n$(cat static/openapi.yaml)"

# Start with a minimal Alpine image
FROM alpine:3.17

# Install extra packages
# See https://github.com/gliderlabs/docker-alpine/issues/136#issuecomment-272703023
RUN apk update \
	&& apk add ca-certificates libc6-compat \
	&& update-ca-certificates \
	&& rm -rf /var/cache/apk/*

# Copy the binary from the previous build stage
COPY --from=0 /app/minimal-go-app /usr/local/bin/minimal-go-app
COPY --from=0 /app/static /static

# Expose the port that the application listens on
EXPOSE 8000

# required for external tools to detect this as a go binary
ENV GOTRACEBACK=all


# Set the default command to start the application
CMD ["/usr/local/bin/minimal-go-app"]
