FROM golang:1.23-alpine

ENV WORK_DIR="/app"
ENV EXEC_NAME="camagru"
ENV EXEC_PATH="${WORK_DIR}/${EXEC_NAME}"

# Set the working directory
WORKDIR ${WORK_DIR}

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the Go application files
COPY . .

# Build the Go application binary and make it executable
RUN go build -o ${EXEC_PATH} .
RUN chmod +x ${EXEC_PATH}

# Expose the port the application runs on
EXPOSE 8000

# Run the binary
CMD ${EXEC_PATH}
