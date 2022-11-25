############################
# STEP 1 build executable binary
############################
ARG ENV
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add  git && apk add  build-base upx

WORKDIR /src/go-graphql-starter/graphql-server/
COPY . .
# Fetch dependencies.
# Using go get.
# RUN cd ./cmd/graphql-server && go build  -mod vendor  -ldflags '-w -s'  -o /go/bin/graphql-server
# Build the binary.
RUN cd ./cmd/graphql-server && go build  -o /go/bin/graphql-server
RUN upx /go/bin/graphql-server

############################
# STEP 2 build a small image
############################
FROM alpine
ARG ENV
# Copy our static executable.
RUN apk update && apk add --no-cache  vips-dev
COPY --from=builder /go/bin/graphql-server /go/bin/graphql-server
COPY --from=builder /src/go-graphql-starter/graphql-server/config/$ENV/graphql-server.yaml /go/bin/
COPY --from=builder /src/go-graphql-starter/graphql-server/pkg/lib/errorMessages.json /etc/
COPY --from=builder /src/go-graphql-starter/graphql-server/emoji_search.bleve /go/bin/emoji_search.bleve
# Run the hello binary.
ENTRYPOINT ["/go/bin/graphql-server",  "-configpath=/go/bin/graphql-server.yaml"]
