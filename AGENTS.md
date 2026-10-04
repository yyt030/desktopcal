# AGENTS.md

## 项目概述

桌面日历（desktopcal）—— 基于 Wails v2 的桌面日历应用，支持农历、节气和事件管理。

## 技术栈

- **后端**: Go 1.25.14, Wails v3 (`github.com/wailsapp/wails/v3`), SQLite (modernc.org/sqlite，纯 Go)
- **前端**: Vue 3 (Composition API), TypeScript, Pinia, Vite
- **构建**: Makefile (后端), npm + Vite (前端)

## 项目结构

```
.
├── main.go              # Wails 入口，初始化窗口、托盘和数据库
├── app.go               # App 结构体与运行时依赖注入（绑定服务根）
├── app_calendar.go      # 月视图数据：日期格子、农历/节气、事件与背景色聚合
├── app_event.go         # 事件 CRUD、日期背景色、今日农历/节气（绑定方法）
├── app_settings.go      # 设置项的声明式编解码与批量读写（绑定方法）
├── app_data.go          # 数据导入导出（绑定方法）
├── app_window.go        # 窗口控制（隐藏/显示/置顶/退出，绑定方法）
├── internal/
│   ├── db/
│   │   ├── db.go            # 连接（WAL/busy_timeout）、连接池与建表迁移
│   │   ├── db_event.go      # 事件 CRUD
│   │   ├── db_background.go # 日期背景色
│   │   ├── db_settings.go   # 设置键值读写
│   │   └── db_import.go     # 导入覆盖写入
│   ├── lunar/
│   │   ├── lunar.go         # 公历转农历、节气查询
│   │   ├── lunar_calc.go    # 历法计算
│   │   └── lunar_data.go    # 1900–2100 农历数据表
│   └── model/model.go   # 数据模型定义、日期格式/布局常量与默认设置
├── frontend/
│   ├── src/
│   │   ├── main.ts      # Vue 入口
│   │   ├── App.vue      # 根组件
│   │   ├── api.ts       # Wails 后端调用封装
│   │   ├── style.css    # 全局样式与主题 CSS 变量
│   │   ├── theme.ts     # 主题定义（CSS 变量 + 应用主题），设置页共用
│   │   ├── types.ts     # 前端 DTO（派生自 Wails 绑定模型）
│   │   ├── utils/       # 通用工具（日期等）
│   │   └── components/  # Vue 组件
│   └── bindings/        # Wails 自动生成的绑定代码（勿手动修改）
├── Makefile             # 构建命令
└── wails.json           # Wails 配置
```

## 常用命令

```bash
# 开发模式（构建前端 + 编译后端 + 运行）
make run

# 生产构建
make build

# 仅构建前端
make frontend

# 清理构建产物
make clean

# 前端开发服务器
cd frontend && npm run dev

# 重新生成前端绑定（改动 internal/model 或绑定方法后必须执行）
make bindings
```

> 构建会触发 Go 全量编译。**务必按「通用」规范限制资源**（`GOMAXPROCS=2 -p 2`
> + `nice -n 15`，且一次只跑一个构建）；交叉编译到 Windows 尤其耗 CPU。

## 编码规范

### Go 后端

- 遵循标准 Go 格式（gofmt）
- 错误信息使用中文（面向用户的校验提示）
- 内部包放在 `internal/` 下
- 数据模型集中在 `internal/model/model.go`
- 数据库操作集中在 `internal/db/` 下，按 `db_*.go` 分文件
- 绑定服务方法按职责分文件（`app_*.go`）
- **单个 Go 文件不超过 150 行**，超出按职责拆分（见 `docs/product-features.md` §8.1）

### 前端

- 使用 Vue 3 Composition API + `<script setup>` 语法
- TypeScript strict 模式，不使用 any
- 路径别名：`@/` 映射到 `frontend/src/`
- 组件使用 PascalCase 命名
- Wails 绑定代码（`frontend/bindings/`）不要手动修改，由  
  `wails3 generate bindings -ts ./...` 生成

### 通用

- 缩进：Go 使用 Tab，前端使用 2 空格
- 编码：UTF-8
- 换行符：LF
- 提交信息使用中文
- 读写操作仅限于该工程目录： `$(HOME)/workspaces/desktopcal`
- 如需要其它权限和操作的话，一定要审核
- **编译/构建必须限制本机资源占用**：`go build` / `go test` / `go vet` 等命令
  一律带 `GOMAXPROCS=2` 与 `-p 2`，并以 `nice -n 15` 降低优先级，例如
  `nice -n 15 env GOMAXPROCS=2 go build -p 2 -o out .`；同时禁止并发执行多个
  构建命令（一次只跑一个）。交叉编译（`GOOS=windows`）会从零编译整套 Wails + cgo，
  CPU 占用极高，尤其必须限流，否则会导致本机风扇狂转

## 注意事项

- `frontend/bindings/` 目录由 Wails 自动生成，不要手动编辑
- 当前使用 Wails **v3**（`wails.json` 的 `$schema` 仍指向 v2 配置，属历史遗留）
- 系统托盘在 `main.go` 的 `setupSystemTray` 中实现：左键切换主窗口显隐  
  （可见且聚焦才隐藏，被遮挡时唤到前台），右键弹出菜单（显示/隐藏/退出）。  
  菜单栏图标是 `build/darwin/trayicon.png`（纯黑 + alpha 的模板图，系统按  
  菜单栏明暗自动反色），改动后运行 `build/gen_icon.sh` 重新生成。  
  退出必须先置 `quitting` 标志再调 `app.Quit()`，否则会被窗口关闭 hook 拦截
- 数据库驱动是**纯 Go** 的 `modernc.org/sqlite`（已弃用需要 CGO 的 `mattn/go-sqlite3`），
  因此 SQLite 部分可以 `CGO_ENABLED=0` 直接交叉编译到 Windows，无需 mingw 工具链：
  `GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o desktopcal.exe .`
- DSN 用 `_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)` 传 PRAGMA，且**文件名不带
  `file:` 前缀**：该驱动仅在带 `file:` 时才把整串交给 SQLite 按 URI 解析，否则会剥掉参数
  按普通路径打开 —— 这样 Windows 的盘符路径（`C:\...`）才不会被 URI 规则误判
- macOS 构建仍需 `CGO_ENABLED=1`（Wails 的 darwin 实现是 Objective-C，走 cgo），
  `CGO_CFLAGS` 与 `CGO_LDFLAGS` 必须带相同的最低版本  
  （`-mmacosx-version-min=12.0`，见 Makefile 的 `MACOS_MIN`），否则会刷屏  
  `built for newer macOS version (12.0) than being linked (10.13)` 链接警告
- Go 使用自定义 GOROOT：`$(HOME)/.gvm/gos/go1.25.14`
- 设置项集中在 `app_settings.go` 的 `settingsSpec` 声明（键名 + 编解码）：  
  新增设置项时同步更新 `internal/model/model.go` 的 `DefaultSettings()`、执行  
  `make bindings` 重新生成绑定，以及 `App.vue` 的初始值
- 主题的 CSS 变量统一在 `frontend/src/theme.ts` 维护，设置页与根组件共用同一份定义
- `frontend/src/types.ts` 的前端 DTO 派生自生成模型，不再手写字段：`Event` /  
  `LunarDay` / `Settings` 是生成模型的直接别名，`DayInfo` 仅把 `events` 从  
  `(Event | null)[]` 收窄为 `Event[]`。因此后端加字段后不用改这个文件，只需  
  重新生成绑定
- `frontend/bindings/` 由 `wails3 generate bindings -ts ./...` 生成。**必须带 `-ts`**，  
  否则默认输出带 JSDoc 的 `.js`（`@ts-check`）而非 `.ts`
- 本机 Wails CLI 在 `$(HOME)/.gvm/pkgsets/go1.25.14/global/bin/`：`wails3` 是 v3  
  （v3.0.0-beta.27），`wails` 是 v2（v2.16.0），**本工程只能用 `wails3`**
- 改动 `internal/model` 或新增/删除绑定方法后，必须重新生成绑定，否则  
  `frontend/bindings/` 的类型会落后于后端（历史上就漏过 `Settings.layoutMode`  
  与 `DayInfo.isEmpty`）
