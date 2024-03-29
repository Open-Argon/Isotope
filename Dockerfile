# Use Alpine as the base image
FROM golang:alpine

# Set the Current Working Directory inside the container
WORKDIR /isotope

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY ./src ./src

# Copy License
COPY LICENSE .

# Build the Go app
RUN go build -trimpath -ldflags="-s -w" -o /argon/bin/isotope ./src

# make the binary executable
RUN chmod +x /argon/bin/isotope

# add the binary to the path
ENV PATH="/argon/bin:${PATH}"