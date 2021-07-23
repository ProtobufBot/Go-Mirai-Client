FROM node:latest AS ui_builder

WORKDIR /build

RUN cd /build \
  && git clone https://github.com/ProtobufBot/pbbot-react-ui.git \
  && cd /build/pbbot-react-ui \
  && npm install \
  && npm run build

FROM golang:1.16-alpine AS gmc_builder

RUN go env -w GO111MODULE=auto \
  && go env -w CGO_ENABLED=0 \
  && go env -w GOPROXY="https://goproxy.io,direct" \
  && mkdir /build

WORKDIR /build

COPY ./ .

COPY --from=ui_builder /build/pbbot-react-ui/build ./pkg/static/static

RUN cd /build \
  && go build -ldflags "-s -w -extldflags '-static'" -o gmc ./service/gmc

FROM alpine:latest

WORKDIR /data

COPY --from=gmc_builder /build/gmc /usr/bin/gmc
RUN chmod +x /usr/bin/gmc

ADD ./scripts/env_run.sh /data/

RUN chmod +x /data/env_run.sh
EXPOSE 9000
CMD /data/env_run.sh