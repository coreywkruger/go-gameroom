FROM golang:1.6

RUN curl https://glide.sh/get | sh
RUN mkdir -p /go/src/app
COPY glide.yaml /go/src/app/
COPY glide.lock /go/src/app/
WORKDIR /go/src/app
RUN glide install
COPY . /go/src/app

RUN chmod u+x entrypoint.sh
RUN make install && make build