FROM openresty/openresty:alpine

RUN apk add --no-cache perl wget curl

RUN opm get openresty/lua-resty-memcached
RUN mkdir -p /usr/local/openresty/lualib/resty && \
    wget -O /usr/local/openresty/lualib/resty/rabbitmqstomp.lua https://raw.githubusercontent.com/wingify/lua-resty-rabbitmqstomp/master/lib/resty/rabbitmqstomp.lua


COPY ../../config/nginx/mime.types /etc/nginx/mime.types
COPY ../../config/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ../../config/nginx/ssl /etc/nginx/ssl

RUN mkdir -p /var/log/nginx && touch /var/log/nginx/access.log /var/log/nginx/error.log

CMD ["openresty", "-c", "/etc/nginx/nginx.conf", "-g", "daemon off;"]
