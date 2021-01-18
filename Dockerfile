FROM golang:alpine AS build-env

ARG NETWORK=testnet

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/crypto-com/chain-indexing

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
  make all

# Final image
FROM alpine:edge

ENV CHAIN_INDEXING /chain-indexing

# Install ca-certificates
RUN apk add --update ca-certificates

RUN addgroup chain-indexing && \
  adduser -S -G chain-indexing chain-indexing -h "$CHAIN_INDEXING"

USER chain-indexing

WORKDIR $CHAIN_INDEXING

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/chain-indexing /usr/bin/chain-indexing

# Run chain-indexing by default
CMD ["chain-indexing"]