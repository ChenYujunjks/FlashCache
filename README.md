
# FlashCache

FlashCache 是一个基于 Go 实现的高并发内存缓存服务，采用面向生产环境的设计思路，集成了并发控制、过期淘汰、限流保护、监控指标以及容器化部署能力，用于模拟真实后端基础组件的设计与实现。

---

## 🚀 项目特点

- 高性能内存 Key-Value 缓存
- 支持 TTL 过期与后台清理机制
- 基于分片锁（Sharded Lock）实现高并发访问
- 提供 HTTP 接口供外部服务调用
- 内置令牌桶限流（Token Bucket）
- 支持 Prometheus 指标采集
- 结构化日志输出
- 支持 Docker / docker-compose 部署
- 基础 CI/CD（GitHub Actions）

---

## 🧱 项目结构

``` id="structure"
flashcache/
├── cmd/server/           # 程序入口
├── internal/cache/       # 缓存核心逻辑
├── internal/api/         # HTTP / TCP 接口层
├── internal/limiter/     # 限流模块
├── internal/metrics/     # 监控指标
├── internal/logger/      # 日志系统
├── internal/config/      # 配置管理
├── configs/              # 配置文件
├── deployments/          # Docker / 部署相关
├── scripts/              # 压测与辅助脚本
├── test/                 # 测试代码
└── README.md
````

---

## ⚙️ 快速开始

### 本地运行

```bash
go run cmd/server/main.go
```

测试接口：

```bash
curl localhost:8080/ping
```

---

### Docker 运行（后续支持）

```bash
docker-compose up --build
```

---

## 📡 核心功能（规划中）

* [x] 基础 HTTP 服务
* [ ] SET / GET / DELETE
* [ ] TTL / EXPIRE
* [ ] 分片锁并发模型
* [ ] 限流机制
* [ ] Prometheus 监控
* [ ] Grafana 可视化
* [ ] TCP 协议支持
* [ ] LRU / LFU 淘汰策略
* [ ] 持久化（Snapshot / AOF）

---

## 📊 监控指标（规划）

* QPS（请求吞吐）
* 请求成功率
* 平均响应时间
* P95 / P99 延迟
* 当前 Key 数量
* 过期 Key 清理次数
* 限流触发次数
* Goroutine 数量
* 内存使用情况

---

## 🧪 测试与压测（规划）

```bash
go test ./...
```

压测脚本：

```bash
bash scripts/load_test.sh
```

---

## 🎯 项目目标

本项目的目标不是实现一个简单的 KV 存储，而是模拟真实后端服务在高并发场景下的设计与工程实践，包括：

* 并发控制与性能优化
* 服务稳定性与限流保护
* 可观测性（日志 + 监控）
* 工程化能力（CI/CD、容器化）
* 系统设计与模块解耦

---

## 🛠 技术栈

* Go
* Prometheus（规划）
* Grafana（规划）
* Docker
* GitHub Actions

---

## 📌 说明

该项目主要用于学习与实践：

* Go 并发模型（goroutine / mutex / channel）
* 高并发系统设计
* 缓存系统核心机制
* 后端服务工程化能力




