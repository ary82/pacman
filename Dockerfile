FROM golang:alpine3.20 AS base

WORKDIR /app

COPY . .

RUN apk add make

RUN make clean
RUN go mod download
RUN make build-ssh

FROM scratch
COPY --from=base /app/ssh-bin /app/

WORKDIR /app

CMD [ "./ssh-bin" ]
