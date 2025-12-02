# 贡献指南

感谢您有兴趣为本项目做出贡献！本文档提供了参与项目的指南和说明。

## 快速开始

1. Fork 本仓库
2. 克隆您的 Fork：`git clone https://github.com/YOUR_USERNAME/latex_ppt_by_claude.git`
3. 创建新分支：`git checkout -b feature/your-feature-name`
4. 进行更改
5. 彻底测试您的更改
6. 提交并推送到您的 Fork
7. 创建 Pull Request

## 开发环境设置

详细的设置说明请参阅 [SETUP.md](SETUP.md)。

开发快速启动：
```bash
# 启动基础设施服务
make dev

# 在单独的终端中运行：
make backend-dev
make frontend-dev
```

## 项目结构

```
├── backend/           # Go 后端
│   ├── cmd/          # 应用程序入口点
│   ├── internal/     # 私有应用程序代码
│   │   ├── api/     # HTTP 处理程序和路由
│   │   ├── service/ # 业务逻辑
│   │   ├── repository/ # 数据访问
│   │   ├── model/   # 数据模型
│   │   └── config/  # 配置
│   └── pkg/         # 公共库
│       ├── ai/      # AI 客户端实现
│       ├── embedding/ # Embedding 生成
│       ├── vectordb/ # 向量数据库客户端
│       ├── parser/  # 文档解析器
│       └── latex/   # LaTeX 编译
├── frontend/        # Vue3 前端
│   └── src/
│       ├── api/     # API 客户端
│       ├── components/ # Vue 组件
│       ├── views/   # 页面组件
│       ├── store/   # Pinia 存储
│       ├── router/  # Vue Router 配置
│       └── utils/   # 工具函数
└── docker/          # Docker 配置
```

## 代码规范

### 后端 (Go)

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 指南
- 使用 `gofmt` 格式化代码
- 提交前运行 `go vet`
- 为导出的函数和类型添加注释
- 保持函数专注且简短
- 使用有意义的变量名

示例：
```go
// GeneratePPT creates a new PPT based on user prompt and optional knowledge base
// GeneratePPT 根据用户提示和可选知识库创建新的 PPT
func (s *PPTService) GeneratePPT(ctx context.Context, userID uint, prompt string) (*model.PPTRecord, error) {
    // Implementation
}
```

### 前端 (Vue3/TypeScript)

- 遵循 Vue 3 Composition API 最佳实践
- 使用 TypeScript 进行类型安全检查
- 使用 `<script setup>` 的 Composition API
- 保持组件小巧且专注
- 使用 Pinia 进行状态管理
- 遵循 Element Plus 组件模式

示例：
```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { User } from '@/types'

const user = ref<User | null>(null)

onMounted(async () => {
  // Implementation
})
</script>
```

### 通用规范

- 编写清晰的提交信息
- 为复杂逻辑添加注释
- 更改功能时更新文档
- 为新功能编写测试（如适用）

## 添加新功能

### 后端功能

1. **添加模型** (如果需要)
   - 在 `backend/internal/model/` 中创建
   - 添加 GORM 标签用于数据库映射

2. **添加 Repository 方法**
   - 在 `backend/internal/repository/` 中创建/更新
   - 添加 CRUD 操作

3. **添加 Service 逻辑**
   - 在 `backend/internal/service/` 中创建/更新
   - 实现业务逻辑

4. **添加 API Handler**
   - 在 `backend/internal/api/handler/` 中创建处理程序
   - 处理 HTTP 请求/响应

5. **更新 Router**
   - 在 `backend/internal/api/router.go` 中添加路由

6. **更新文档**
   - 将 API 端点添加到 `API.md`

### 前端功能

1. **添加类型**
   - 在 `frontend/src/types/` 中定义 TypeScript 类型

2. **添加 API 客户端**
   - 在 `frontend/src/api/` 中创建 API 函数

3. **添加 Store** (如果需要)
   - 在 `frontend/src/store/` 中创建 Pinia store

4. **添加组件/视图**
   - 在适当的目录中创建 Vue 组件
   - 使用带 TypeScript 的 Composition API

5. **更新 Router** (如果是新页面)
   - 在 `frontend/src/router/index.ts` 中添加路由

6. **更新文档**
   - 在 `README.md` 中更新功能描述

## 测试

### 后端测试

```bash
cd backend
go test ./...
```

### 前端测试

目前未设置自动化测试。需要手动测试：
1. 测试所有 UI 交互
2. 测试 API 集成
3. 测试错误处理
4. 在不同屏幕尺寸上测试

### 集成测试

1. 启动所有服务：`make up`
2. 测试完整工作流程：
   - 用户注册和登录
   - 文档上传和处理
   - PPT 生成
   - PDF 编译和下载

## 常见开发任务

### 添加新的 AI 提供商

1. 在 `backend/pkg/ai/` 中创建客户端
2. 实现 `GenerateLaTeX` 方法
3. 更新 `AIService` 以使用新提供商
4. 在 `.env.example` 中添加配置
5. 更新文档

### 添加新的文档解析器

1. 在 `backend/pkg/parser/` 中创建解析器
2. 实现 `Parse` 方法
3. 更新 `GetParser` 函数
4. 将文件类型添加到上传验证
5. 更新文档

### 添加新的 LaTeX 模板

1. 在 `backend/pkg/latex/templates.go` 中添加模板字符串
2. 将模板名称添加到 `ListTemplates()`
3. 更新前端模板选择器
4. 记录模板功能

## Pull Request 指南

### 提交前

- [ ] 代码遵循项目风格指南
- [ ] 代码编译无错误
- [ ] 完成手动测试
- [ ] 文档已更新
- [ ] 提交信息清晰

### PR 描述

包括：
1. **What**: 更改描述
2. **Why**: 更改原因
3. **How**: 实现方法
4. **Testing**: 如何测试更改
5. **Screenshots**: UI 更改的截图

示例：
```markdown
## Add support for Markdown tables in PPT generation

### What
Added support for parsing and rendering Markdown tables in generated LaTeX presentations.

### Why
Users requested ability to include tabular data in presentations.

### How
- Extended Markdown parser to detect tables
- Added LaTeX `tabular` environment generation
- Updated AI prompt to include table formatting instructions

### Testing
- Tested with various table sizes
- Verified LaTeX compilation
- Checked PDF output quality

### Screenshots
[Attach screenshots of generated tables]
```

## 代码审查流程

1. 维护者审查 PR
2. 请求更改或批准
3. 批准后，合并 PR
4. 删除分支

## 报告 Bug

使用 GitHub Issues 并提供以下信息：

```markdown
**Description**
Clear description of the bug

**Steps to Reproduce**
1. Go to '...'
2. Click on '...'
3. See error

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: [e.g., Ubuntu 22.04]
- Docker version: [e.g., 24.0.0]
- Browser: [e.g., Chrome 120]

**Logs**
Relevant error logs or screenshots
```

## 请求功能

使用 GitHub Issues 并提供：

```markdown
**Feature Description**
Clear description of the proposed feature

**Use Case**
Why this feature would be useful

**Proposed Implementation**
(Optional) Ideas on how to implement

**Alternatives Considered**
Other solutions you've thought about
```

## 社区

- 保持尊重和建设性
- 尽力帮助他人
- 分享知识和经验
- 遵循 [行为准则](CODE_OF_CONDUCT.md) (如果存在)

## 许可证

通过贡献，您同意您的贡献将根据与项目相同的许可证（MIT 许可证）进行许可。

## 有问题？

- 检查现有问题和文档
- 在 GitHub Discussions 中提问（如果已启用）
- 创建带有 "question" 标签的新 issue

## 谢谢！

您的贡献有助于让这个项目对每个人都变得更好！
