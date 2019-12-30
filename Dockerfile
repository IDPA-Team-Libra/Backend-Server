# Dockerfile References: https://docs.docker.com/engine/reference/builder/
# Start from the latest golang base image
FROM golang:latest
# Set the Current Working Directory inside the container
WORKDIR /app
COPY main_linux_amd64 ./main
EXPOSE 3440
CMD ["./main"]
