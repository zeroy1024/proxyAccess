FROM golang

WORKDIR /app

COPY . .

RUN go build -o ProxyAccess .


CMD [ "ProxyAccess" ]