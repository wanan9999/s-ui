#!/bin/sh
# One protocol profile for tests, local builds, releases and containers.
BASE_TAGS="with_quic,with_grpc,with_utls,with_acme,with_gvisor,with_tailscale,with_cloudflared,with_openconnect,with_openvpn"
tags_for() {
    case "$1" in
        test) echo "$BASE_TAGS" ;;
        dev|release|docker) echo "$BASE_TAGS,with_naive_outbound,with_purego" ;;
        *) echo "unknown build profile: $1" >&2; return 1 ;;
    esac
}
ldflags_for() { echo "-s -w -buildid="; }
