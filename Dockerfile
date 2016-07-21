FROM golang
ADD . /go/src/github.com/ziscky/mock-pesa
RUN go get github.com/BurntSushi/toml
RUN go get github.com/asaskevich/govalidator
RUN go get github.com/gorilla/mux
RUN go get github.com/braintree/manners
RUN go get gopkg.in/mgo.v2/bson
RUN go get github.com/icub3d/graceful
RUN go install github.com/ziscky/mock-pesa

ENTRYPOINT /go/bin/mock-pesa
EXPOSE 7000