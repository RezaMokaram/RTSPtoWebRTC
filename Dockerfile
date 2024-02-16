FROM ubuntu:22.04

RUN apt-get update && apt-get install -y golang

WORKDIR /
COPY . .

RUN apt install gcc -y
RUN apt-get install git -y
RUN apt install make -y
RUN apt-get install -y ca-certificates && update-ca-certificates 2>/dev/null || true

RUN go get -u -d gocv.io/x/gocv

WORKDIR /root/go/pkg/mod/gocv.io/x/gocv@v0.35.0
RUN sed -i 's/sudo//g' Makefile
RUN make install

WORKDIR /
RUN go build -o app main.go

EXPOSE 8080
ENTRYPOINT ["./app"]