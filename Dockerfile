FROM golang:1.22.0
WORKDIR /
COPY . .

RUN go get -u -d gocv.io/x/gocv
RUN cd $GOPATH/src/gocv.io/x/gocv
RUN make install
RUN go install gocv.io/x/gocv


RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

EXPOSE 8080
ENTRYPOINT ["./app"]