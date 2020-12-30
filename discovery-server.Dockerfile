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

RUN cp /build/bin/discovery-server /discovery-server

# Build a small image
FROM scratch

COPY --from=builder /build/bin/discovery-server /discovery-server

# Command to run
ENTRYPOINT ["/discovery-server"]