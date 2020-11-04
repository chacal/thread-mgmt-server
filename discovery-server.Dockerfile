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

# Copy the code into the container
COPY cmd/ ./cmd

# Build the application
RUN go build -o discovery-server ./cmd/discovery-server/*

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/discovery-server .

# Build a small image
FROM scratch

COPY --from=builder /dist/discovery-server /

# Command to run
ENTRYPOINT ["/discovery-server"]