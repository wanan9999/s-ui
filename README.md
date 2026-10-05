# S-UI 中文版

Linux 代理管理面板，基于 sing-box，支持纯 Go **L2TP/IPsec 入站**、多用户分流和统一 DNS。

- 默认简体中文、北京时间（`Asia/Shanghai`）。
- 支持 Linux amd64 / arm64，纯 Go SQLite，`CGO_ENABLED=0` 构建。
- 前端源码在 `frontend/` 本地维护，无需拉取子模块。

[安装](#安装) · [连接 L2TP](#连接-l2tp) · [多用户分流](#多用户分流) · [统一 DNS](#统一-dns) · [开发维护](AGENTS.md)

## 安装

```bash
SUI_LANG=zhcn bash <(curl -fsSL https://raw.githubusercontent.com/wanan9999/s-ui/main/install.sh)
```


Docker 使用仓库中的 [docker-compose.yml](docker-compose.yml)，采用 Linux host 网络，直接使用宿主机端口。运行 `docker compose up -d`。镜像为 `ghcr.io/wanan9999/s-ui:latest`


## 默认安装信息
- 面板端口：2095
- 面板路径：/app/
- 订阅端口：2096
- 订阅路径：/sub/
- 用户名/密码：admin

## 连接 L2TP

1. 在 **入站管理** 添加类型为 `L2TP/IPsec` 的入站，标签例如 `l2tp-in`。
2. 填写 **服务器公网 IPv4**、**IPsec 预共享密钥（PSK）**。
3. 在 **用户管理** 添加用户，设置 L2TP 用户名、密码并关联该入站；没有启用用户时入站不监听。
4. 按下文配置出站、路由和 DNS。服务器防火墙及云安全组开放 **UDP 500、4500**。
5. 客户端选择 L2TP/IPsec，填写服务器公网 IP、PSK、用户名和密码。

无需开放公网裸 UDP 1701，无需安装 strongSwan、xl2tpd、pppd，也无需 `/dev/net/tun`、系统转发或 NAT。已有 VPN 服务须先停止，避免端口冲突。

## 多用户分流

1. 创建多个用户
2. 根据对应的用户名路由到指定出站即可

## 统一 DNS

**所有入站进入核心的普通 DNS 查询统一交给 DoH。** 绕过 VPN 的流量和隧道外 IPv6 不受普通 DNS 劫持控制。

### 1. 添加一个 DNS 服务器

打开 **DNS → 添加 DNS 服务器**，填写：

| 面板字段 | 值 |
|---|---|
| 类型 | `HTTPS` |
| 标签 | `global-doh` |
| 地址 | `1.1.1.1` |
| HTTP 请求路径 | `/dns-query` |

在同一弹窗继续设置：

1. 开启 **启用 TLS**，点击 **TLS 选项**并开启 `SNI`，填写 `cloudflare-dns.com`。
2. 在 **拨号选项**里开启 **转发**，将dns请求转发到出站节点，使用出站代理来处理dns流量。(如果当前服务器为国外，则无需转发)
3. 页面顶部**最终**：选择 `global-doh`，**域名解析策略**：选择 `ipv4_only`，然后点击 **保存**。

### 2. 全局接管普通 DNS

1. 打开 **路由列表 → 添加规则**。
2. 在 **规则选项**开启 **端口**，填写 `53`；开启 **网络**，选择 `tcp` 和 `udp`。
3. **操作**选择 `Hijack DNS`，点击弹窗 **保存**。
4. 将此规则卡片拖到第一位，再点击页面顶部 **保存**。

已有按 `dns` 协议匹配的 `Hijack DNS` 规则可删除，避免重复。


## 使用边界

- 每个 L2TP 会话自动使用 PPP 协商的 MTU，无需在面板手动设置。
- L2TP 当前转发 **IPv4 TCP/UDP**，不转发 ICMP；客户端应使用全局 VPN 并阻断隧道外 IPv6。
- 修改 L2TP 账号会重建对应入站并断开连接。禁用、到期和流量限额沿用用户管理。
- L2TP 客户端须手动配置，不生成 L2TP 出站订阅。
- L2TP 的 IKE/IPsec、PPP 协商日志接入面板日志，排障按时间和对端地址区分客户端；`NO-PROPOSAL-CHOSEN` 不等于密码错误。
- L2TP 支持在保留 PPP 会话的情况下更新 IKE/ESP 密钥；到期、失联或当前 SA 被删除时清理连接。Windows、爱快的长期换钥互通仍需实机验收，不代表其默认算法全部兼容。
- L2TP 监听异常退出会被核心健康检查发现，并由现有守护任务重启核心恢复；恢复时其他入站连接也会中断。

## 来源与许可

基于 [alireza0/s-ui](https://github.com/alireza0/s-ui)；前端基于 [alireza0/s-ui-frontend](https://github.com/alireza0/s-ui-frontend) 的 `f859e16953cd733293618f626cc19b8466e00fd3`，保留原作者版权与 GPL-3.0 许可。

L2TP/IPsec 使用 [wanan9999/veepin](https://github.com/wanan9999/veepin)（MIT，上游为 xen0bit/veepin）。
