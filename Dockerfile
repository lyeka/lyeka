FROM klakegg/hugo as builder

WORKDIR /app

COPY . .

RUN ls

# hugo渲染生产站点静态文件
RUN hugo -D

FROM nginx:latest
COPY --from=0 /app/nginx.conf /etc/nginx
COPY --from=0 /app/public  /usr/share/nginx/html