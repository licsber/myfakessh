FROM golang:alpine AS builder
WORKDIR /licsber
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -ldflags="-w -s" .

FROM scratch
COPY --from=builder /licsber/myfakessh /licsber/myfakessh
EXPOSE 22
ENTRYPOINT ["/licsber/myfakessh"]
