FROM golang

ADD . /go/src/github.com/corest/podkiller
WORKDIR /go/src/github.com/corest/podkiller

RUN go get \
    && go test \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o podkiller .


ENTRYPOINT /go/bin/giantswarmdemo


FROM alpine:latest  
RUN apk --no-cache add ca-certificates && \
    mkdir /etc/pod-killer /app

WORKDIR /app/
COPY --from=0 /go/src/github.com/corest/podkiller/podkiller .
COPY --from=0 /go/src/github.com/corest/podkiller/config.toml /etc/pod-killer/
ENTRYPOINT /app/podkiller