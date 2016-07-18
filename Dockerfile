FROM golang
ADD .
RUN go get github.com/BurntSushi/toml
RUN go get github.com/asaskevich/govalidator
RUN go get github.com/gorilla/mux
RUN go get github.com/braintree/manners
RUN go install github.com/ziscky/mock-pesa

ENTRYPOINT /go/bin/mock-pesa
EXPOSE 7000