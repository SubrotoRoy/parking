FROM golang:1.16 as builder

# Add Maintainer Info
LABEL maintainer="Subroto Roy"

WORKDIR /app

# Copy everything from the current directory to container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/SubrotoRoy/parking

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# This container exposes port 8090 to the outside world
EXPOSE 8090

# Run the executable
CMD ["./main"]
