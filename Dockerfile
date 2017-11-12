FROM debian

RUN apt-get update -y && apt-get install ca-certificates -y
   
COPY oauth2_proxy /usr/local/bin/oauth2_proxy

ENTRYPOINT ["oauth2_proxy"]
