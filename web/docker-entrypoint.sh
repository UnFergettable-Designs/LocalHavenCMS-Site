#!/bin/sh

# Set default values for environment variables
: "${API_HOST:=backend}"
: "${API_PORT:=8090}"
: "${API_PROTOCOL:=https}"

# Replace environment variables in nginx.conf
envsubst '${API_HOST} ${API_PORT} ${API_PROTOCOL}' < /etc/nginx/nginx.conf > /tmp/nginx-temp/nginx.conf.tmp
cp /tmp/nginx-temp/nginx.conf.tmp /etc/nginx/nginx.conf

# Start nginx
exec nginx -g 'daemon off;'