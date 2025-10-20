FROM golang:1.24.5

COPY ./ /go/src/backend-test/
WORKDIR /go/src/backend-test/

# Install dependencies
RUN go mod download

# Build the application
RUN go build -o main .

# Expose port 1112
EXPOSE 1112

# Run the application
CMD ["/go/src/backend-test/main"]