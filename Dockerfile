FROM golang:onbuild as build

FROM alpine
RUN apk --update upgrade && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build /go/bin/app /usr/local/bin/oauth2_proxy
   
ENTRYPOINT ["oauth2_proxy"]
