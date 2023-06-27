#Build context is ../, which is the ui/ directory
FROM nginx:alpine
WORKDIR /usr/share/nginx/html
COPY /src .
COPY /conf/nginx.conf /etc/nginx/nginx.conf
