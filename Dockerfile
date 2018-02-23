FROM nginx:alpine

RUN apk add --no-cache supervisor

COPY ./nginx.conf /etc/nginx/conf.d/default.conf
COPY ./ui/build /var/www/app

COPY ./server/docker-manager-linux /app/server
COPY ./server/config.json /config/config.json

RUN mkdir -p /var/log/supervisor
COPY supervisord.conf /supervisor/supervisord.conf

CMD ["/usr/bin/supervisord", "-c", "/supervisor/supervisord.conf"]