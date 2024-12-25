FROM node:22-alpine AS builder

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Install dependencies
COPY package.json pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy source files and build
COPY . .
RUN pnpm build

# Production image
FROM nginx:alpine

# Create temp directories and set permissions
RUN mkdir -p /tmp/nginx && \
    mkdir -p /var/cache/nginx && \
    mkdir -p /var/log/nginx && \
    mkdir -p /tmp/client_temp /tmp/proxy_temp /tmp/fastcgi_temp /tmp/uwsgi_temp /tmp/scgi_temp && \
    chmod 777 -R /tmp/nginx /var/cache/nginx /var/log/nginx /tmp && \
    rm -rf /usr/share/nginx/html/* 

# Copy nginx config and entrypoint script
COPY nginx.conf /etc/nginx/nginx.conf
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh && \
    chmod 644 /etc/nginx/nginx.conf && \
    chown -R nginx:nginx /etc/nginx /var/cache/nginx /var/log/nginx

# Copy built assets and set permissions
COPY --from=builder /app/dist /usr/share/nginx/html/
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html && \
    chmod 755 /usr/share/nginx/html

# Create a directory for temporary nginx config
RUN mkdir -p /tmp/nginx-temp && \
    chown -R nginx:nginx /tmp/nginx-temp && \
    chmod 755 /tmp/nginx-temp

USER nginx

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
