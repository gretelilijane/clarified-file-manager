# Start from golang base image
FROM golang:alpine

# Add Maintainer info
LABEL maintainer="Grete Kuppas"

# Setup folders
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# Build the Go app
RUN go build -o /build

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD [ "/build" ]
