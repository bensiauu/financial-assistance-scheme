FROM nginx:alpine

# Install bash
RUN apk add --no-cache bash

# Copy nginx.conf
COPY ./frontend/nginx/nginx.conf /etc/nginx/nginx.conf

# Copy wait-for-it.sh
COPY ./frontend/nginx/wait-for-it.sh /wait-for-it.sh

# Make the wait-for-it.sh script executable
RUN chmod +x /wait-for-it.sh

