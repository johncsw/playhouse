#Build context is ../, which is the server/ directory
FROM golang:1.20.4
# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*
WORKDIR /app
RUN mkdir /app/src /app/conf /app/chunks
COPY ./src /app/src
COPY ./conf /app/conf
ENV ENV_PATH=./conf/.env_docker
RUN cd /app/src && go mod download && go build -o main .
CMD ["./src/main"]