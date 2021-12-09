FROM golang:1.17.5-stretch AS builder

RUN apt update && apt upgrade -y && apt install -y make git zlib1g-dev libssl-dev gperf php-cli cmake clang libc++-dev libc++abi-dev

RUN git clone https://github.com/tdlib/td.git

WORKDIR td/build

RUN CXXFLAGS="-stdlib=libc++" CC=/usr/bin/clang CXX=/usr/bin/clang++ cmake -DCMAKE_BUILD_TYPE=Release .. &&\
    cmake --build . --target prepare_cross_compiling &&\
    cd .. &&\
    php SplitSource.php &&\
    cd build && \
    cmake --build . --target install -- -j$(nproc)

WORKDIR /NeoPostAdminBot

COPY ./vendor ./vendor
COPY ./go.mod ./go.sum ./
COPY ./go-tdlib.patch ./go-tdlib.patch

RUN git apply go-tdlib.patch
RUN go build -v ./vendor/...

COPY ./pkg ./pkg
COPY ./cmd ./cmd

RUN go build -v ./cmd/bot

FROM debian:stretch-slim

RUN apt update && apt upgrade -y && apt install -y zlib1g-dev libssl-dev libc++-dev

COPY --from=builder /NeoPostAdminBot/bot ./bot

CMD ["./bot"]
