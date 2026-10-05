# Linux build; Go and SQLite run with CGO_ENABLED=0.
FROM --platform=$BUILDPLATFORM node:26-alpine@sha256:ef24c5053d50fdc3e4e56eb4e7ddb7861874ab0fdc797046ba897581deb8e868 AS front-builder
WORKDIR /app
COPY frontend/ ./
RUN npm ci && npm run build

FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS backend-builder
WORKDIR /app
ARG TARGETARCH
ARG TARGETVARIANT
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=$TARGETARCH

RUN apk upgrade --no-cache --scripts=no apk-tools && \
    apk add --no-cache \
    make \
    git \
    wget \
    unzip \
    bash \
    curl

RUN CRONET_ARCH="$TARGETARCH" && \
    CRONET_URL="https://github.com/SagerNet/cronet-go/releases/latest/download/libcronet-linux-${CRONET_ARCH}.so"; \
    echo "Downloading $CRONET_URL" && \
    wget -q -O ./libcronet.so "$CRONET_URL" && \
    chmod 755 ./libcronet.so

COPY . .
COPY --from=front-builder /app/dist/ /app/web/html/

RUN if [ "$TARGETARCH" = "arm" ]; then export GOARM=7; [ "$TARGETVARIANT" = "v6" ] && export GOARM=6; fi; \
    . ./build-tags.sh && \
    TAGS=$(tags_for docker) && \
    LDFLAGS=$(ldflags_for docker) && \
    go build -trimpath -buildvcs=false -ldflags="$LDFLAGS" -tags "$TAGS" -o sui .

FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b
ENV TZ=Asia/Shanghai
WORKDIR /app
RUN set -ex && apk upgrade --no-cache --scripts=no apk-tools && \
    apk add --no-cache --upgrade bash ca-certificates nftables su-exec && \
    addgroup -S -g 10001 sui && \
    adduser -S -u 10001 -G sui -h /app -s /sbin/nologin sui
COPY --from=backend-builder /app/sui /app/libcronet.so /app/
COPY entrypoint.sh /app/

# Asks the binary, which reads the port the operator actually configured. A
# check with the port written in here goes red the moment they change it in the
# panel, and an orchestrator then kills a container that was working.
HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD ["./sui", "healthcheck"]

# The container still starts as root by default: a TUN inbound needs
# CAP_NET_ADMIN in the process's permitted set, and a panel port below 1024
# needs CAP_NET_BIND_SERVICE. Flipping the default would break both, silently,
# on every existing deployment. Set SUI_UID (and optionally SUI_GID) to have
# entrypoint.sh hand over to an unprivileged user instead -- see entrypoint.sh.
ENTRYPOINT [ "./entrypoint.sh" ]
