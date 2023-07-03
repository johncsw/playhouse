#Build context is ../, which is the server/ directory
FROM golang:latest
WORKDIR /app
RUN mkdir /app/src /app/conf
COPY ./src /app/src
COPY ./conf /app/conf
ENV ENV_PATH=./conf/.env_docker
RUN cd /app/src && go mod download && go build -o main .
CMD ["./src/main"]
