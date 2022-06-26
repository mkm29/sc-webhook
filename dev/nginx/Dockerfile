FROM nginx:1.23.0

## install python3
RUN apt-get update && \
  apt-get install -y --no-install-recommends python3

COPY ./nginx.conf /etc/nginx/nginx.conf
COPY ./default.conf /etc/nginx/conf.d/default.conf

# set working directory
WORKDIR /usr/share/nginx/html

# add permissions for nginx user
RUN chown -R nginx:nginx /var/cache/nginx && \
  chown -R nginx:nginx /var/log/nginx && \
  chown -R nginx:nginx /etc/nginx/conf.d
RUN touch /var/run/nginx.pid && \
  chown -R nginx:nginx /var/run/nginx.pid

## switch to non-root user
USER 101