FROM nginx:latest

RUN apt-get update && apt-get install -y \
    libnginx-mod-http-memcached

COPY ../../config/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ../../config/nginx/ssl /etc/nginx/ssl

CMD ["nginx", "-g", "daemon off;"]
