# AGENTS.md

## 项目概述

桌面日历（desktopcal）—— 基于 Wails v2 的桌面日历应用，支持农历、节气和事件管理。

## 技术栈

- **后端**: Go 1.25.14, Wails v3 (`github.com/wailsapp/wails/v3`), SQLite (go-sqlite3)
- **前端**: Vue 3 (Composition API), TypeScript, Pinia, Vite
- **构建**: Makefile (后端), npm + Vite (前端)

## 项目结构

```
.
├── main.go              # Wails 入口，初始化窗口、托盘和数据库
├── app.go               # App 结构体，日历/事件等绑定到前端的方法
├── app_data.go          # 设置读写、数据导入导出（绑定方法）
├── app_window.go        # 窗口控制（隐藏/显示/置顶/退出，绑定方法）
├── internal/
│   ├── db/
│   │   ├── db.go            # 连接与建表迁移
│   │   ├── db_event.go      # 事件 CRUD
│   │   ├── db_background.go # 日期背景色
│   │   ├── db_settings.go   # 设置键值读写
│   │   └── db_import.go     # 导入覆盖写入
│   ├── lunar/
│   │   ├── lunar.go         # 公历转农历、节气查询
│   │   ├── lunar_calc.go    # 历法计算
│   │   └── lunar_data.go    # 1900–2100 农历数据表
│   └── model/model.go   # 数据模型定义
├── frontend/
│   ├── src/
│   │   ├── main.ts      # Vue 入口
│   │   ├── App.vue      # 根组件
│   │   ├── api.ts       # Wails 后端调用封装
│   │   ├── style.css    # 全局样式与主题 CSS 变量
│   │   ├── types.ts     # TypeScript 类型定义
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
```

## 编码规范

### Go 后端

- 遵循标准 Go 格式（gofmt）
- 错误信息使用中文（面向用户的校验提示）
- 内部包放在 `internal/` 下
- 数据模型集中在 `internal/model/model.go`
- 数据库操作集中在 `internal/db/db.go`

### 前端

- 使用 Vue 3 Composition API + `<script setup>` 语法
- TypeScript strict 模式，不使用 any
- 路径别名：`@/` 映射到 `frontend/src/`
- 组件使用 PascalCase 命名
- Wails 绑定代码（`frontend/bindings/`）不要手动修改，由 `wails` 生成

### 通用

- 缩进：Go 使用 Tab，前端使用 2 空格
- 编码：UTF-8
- 换行符：LF
- 提交信息使用中文
- 读写操作仅限于该工程目录： `$(HOME)/workspaces/desktopcal`
- 如需要其它权限和操作的话，一定要审核

## 注意事项

- `frontend/bindings/` 目录由 Wails 自动生成，不要手动编辑
- 当前使用 Wails **v3**（`wails.json` 的 `$schema` 仍指向 v2 配置，属历史遗留）
- `main.go` 中系统托盘目前被注释掉（`TEMP: tray disabled for isolation`），恢复时取消相关注释即可
- SQLite 使用 CGO，构建时需要 `CGO_ENABLED=1`
- macOS 构建的 `CGO_CFLAGS` 与 `CGO_LDFLAGS` 必须带相同的最低版本
  （`-mmacosx-version-min=12.0`，见 Makefile 的 `MACOS_MIN`），否则会刷屏
  `built for newer macOS version (12.0) than being linked (10.13)` 链接警告
- Go 使用自定义 GOROOT：`$(HOME)/.gvm/gos/go1.25.14`
