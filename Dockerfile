FROM golang:1.17-alpine AS glc_builder

RUN go env -w GO111MODULE=auto \
  && go env -w CGO_ENABLED=0 \
  && go env -w GOPROXY="https://goproxy.io,direct" \
  && mkdir /build

WORKDIR /build

COPY ./ .

RUN wget https://github.com/2mf8/pbbot-react-ui/releases/latest/download/static.zip && unzip -o static.zip -d ./pkg/static/ \
  && cd /build \
  && go build -ldflags "-s -w -extldflags '-static'" -o glc ./service/glc

FROM alpine:latest

WORKDIR /data

COPY --from=glc_builder /build/glc /usr/bin/glc
RUN chmod +x /usr/bin/glc

ADD ./scripts/env_run.sh /data/

RUN chmod +x /data/env_run.sh
EXPOSE 9000
CMD /data/env_run.sh
