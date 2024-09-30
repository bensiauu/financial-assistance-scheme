# Stage 1: Build the static files
FROM node:18-alpine AS node-builder

WORKDIR /app

COPY ./frontend/package.json ./frontend/package-lock.json ./
RUN npm install --frozen-lockfile

COPY ./frontend/ ./
RUN npm run build

# Stage 2: Nginx to serve the static files
FROM nginx:alpine

# Install bash (optional if you need wait-for-it)
RUN apk add --no-cache bash

# Copy nginx.conf
COPY ./frontend/nginx/nginx.conf /etc/nginx/nginx.conf

# Copy the built static files from node-builder stage
COPY --from=node-builder /app/build /usr/share/nginx/html

# Copy wait-for-it.sh
COPY ./frontend/nginx/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 80

CMD ["/wait-for-it.sh", "app:8080", "--", "nginx", "-g", "daemon off;"]
