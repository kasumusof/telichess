# FROM golang

# ADD . /go/src/github.com/kasumusof/telichess

# RUN go install github.com/kasumusof/telichess

# ENTRYPOINT /go/bin/telichess

FROM golang:onbuild

EXPOSE 3000

