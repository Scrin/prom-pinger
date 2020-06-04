FROM golang:1.14

COPY . /go/src/prom-pinger/
RUN go get prom-pinger/...
RUN go install prom-pinger

ENTRYPOINT ["prom-pinger"]
CMD ["1.1.1.1","1.0.0.1","8.8.8.8","8.8.4.4"]
