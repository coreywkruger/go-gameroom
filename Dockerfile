FROM golang:1.6

RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install build-essential -y

COPY . /go/src/gameroom
WORKDIR /go/src/gameroom

RUN chmod u+x entrypoint.sh
RUN make install && make build

EXPOSE 3334

CMD ["/go/src/gameroom/entrypoint.sh"]
