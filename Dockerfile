# Use the official Golang image as the base image
FROM golang:1.22.0
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source code into the container
COPY . . 
# Build the Go app
RUN go build -o main ./cmd


# # Use the official Node.js image as the base image
# FROM node:18-alpine

# # Set the working directory inside the container
# WORKDIR /app

# # Copy package.json and package-lock.json files
# COPY package*.json ./

# # Install dependencies
# RUN npm install

# # Copy the source code into the container
# COPY . .

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the application
CMD [ "./main" ]