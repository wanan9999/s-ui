# 开发与维护

用户安装和 Web 操作见 [README.md](README.md)。前端在 `frontend/` 直接维护，不使用子模块，保留原始许可证。

## 构建

Go 版本见 `go.mod`，Node.js 26。仅发布 Linux amd64 / arm64，发行构建必须设置 `CGO_ENABLED=0`。

```bash
sh build.sh
```

协议构建标签统一由 `build-tags.sh` 管理。Naive 出站通过 purego 加载 `libcronet.so`，发布包和镜像附带该库；`CGO_ENABLED=0` 不代表所有可选协议都没有动态库依赖。

## 实现约定

- L2TP 通过固定版本 veepin 库处理 IKE、IPsec、L2TP、PPP，不复制协议实现。
- 数据路径：veepin → 每个认证会话独立的 gVisor 栈 → sing-box DNS、路由和出站。
- 进入 sing-box 必须携带认证账号；禁止增加 host NAT、默认路由或未经认证的数据转发。
- 不增加系统 DNS 回退。当前上游核心在 DNS 服务器列表为空时仍隐式使用系统 DNS，不得把空配置宣传为防泄漏。
- 账号变更目前重建入站，不能只修改认证表而保留旧连接。
- 功能变更覆盖配置、前端、数据库、生命周期和失败清理；文档中的按钮、字段、选项以组件及 `frontend/src/locales/zhcn.ts` 为准，硬编码英文选项保留实际显示名称。

## 验证

在 Linux 执行：

```bash
. ./build-tags.sh
CGO_ENABLED=0 go test -tags "$(tags_for test)" ./...
CGO_ENABLED=0 go vet -tags "$(tags_for test)" ./...
(cd frontend && npm ci && npm run lint && npm test && npm run build)
```

L2TP 真实拨号集成测试还需 `SUI_L2TP_INTEGRATION=1`、`/dev/net/tun` 及独立网络命名空间，按 `.github/workflows/test.yml` 的隔离步骤运行。不得直接在开发主机默认网络中执行。

单元测试、隔离网络集成测试和公网验收分别报告。公网独立客户端、同 NAT 并发、NAT 重绑定、长期重连、rekey，以及 DNS 出口故障不直连均须单独验收，不能用配置保存成功代替。

## veepin 依赖更新

生产构建使用 `go.mod` 固定的远端 veepin 版本，并由 `go.sum` 校验。
不使用本地路径 `replace` 或相邻仓库。更新时先推送 veepin 提交，再用
`go get github.com/wanan9999/veepin@<提交 SHA>` 更新依赖，执行 `go mod tidy`，
确认 `GOWORK=off` 下可独立构建，并重跑 Linux 测试及独立客户端验收。

协议换钥、SA 寿命、探活和资源释放由 veepin 负责。面板检查入站 worker
的异常退出，保留失败原因，关闭原核心实例后由既有守护任务重建；维护模式
仍禁止自动启动。换钥不重建 gVisor 栈，不改变账号路由、DNS 或出站配置。

## 发布

- 普通主分支推送：测试并构建 Linux 产物。
- 推送 `v*` 标签：发布 Linux 压缩包、`SHA256SUMS` 和 GHCR 镜像。
- 已有 Release 时上传并替换同名附件，支持失败重试；不得移动已有发布标签。
- Docker 手动工作流默认只验证构建与启动，明确选择发布后才推送镜像。
