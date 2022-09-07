FROM node:17 as node-builder
WORKDIR  /root
COPY frontend .
RUN npm install && npm run build


FROM golang:1.17.3 as go-builder

WORKDIR  /root

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

COPY . .

RUN set -eux; \
    go mod tidy ; \
    go build


FROM docker

WORKDIR /root

## test evn
ENV CHAIN_ADDRESS 183.66.65.207:49944
ENV CPU 1
ENV MEMORY 1

COPY --from=node-builder /root/dist frontend/dist
COPY --from=go-builder /root/hamster-provider .
COPY ./templates ./templates/

CMD  ./hamster-provider
