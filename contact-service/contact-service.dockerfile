FROM alpine:latest

RUN mkdir /app

COPY contactApp /app
COPY templates /templates

CMD ["/app/contactApp"]