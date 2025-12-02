FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY frontend/package.json ./
RUN npm install

# Copy source code
COPY frontend/ ./

# Build
RUN npm run build

# Runtime stage with nginx
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
