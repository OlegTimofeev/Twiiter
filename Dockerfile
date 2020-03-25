FROM golang:1.14 as builder

# Set the Current Working Directory inside the container

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . /build

WORKDIR /build

# Download all the dependencies
RUN make build

# Run the executable
CMD ["/build/twitter"]

