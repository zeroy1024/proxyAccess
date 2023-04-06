FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go build -o ProxyAccess .
RUN pwd
RUN ls -l

CMD [ "/app/ProxyAccess" ]