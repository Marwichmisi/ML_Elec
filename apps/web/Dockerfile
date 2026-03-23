# Web Dashboard Dockerfile - Multi-arch (linux/amd64, linux/arm64)
# Simplified for MVP

FROM node:20-alpine AS builder

WORKDIR /app

# Install system dependencies
RUN apk add --no-cache git

# Enable corepack for pnpm
RUN corepack enable

# Copy monorepo root files
COPY package.json pnpm-workspace.yaml turbo.json ./

# Copy web and shared packages
COPY apps/web/package.json ./apps/web/
COPY packages/shared/package.json ./packages/shared/
COPY packages/ui/package.json ./packages/ui/ 2>/dev/null || true

# Install dependencies
RUN pnpm install --filter web... --ignore-scripts

# Copy source code
COPY packages/shared ./packages/shared
COPY packages/ui ./packages/ui 2>/dev/null || true
COPY apps/web ./apps/web

# Build shared types first, then web
RUN pnpm --filter shared run build 2>/dev/null || true
RUN pnpm --filter web run build

# Stage 2: Serve with nginx
FROM nginx:alpine AS production

# Copy custom nginx config
COPY apps/web/nginx.conf /etc/nginx/conf.d/default.conf

# Copy built React app from builder
COPY --from=builder /app/apps/web/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:80/ || exit 1

# nginx runs as non-root by default on alpine
CMD ["nginx", "-g", "daemon off;"]
