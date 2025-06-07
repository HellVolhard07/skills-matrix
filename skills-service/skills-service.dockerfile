FROM alpine:latest

RUN mkdir /app

COPY skillsApp /app

CMD ["/app/skillsApp"]