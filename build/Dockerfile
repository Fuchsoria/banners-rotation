FROM golang:1.16 as build

ENV BIN_FILE /opt/banners-rotation/banners-rotation
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ${CODE_DIR}/cmd/banners-rotation/

FROM alpine:3.9

LABEL ORGANIZATION="Private"
LABEL SERVICE="banners-rotation"
LABEL MAINTAINERS="fuchsoria@pm.me"

ENV BIN_FILE "/opt/banners-rotation/banners-rotation"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_FILE_NAME

ENV CONFIG_FILE /etc/banners-rotation/${CONFIG_FILE_NAME}.json
COPY ./configs/${CONFIG_FILE_NAME}.json ${CONFIG_FILE}

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

CMD /wait && ${BIN_FILE} -config ${CONFIG_FILE}
