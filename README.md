# S-UI 中文版

基于 sing-box 的 Linux 代理管理面板，增加纯 Go L2TP/IPsec 入站。
默认简体中文、`Asia/Shanghai` 时区；前端源码在 `frontend/` 本地维护。
后端 SQLite 使用纯 Go 驱动，构建使用 `CGO_ENABLED=0`。

## 安装

支持 Linux amd64、arm64。安装脚本及升级包来自本仓库：

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/wanan9999/s-ui/main/install.sh)
```

Docker：下载本仓库的 `docker-compose.yml` 后运行 `docker compose up -d`。
镜像为 `ghcr.io/wanan9999/s-ui:latest`，由本仓库发布流水线生成。
尚未发布镜像时，可在源码目录运行 `docker build -t ghcr.io/wanan9999/s-ui:latest .`。
默认面板端口 `2095`，订阅端口 `2096`；登录后立即设置管理员凭据。

## L2TP/IPsec 入站

1. 添加 `L2TP/IPsec` 入站，填写服务器公网 IPv4、PSK 和私网地址池。
2. 在用户页设置 L2TP 账号密码，绑定该入站。没有启用用户时不监听。
3. 设置 sing-box DNS 服务器、出站与路由；按 `auth_user` 为账号指定出口。
   默认配置保留 DNS 劫持规则，但需自行配置可用的 DNS 上游。
4. 服务器防火墙与云安全组开放 UDP **500、4500**，客户端填写服务器、PSK、账号与密码。
   不映射公网裸 UDP 1701，不需要安装 strongSwan、xl2tpd 或 pppd。

数据路径：veepin IKEv1/IPsec → L2TP/PPP → 每会话独立 gVisor 栈 → sing-box DNS/路由/出站。
服务端无需 `/dev/net/tun`、系统转发或 NAT。认证账号随连接传入路由和流量统计；
修改账号会重建入站，当前连接随之断开。禁用、到期和流量限额沿用面板用户管理。

当前支持 **IPv4 TCP/UDP**；ICMP 等其他 IP 协议拒绝转发。原生 L2TP 配置须手动填写，
不生成 sing-box/Clash 不支持的 L2TP 出站订阅。客户端隧道外 IPv6 不在本入站控制范围。
公网手机、Windows、长时间重连及 rekey 仍须按实际部署验收；自动测试不等于公网验收。

## 开发与验证

Go 版本见 `go.mod`，Node.js 26。Linux 构建：

```bash
sh build.sh
. ./build-tags.sh
CGO_ENABLED=0 go test -tags "$(tags_for test)" ./...
CGO_ENABLED=0 go vet -tags "$(tags_for test)" ./...
cd frontend && npm ci && npm run lint && npm test && npm run build
```

Naive 出站通过 purego 加载 `libcronet.so`；发布包和镜像附带该库。
因此 `CGO_ENABLED=0` 不表示所有可选协议都没有动态库依赖。
Linux CI 验证完整 L2TP 拨号、账号分流、DNS 和踢下线，运行在独立网络命名空间。
推送 `v*` 标签发布 Linux 包和 GHCR 镜像；普通主分支推送只测试、构建。

## 来源与许可

本项目基于 [alireza0/s-ui](https://github.com/alireza0/s-ui)，前端基于
[alireza0/s-ui-frontend](https://github.com/alireza0/s-ui-frontend) 的
`f859e16953cd733293618f626cc19b8466e00fd3`。保留原作者版权与 GPL-3.0 许可。
L2TP/IPsec 使用 [wanan9999/veepin](https://github.com/wanan9999/veepin)（MIT），
其上游为 xen0bit/veepin。
