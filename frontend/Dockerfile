# Build stage
FROM oven/bun:1 AS builder

WORKDIR /app

# Copy package files
COPY frontend/package.json frontend/bun.lock* ./
COPY frontend/bunfig.toml ./
RUN bun install --frozen-lockfile

# Copy source code
COPY frontend/ ./

# Build the application
ARG VITE_BACKEND_HOST=localhost
ARG VITE_BACKEND_PORT=8080
ARG VITE_BACKEND_SCHEME=http
ARG APP_VERSION=1.0.0

ENV VITE_BACKEND_HOST=${VITE_BACKEND_HOST}
ENV VITE_BACKEND_PORT=${VITE_BACKEND_PORT}
ENV VITE_BACKEND_SCHEME=${VITE_BACKEND_SCHEME}
ENV APP_VERSION=${APP_VERSION}

RUN bun run build

# Production stage with nginx
FROM nginx:alpine

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY docker/etc/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

