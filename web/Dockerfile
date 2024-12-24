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

# Create necessary directories and set permissions
RUN mkdir -p /var/cache/nginx /var/run && \
    chmod -R 777 /var/cache/nginx /var/run && \
    # Remove default nginx static files
    rm -rf /usr/share/nginx/html/* && \
    # Make sure nginx can write to these directories
    chown -R nginx:nginx /var/cache/nginx /var/run /usr/share/nginx/html

# Copy built assets and config
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Switch to non-root user
USER nginx

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
