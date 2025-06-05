FROM alpine:latest

RUN mkdir /app

COPY contactApp /app

CMD ["/app/contactApp"]