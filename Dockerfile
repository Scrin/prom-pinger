FROM golang:1.16

WORKDIR /go/src/github.com/Scrin/prom-pinger/
COPY . ./
RUN go install .

ENTRYPOINT ["prom-pinger"]
CMD ["1.1.1.1","1.0.0.1","8.8.8.8","8.8.4.4"]
