# Start with the nginx base image
FROM nginx:latest

# Install envsubst for dynamic configuration
RUN apt-get update && apt-get install -y gettext-base && apt-get clean

# Copy the template configuration file
COPY nginx/nginx.conf.template /etc/nginx/templates/nginx.conf.template

