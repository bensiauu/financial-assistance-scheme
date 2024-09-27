# Stage 1: Build the React app
FROM node:18-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install --frozen-lockfile

COPY . .
RUN npm run build

# Stage 2: Just provide the built static files
# The build output will be served by Nginx, so we don't need to expose ports or run anything here.
FROM alpine:latest

WORKDIR /app

# Copy the build output from the first stage
COPY --from=build /app/build /app
