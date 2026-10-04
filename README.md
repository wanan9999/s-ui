# S-UI 中文版

Linux 代理管理面板，基于 sing-box，支持纯 Go **L2TP/IPsec 入站**、多用户分流和统一 DNS。

- 默认简体中文、北京时间（`Asia/Shanghai`）。
- 支持 Linux amd64 / arm64，纯 Go SQLite，`CGO_ENABLED=0` 构建。
- 前端源码在 `frontend/` 本地维护，无需拉取子模块。

[安装](#安装) · [连接 L2TP](#连接-l2tp) · [多用户分流](#多用户分流) · [统一 DNS](#统一-dns) · [开发说明](CONTRIBUTING.md)

## 安装

```bash
bash <(curl -fsSL https://raw.githubusercontent.com/wanan9999/s-ui/main/install.sh)
```

也可下载 [发行包](https://github.com/wanan9999/s-ui/releases)。默认面板端口 **2095**，订阅端口 **2096**；首次登录后修改管理员凭据。

Docker 使用仓库中的 [docker-compose.yml](docker-compose.yml)，采用 Linux host 网络，直接使用宿主机端口。运行 `docker compose up -d`。镜像为 `ghcr.io/wanan9999/s-ui:latest`

## 连接 L2TP

1. 在 **入站管理** 添加类型为 `L2TP/IPsec` 的入站，标签例如 `l2tp-in`。
2. 填写 **服务器公网 IPv4**、**IPsec 预共享密钥（PSK）**。
3. 在 **用户管理** 添加用户，设置 L2TP 用户名、密码并关联该入站；没有启用用户时入站不监听。
4. 按下文配置出站、路由和 DNS。服务器防火墙及云安全组开放 **UDP 500、4500**。
5. 客户端选择 L2TP/IPsec，填写服务器公网 IP、PSK、用户名和密码。

无需开放公网裸 UDP 1701，无需安装 strongSwan、xl2tpd、pppd，也无需 `/dev/net/tun`、系统转发或 NAT。已有 VPN 服务须先停止，避免端口冲突。

## 多用户分流

示例：`alice` 走 `exit-a`，`bob` 走 `exit-b`。请将示例名称替换成自己的账号和标签。

1. 在 **出站管理** 添加两个代理出站，标签分别为 `exit-a`、`exit-b`。
2. 打开 **路由列表 → 添加规则**。
3. 点击 **规则选项**，开启 **入站管理**和 **用户管理**，然后按下表添加两条规则。每条完成后点击弹窗的 **保存**。

| 字段 | 第一条规则 | 第二条规则 |
|---|---|---|
| 入站管理 | `l2tp-in` | `l2tp-in` |
| 用户管理 | `alice` | `bob` |
| 操作 | `Route` | `Route` |
| 出站 | `exit-a` | `exit-b` |

**用户管理**匹配的是认证用户名，L2TP 用户名应与所选用户名称一致；不要填写备注或客户端 IP。当前面板的路由动作仍显示英文 `Route`、`Reject`、`Hijack DNS`。

4. 再添加一条兜底规则：只开启 **入站管理**并选择 `l2tp-in`，**操作**选择 `Reject`，防止没有分流规则的账号使用默认直连。
5. 拖动规则卡片排序：**DNS 劫持 → alice → bob → l2tp-in 的 Reject**，将这组规则放在通用直连规则前面。DNS 劫持按下一节添加。
6. 点击 **路由列表**页面顶部的 **保存**。

## 统一 DNS

**配置一次，所有入站进入核心的普通 DNS 查询统一交给 DoH。** 由服务器直接通过 HTTPS 连接 DNS 服务；解析器看到的是服务器出口 IP，业务按账号走各自出站。

> 新安装的 DNS 列表为空，核心会使用系统 DNS，并非默认防泄漏。请完成以下设置。客户端自带 DoH/DoT、绕过 VPN 的流量和隧道外 IPv6 不受普通 DNS 劫持控制。

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
2. 点击 **保存**。

### 2. 设为全局 DNS

回到 **DNS**页面的 **基础信息**：

- **最终**：选择 `global-doh`。
- **域名解析策略**：选择 `ipv4_only`。

若有已有的 **DNS 规则**，将其中 `路由` 操作的服务器统一改为 `global-doh`；不再需要的规则可删除。保留需要的拒绝规则，但它们匹配的查询不会发送到 DoH。移除规则中自行设置的客户端子网，避免规则覆盖全局设置。点击页面顶部的 **保存**。

### 3. 全局接管普通 DNS

1. 打开 **路由列表 → 添加规则**。
2. 在 **规则选项**开启 **端口**，填写 `53`；开启 **网络**，选择 `tcp` 和 `udp`。
3. **操作**选择 `Hijack DNS`，点击弹窗 **保存**。
4. 将此规则卡片拖到第一位，再点击页面顶部 **保存**。

已有按 `dns` 协议匹配的 `Hijack DNS` 规则可删除，避免重复。

### 4. 确认生效

页面顶部 **保存**会触发核心重新加载，可能断开已有连接；弹窗 **保存**只更新当前页面，不能省略页面顶部的保存。等待核心运行正常后重连客户端。

- 分别连接两个账号，确认业务公网 IP 对应各自出站。
- 测试普通 DNS 解析，确认使用配置的 DoH 服务。
- 测试时临时将 `global-doh` 的 **地址**改为不可达的测试地址，用未缓存的新域名查询：应解析失败。测试后恢复地址并保存。

DNS 检测显示 Cloudflare 解析器地址是正常的，不要求等于代理出口 IP。此方案避免普通 DNS 使用系统解析器和明文上游，但不会隐藏服务器出口 IP；严格防泄漏验收还需抓包。

## 使用边界

- L2TP 当前转发 **IPv4 TCP/UDP**，不转发 ICMP；客户端应使用全局 VPN 并阻断隧道外 IPv6。
- 修改 L2TP 账号会重建对应入站并断开连接。禁用、到期和流量限额沿用用户管理。
- L2TP 客户端须手动配置，不生成 L2TP 出站订阅。

## 来源与许可

基于 [alireza0/s-ui](https://github.com/alireza0/s-ui)；前端基于 [alireza0/s-ui-frontend](https://github.com/alireza0/s-ui-frontend) 的 `f859e16953cd733293618f626cc19b8466e00fd3`，保留原作者版权与 GPL-3.0 许可。

L2TP/IPsec 使用 [wanan9999/veepin](https://github.com/wanan9999/veepin)（MIT，上游为 xen0bit/veepin）。
