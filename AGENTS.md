# AGENTS.md

## 项目概述

桌面日历（desktopcal）—— 基于 Wails v2 的桌面日历应用，支持农历、节气和事件管理。

## 技术栈

- **后端**: Go 1.25.14, Wails v2, SQLite (go-sqlite3)
- **前端**: Vue 3 (Composition API), TypeScript, Pinia, Vite
- **构建**: Makefile (后端), npm + Vite (前端)

## 项目结构

```
.
├── main.go              # Wails 入口，初始化窗口和数据库
├── app.go               # App 结构体，绑定到前端的方法
├── internal/
│   ├── db/db.go         # SQLite 数据库操作
│   ├── lunar/lunar.go   # 农历/节气计算
│   └── model/model.go   # 数据模型定义
├── frontend/
│   ├── src/
│   │   ├── main.ts      # Vue 入口
│   │   ├── App.vue      # 根组件
│   │   ├── api.ts       # Wails 后端调用封装
│   │   ├── types.ts     # TypeScript 类型定义
│   │   ├── components/  # Vue 组件
│   │   └── stores/      # Pinia 状态管理
│   └── wailsjs/         # Wails 自动生成的绑定代码（勿手动修改）
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
- Wails 绑定代码（`wailsjs/`）不要手动修改

### 通用

- 缩进：Go 使用 Tab，前端使用 2 空格
- 编码：UTF-8
- 换行符：LF
- 提交信息使用中文

## 注意事项

- `frontend/wailsjs/` 目录由 Wails 自动生成，不要手动编辑
- SQLite 使用 CGO，构建时需要 `CGO_ENABLED=1`
- macOS 构建需要 `CGO_LDFLAGS="-framework UniformTypeIdentifiers"`
- Go 使用自定义 GOROOT：`$(HOME)/.gvm/gos/go1.25.14`
