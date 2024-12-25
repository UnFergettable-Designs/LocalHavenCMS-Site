#!/bin/sh

# Replace environment variables in nginx.conf
envsubst '${API_HOST} ${API_PORT} ${API_PROTOCOL}' < /etc/nginx/nginx.conf > /etc/nginx/nginx.conf.tmp
mv /etc/nginx/nginx.conf.tmp /etc/nginx/nginx.conf

# Start nginx
exec nginx -g 'daemon off;'