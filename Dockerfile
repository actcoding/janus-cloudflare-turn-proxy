FROM alpine

WORKDIR /app

COPY ./dist/jctp jctp

CMD [ "/app/jctp" ]
