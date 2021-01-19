FROM golang:alpine AS build-env

ARG NETWORK=testnet

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3

# Set working directory for the build
WORKDIR /go/src/github.com/crypto-com/chain-indexing

# Install minimum necessary dependencies
RUN apk add --no-cache $PACKAGES && \
  go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate

# Add source files
COPY . .

# build Cosmos SDK, remove packages
RUN make all

# Final image
FROM alpine:edge

ENV CHAIN_INDEXING /chain-indexing

# Install ca-certificates bash
RUN apk add --update ca-certificates bash gettext curl

RUN addgroup chain-indexing && \
  adduser -S -G chain-indexing chain-indexing -h "$CHAIN_INDEXING"

USER chain-indexing

WORKDIR $CHAIN_INDEXING

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/chain-indexing /go/bin/migrate /usr/bin/

ADD migrations migrations
ADD pgmigrate.sh ./

# Run chain-indexing by default
CMD ["chain-indexing"]