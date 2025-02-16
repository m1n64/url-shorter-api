FROM openresty/openresty:latest

RUN opm get openresty/lua-resty-memcached

COPY ../../config/nginx/nginx.conf /etc/nginx/nginx.conf
COPY ../../config/nginx/ssl /etc/nginx/ssl

CMD ["openresty", "-g", "daemon off;"]
