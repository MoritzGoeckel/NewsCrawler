# Static-content-server 
#
# VERSION               0.1

FROM httpd:latest

# RUN ON BUILD
# RUN apt-get update
# RUN apt-get install -y wget

# CMD ON START

LABEL name="static-content-server"

EXPOSE 80/tcp
EXPOSE 80/udp

COPY ./static-content/ /usr/local/apache2/htdocs/

RUN find /usr/local/apache2/htdocs/ -exec chmod 775 {} \;
