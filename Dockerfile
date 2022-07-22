FROM golang:1.17.3 as builder

# install cgo-related dependencies
RUN set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		libvirt-dev \
		gcc \
	; \
	rm -rf /var/lib/apt/lists/*

WORKDIR  /usr/local/go/src/github.com/hamster-shared/hamster-provider/

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

COPY . .

RUN set -eux; \
    go mod tidy ; \
    go build


FROM docker:20

RUN set -eux; \
    apk update ; \
    apk add libvirt-dev

WORKDIR /root/

## test evn
ENV CHAIN_ADDRESS 183.66.65.207:49944
ENV CPU 1
ENV MEMORY 1

COPY --from=builder /usr/local/go/src/github.com/hamster-shared/hamster-provider/hamster-provider /usr/local/bin/
ADD ./frontend/dist frontend/dist

RUN set -eux ;\
    hamster-provider init


CMD sed -i 's/"none"/"done"/' ~/.hamster-provider/config \
    && sed -i "s/127.0.0.1:9944/$CHAIN_ADDRESS/" ~/.hamster-provider/config \
    && sed -i "s/cpu\"\:1/cpu\"\:$CPU/"  ~/.hamster-provider/config\
    && sed -i "s/mem\"\:1/mem\"\:$MEMORY/"  ~/.hamster-provider/config\
    && hamster-provider daemon
