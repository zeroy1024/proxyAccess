FROM golang

WORKDIR /app

COPY . .

RUN go build -o ProxyAccess .
RUN pwd
RUN ls -l

CMD [ "/app/ProxyAccess" ]