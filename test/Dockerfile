FROM --platform=linux/arm64/v8 grafana/k6:latest
COPY ./ /test
WORKDIR /test
CMD ["run", "uploadchunk.js"]

# docker build --platform linux/arm64/v8 -t playhouse-test .
# - this image can only be used with playhouse project started with docker, we suggest installing k6 if you're testing
#   on your local
# - change --platform for your own platform, or remove the param
#docker run --name playhouse-test -e HOST=playhouse-playhouse-server-1:2345 --network=playhouse_playhouse-net playhouse-test
# - change HOST to the name of the video upload host
# - change --network to the docker network if you have any
