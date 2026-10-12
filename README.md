# 桌面日历 (desktopcal)

将月历直接嵌入桌面的跨平台日历应用，支持农历、节气、事件记录、数据导入导出。

- 后端：Go 1.25 + Wails v3 + SQLite (纯 Go 驱动 `modernc.org/sqlite`)
- 前端：Vue 3 (Composition API) + TypeScript + Vite
- 平台：macOS / Windows（帧无边框窗口 + 系统托盘 + 全局快捷键）

## 功能概览

- 月视图嵌入桌面，宽松/紧凑两种布局
- 农历、节气显示（1900–2100 数据）
- 双击日期直接记录任务，支持完成态、排序、背景色
- 系统托盘 + 全局快捷键 (`Cmd/Ctrl+Shift+C`) 显隐窗口
- JSON 备份 与 桌面日历官方文本格式 导入导出
- 多主题、透明度、字号可调，开机自启动、关闭到托盘

完整功能清单见 [docs/product-features.md](docs/product-features.md)。

## 快速开始

### 环境准备

- Go ≥ 1.25（本工程使用 GVM 的 `$(HOME)/.gvm/gos/go1.25.14`）
- Node.js ≥ 20，npm
- Wails CLI v3：`wails3`（本工程**不使用** v2）
- macOS 构建需 Xcode 命令行工具（cgo + Objective-C）

### 开发模式

```bash
make run          # 构建前端 + 编译后端 + 启动应用
```

### 生产构建

```bash
make build        # 输出到 build/bin/desktopcal.app
make frontend     # 仅构建前端
make clean        # 清理构建产物
```

### 交叉编译到 Windows

SQLite 驱动是纯 Go，可直接关闭 cgo：

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  nice -n 15 env GOMAXPROCS=2 \
  go build -p 2 -o desktopcal.exe .
```

### 重新生成前端绑定

改动 `internal/model` 或新增/删除 App 绑定方法后**必须执行**：

```bash
make bindings
# 等价于：wails3 generate bindings -ts ./...
```

> `frontend/bindings/` 目录由工具生成，不要手动修改。

## 项目结构

```
.
├── main.go                    # Wails 入口：窗口、托盘、数据库初始化
├── app.go                     # App 结构体与依赖注入
├── app_calendar.go            # 月视图数据聚合（农历/节气/事件/背景色）
├── app_event.go               # 事件 CRUD、背景色
├── app_settings.go            # 设置项声明式编解码
├── app_data.go                # 数据导入导出
├── app_hotkey.go              # 全局快捷键注册
├── app_import_text.go         # 官方文本格式解析
├── app_window.go              # 窗口控制
├── internal/
│   ├── db/                    # SQLite 访问层（按 db_*.go 分文件）
│   ├── lunar/                 # 农历/节气计算 + 数据表
│   └── model/                 # 共享数据模型
├── frontend/
│   ├── src/                   # Vue 源码（App.vue / api.ts / theme.ts / components）
│   └── bindings/              # Wails 自动生成，勿手动编辑
├── build/darwin/              # macOS Info.plist、图标、托盘图
├── docs/product-features.md   # 详细功能规格
├── Makefile
└── wails.json
```

## 关键约定

- **单文件不超过 150 行**（Go），超出按职责拆分
- 编译/构建命令必须限流：`GOMAXPROCS=2 -p 2 nice -n 15`，且**串行**执行
- 提交信息使用中文
- 用户可见错误提示使用中文
- 设置项集中在 `app_settings.go` 的 `settingsSpec` 声明，新增字段需同步更新 `model.DefaultSettings()`、`make bindings`、`App.vue` 初始值
- 主题 CSS 变量在 `frontend/src/theme.ts` 单点维护

更多约定见 [AGENTS.md](AGENTS.md)。

## 作者

yueyt
