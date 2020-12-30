##################################
##
## Management server builder
##
##################################
FROM golang AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download

RUN mkdir bin

# Copy the code into the container
COPY cmd/ ./cmd
COPY pkg/ ./pkg

# Build the application
RUN go build -o bin ./...



##################################
##
## Webpack frontend builder
##
##################################
FROM node:current AS nodebuilder

COPY ./package.json ./package-lock.json ./tsconfig.json ./webpack.config.js ./
RUN npm install
COPY public/ ./public
RUN npx webpack --mode production



##################################
##
## Final minimal combined image
##
##################################
FROM alpine:latest

COPY --from=builder /build/bin/mgmt-server /management-server
COPY --from=nodebuilder dist/ ./dist

# Command to run
ENTRYPOINT ["/management-server"]