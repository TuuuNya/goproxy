# goproxy

goproxy 是一个自动化代理抓取、验证与分发的高性能 SOCKS5 代理服务器。只需一行命令即可在本地快速启动 SOCKS5 代理，并智能聚合全球公开代理资源，自动检测、优选最优链路，极大提升匿名性与连接质量。

---

## 快速启动本地 SOCKS5 代理

**最常用用法：一行命令即可本地部署 SOCKS5 代理**

```bash
go run main.go serve
```

或编译后直接运行：

```bash
./goproxy serve
```

默认将会在本地 `7777` 端口启动 SOCKS5 代理服务，自动抓取大量全球代理，实时检测和分发。

---

### serve 命令参数说明

你可以通过如下参数自定义代理行为：

- `-p, --port int`              监听端口（默认 7777）
- `-t, --type string`           代理类型，仅支持 socks5（默认 "socks5"）
- `--max_delay duration`        允许的最大代理延迟（默认 10s，支持如 2s、500ms 这种 Go 时间格式）
- `-c, --check_pool_size int`   检测代理的并发池大小（默认 100000，建议根据机器性能调整）
- `-d, --debug`                 启用 debug 日志

#### 示例用法

- 在 1080 端口启动，限制最大延迟 2 秒，检测池 200 并发
  ```bash
  ./goproxy serve -p 1080 --max_delay 2s -c 200
  ```

- 启用 debug 日志，便于问题排查
  ```bash
  ./goproxy serve -d
  ```

- 查看所有 serve 参数
  ```bash
  ./goproxy serve -h
  ```

---

## 其他命令简介

- `goproxy find`  
  仅抓取并输出代理列表，不启动本地代理。支持 `-t` 指定类型。
  ```bash
  ./goproxy find -t socks5,https
  ```

- `goproxy completion`  
  生成 shell 自动补全脚本。

---

## 插件机制与自定义 Provider 扩展

goproxy 支持通过自定义 Provider 实现新代理源扩展。你可以轻松集成任意网站或 API 的代理抓取能力。

### 新增 Provider 步骤

1. **实现 Provider 接口**

   在 `internal/provider/` 新建文件（比如 `myprovider.go`），实现如下接口：

   ```go
   type Provider interface {
       Name() string
       SupportedTypes() []string
       FetchProxies(types []string) ([]Proxy, error)
   }
   ```

2. **实现核心方法**
   ```go
   type MyProvider struct{}

   func (m *MyProvider) Name() string { return "myprovider" }
   func (m *MyProvider) SupportedTypes() []string { return []string{"socks5"} }
   func (m *MyProvider) FetchProxies(types []string) ([]Proxy, error) {
       // 这里实现你的代理抓取逻辑，返回 Proxy 切片
   }
   ```

3. **注册新 Provider**
   在 `internal/engine/finder.go` 的 providers 列表中加入你的 Provider，例如：

   ```go
   providers := []provider.Provider{
       &provider.ProxyListOrg{},
       &provider.XseoIn{},
       // ...
       &provider.MyProvider{},  // 加入你自己的
   }
   ```

4. **编译运行即可自动生效**
   新 Provider 会被自动调用并合并到代理池，实现聚合高质量代理。

---

## 目录结构简述

- `internal/engine/`：主逻辑（代理抓取、检测、分发）
- `internal/provider/`：各类代理源插件（可自定义扩展）
- `pkg/logger/`：日志模块
- `main.go`：程序入口

---

## 声明

本项目仅供学习与研究使用，请勿用于任何非法用途。请遵循相关法律法规。

---

如需更多帮助或功能扩展建议，欢迎提 issue 或 PR！
