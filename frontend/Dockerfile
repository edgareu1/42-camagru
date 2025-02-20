FROM nginx:alpine

# Copy the frontend files
COPY ./files /share/nginx/html
# Copy the NGINX configuration file
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000
