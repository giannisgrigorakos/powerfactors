FROM golang:1.20.1 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code to the container
COPY . ./

# Build the Go application
RUN go build -o powerfactors ./cmd/powerfactors/

# Set the environment variables
ENV ADDRESS=127.0.0.1
ENV PORT=3000

# Expose the port that the application listens on
EXPOSE $PORT

ENTRYPOINT ["/app/powerfactors"]
CMD ["./powerfactors", "-address={$ADDRESS}", "-port={$PORT}"]