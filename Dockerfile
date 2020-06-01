FROM nginx:latest
COPY nginx.conf /etc/nginx
COPY /public  /usr/share/nginx/html