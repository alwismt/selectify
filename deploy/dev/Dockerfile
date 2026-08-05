FROM nginx:1.24

EXPOSE 80
EXPOSE 443

COPY conf/nginx.conf /etc/nginx/nginx.conf
#COPY cert1/selectify.alwis.dev.pem /etc/ssl/selectify.alwis.dev.pem
#COPY cert1/selectify.alwis.dev.key /etc/ssl/selectify.alwis.dev.key

COPY /cert/live/selectify.alwis.dev/cert.pem /etc/ssl/cert.pem
COPY /cert/live/selectify.alwis.dev/privkey.pem /etc/ssl/privkey.pem
COPY /cert/live/selectify.alwis.dev/fullchain.pem /etc/ssl/fullchain.pem


CMD ["nginx", "-g", "daemon off;"]
