# 开发约定

仅发布 Linux amd64/arm64，所有发行构建设置 CGO_ENABLED=0。
前端在 frontend/ 直接维护，不使用子模块；保留原始 LICENSE。
构建命令和运行条件见 README.md，协议标签统一在 build-tags.sh。

L2TP 只通过固定版本 veepin 库处理 IKE、IPsec、L2TP、PPP，不复制协议实现。
每个认证会话独立网络栈；进入 sing-box 必须携带账号身份。
不得添加 host NAT、默认路由、系统 DNS 回退或未经认证的数据转发。
新增功能必须覆盖注册、配置、前端、数据库、生命周期及失败清理。
账号变更目前采用入站重建，不能仅修改认证表而保留旧连接。

测试包括前端 lint/单测/构建、纯 Go SQLite 迁移/备份、Linux 静态检查与测试。
L2TP 网络验收在隔离命名空间中运行，禁止修改开发主机默认路由。
公网独立客户端、NAT 重绑定、长期重连与 rekey 需单独记录验收结果。
