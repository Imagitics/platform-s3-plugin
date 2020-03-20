# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Nikhil Bhide <nikhilbhide15@gmail.com>"

RUN mkdir -p /go/src/github.com/nik/app

ADD . /go/src/github.com/nik/app

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed

# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/nik/app

RUN go mod download

# Go dep!
#RUN go get -u github.com/golang/dep/...
#RUN dep ensure

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]