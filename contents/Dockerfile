FROM alpine:latest
COPY ./apk-repositories.txt /etc/apk/repositories
RUN apk update \
    && apk add sqlite git openssh-client
ADD ./ /