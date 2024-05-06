FROM ubuntu:24.04

RUN apt update && apt install -y libsodium-dev

ENTRYPOINT ["/pigeon"]
COPY pigeon /
