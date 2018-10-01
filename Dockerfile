FROM golang

WORKDIR /go/src/app
COPY . .

RUN go get "github.com/dgrijalva/jwt-go"\
            "github.com/gorilla/mux"\
            "github.com/gorilla/context"\
            "github.com/gorilla/handlers"\
            "gopkg.in/mgo.v2"\
            "github.com/night-codes/mgo-ai"
