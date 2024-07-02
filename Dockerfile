# Use the official Go image as the base image
FROM golang:1.22 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download and cache Go modules
RUN GOPROXY=direct go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Use a lightweight Alpine Linux image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built executable from the previous build stage
COPY --from=build /app/app .
COPY --from=build /app/.env .

# Expose the application port (adjust if your application uses a different port)
EXPOSE 8016

# Run the application
CMD ["./app"]
