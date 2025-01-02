FROM node:22-slim AS builder

WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy source files and build
COPY . .
RUN pnpm build

# Development stage
FROM node:22-slim AS development

WORKDIR /app

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy package files and install dependencies
COPY package.json pnpm-lock.yaml ./
RUN pnpm install

# Copy source code
COPY . .

# Expose dev server port
EXPOSE 4321

# Default to dev command (will be overridden in docker-compose.dev.yml)
CMD ["pnpm", "dev", "--host"]

# Production stage
FROM nginx:alpine AS production

# Copy built files from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Add nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
