# LaTeX PPT Generator by AI

基于AI生成LaTeX制作PPT的完整解决方案，支持RAG（检索增强生成）、多种AI模型集成，以及完整的前后端实现。

## 项目简介

这是一个基于AI的LaTeX Beamer演示文稿生成系统，提供以下核心功能：

- 🤖 **AI驱动**: 支持OpenAI GPT-4和Claude API生成高质量LaTeX代码
- 📚 **RAG知识库**: 上传文档构建知识库，生成PPT时自动检索相关内容
- 🎨 **多种模板**: 提供多种Beamer主题模板，支持中文
- 📄 **自动编译**: 自动将LaTeX编译为PDF，支持预览和下载
- 💻 **现代化界面**: 基于Vue3 + Element Plus的响应式Web界面

## 技术栈

### 后端
- **框架**: Gin (Go Web Framework)
- **数据库**: PostgreSQL (关系数据库)
- **向量数据库**: Milvus (用于RAG)
- **AI集成**: OpenAI API / Claude API
- **ORM**: GORM
- **文档解析**: 支持PDF、DOCX、TXT、Markdown

### 前端
- **框架**: Vue 3 + Composition API + TypeScript
- **UI库**: Element Plus
- **状态管理**: Pinia
- **构建工具**: Vite
- **路由**: Vue Router
- **HTTP客户端**: Axios

### 基础设施
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **LaTeX编译**: XeLaTeX (支持中文)

## 快速开始

### 前置要求

- Docker 和 Docker Compose
- 至少 4GB 可用内存
- OpenAI API Key 或 GitHub Copilot 订阅

### 一键部署

```bash
# 克隆项目
git clone https://github.com/qingbingwei/latex_ppt_by_claude.git
cd latex_ppt_by_claude

# 运行一键部署脚本
./scripts/deploy.sh
```

脚本会自动：
- ✅ 检查系统依赖 (Docker, Docker Compose)
- ✅ 配置 API Token (支持 GitHub Copilot / OpenAI)
- ✅ 生成安全的 JWT 密钥
- ✅ 构建 Docker 镜像
- ✅ 启动所有服务
- ✅ 等待服务就绪

### 部署脚本参数

```bash
./scripts/deploy.sh              # 完整部署
./scripts/deploy.sh --dev        # 开发模式 (只启动基础服务)
./scripts/deploy.sh --clean      # 清理后重新部署
./scripts/deploy.sh --skip-token # 跳过 Token 配置
```

### 手动部署

1. 克隆项目:
```bash
git clone https://github.com/qingbingwei/latex_ppt_by_claude.git
cd latex_ppt_by_claude
```

2. 配置环境变量:
```bash
cp .env.example .env
# 编辑 .env 文件，填入你的 API Keys
```

3. 启动所有服务:
```bash
make up
# 或者
docker-compose up -d
```

4. 访问应用:
- 前端: http://localhost:3000
- 后端API: http://localhost:8080

### 开发模式

启动基础服务（数据库、Milvus）:
```bash
make dev
```

在单独的终端运行后端:
```bash
make backend-dev
```

在另一个终端运行前端:
```bash
make frontend-dev
```

## 使用指南

### 1. 注册/登录

首次使用需要注册账号，之后可以使用用户名和密码登录。

### 2. 上传知识库文档（可选）

在"Knowledge Base"页面上传相关文档（PDF、DOCX、TXT、Markdown），系统会自动：
- 解析文档内容
- 分割成chunks
- 生成向量embeddings
- 存储到Milvus向量数据库

### 3. 生成PPT

在"Generate"页面：
1. 输入PPT标题和详细要求
2. 选择Beamer模板（default、madrid、modern）
3. 选择AI模型（OpenAI或Claude）
4. 如果开启"Use Knowledge Base"，系统会从知识库检索相关内容
5. 点击"Generate PPT"

生成的LaTeX代码会自动编译为PDF，可以在线预览和下载。

### 4. 查看历史

在"History"页面可以查看所有生成的PPT记录，支持查看详情、重新下载和删除。

## API文档

### 认证相关

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/auth/profile` - 获取用户信息

### 知识库管理

- `POST /api/v1/knowledge/upload` - 上传文档
- `GET /api/v1/knowledge/list` - 文档列表
- `GET /api/v1/knowledge/:id` - 文档详情
- `DELETE /api/v1/knowledge/:id` - 删除文档
- `POST /api/v1/knowledge/search` - 向量检索

### PPT生成

- `POST /api/v1/ppt/generate` - 生成PPT
- `GET /api/v1/ppt/templates` - 获取模板列表
- `POST /api/v1/ppt/compile` - 编译LaTeX
- `GET /api/v1/ppt/history` - 生成历史
- `GET /api/v1/ppt/:id` - PPT详情
- `GET /api/v1/ppt/:id/download` - 下载PPT
- `DELETE /api/v1/ppt/:id` - 删除PPT记录

## 项目结构

```
latex_ppt_by_claude/
├── backend/                    # Go后端
│   ├── cmd/server/            # 主程序入口
│   ├── internal/              # 内部包
│   │   ├── api/              # API层
│   │   ├── service/          # 业务逻辑层
│   │   ├── repository/       # 数据访问层
│   │   ├── model/            # 数据模型
│   │   └── config/           # 配置管理
│   └── pkg/                  # 公共包
│       ├── ai/               # AI客户端
│       ├── embedding/        # Embedding生成
│       ├── vectordb/         # 向量数据库
│       ├── parser/           # 文档解析
│       └── latex/            # LaTeX编译
├── frontend/                  # Vue3前端
│   ├── src/
│   │   ├── api/              # API调用
│   │   ├── components/       # Vue组件
│   │   ├── views/            # 页面视图
│   │   ├── store/            # Pinia状态管理
│   │   ├── router/           # 路由配置
│   │   └── utils/            # 工具函数
│   └── public/               # 静态资源
├── docker/                    # Docker配置
├── scripts/                   # 脚本文件
└── docker-compose.yml        # Docker Compose配置
```

## 常用命令

```bash
# 查看所有可用命令
make help

# 构建镜像
make build

# 启动服务
make up

# 停止服务
make down

# 查看日志
make logs

# 查看服务状态
make status

# 重启服务
make restart

# 清理所有数据
make clean
```

## 环境变量说明

在 `.env` 文件中配置以下环境变量：

```bash
# AI API配置
# 支持 OpenAI 兼容的 API，包括：
# - GitHub Copilot: https://api.githubcopilot.com (推荐)
# - OpenAI 官方: https://api.openai.com/v1
OPENAI_API_KEY=your-github-copilot-token-or-openai-api-key
OPENAI_BASE_URL=https://api.githubcopilot.com
CLAUDE_API_KEY=your-claude-api-key

# JWT配置
JWT_SECRET=your-jwt-secret-key-change-this-in-production
```

### 获取 GitHub Copilot Token

如果使用 GitHub Copilot API，可以通过 GitHub CLI 获取 Token：

```bash
# 1. 安装 GitHub CLI
brew install gh

# 2. 登录 GitHub（首次使用需要授权）
gh auth login

# 3. 获取 Token
gh auth token
```

将获取到的 Token（格式如 `gho_xxxx`）填入 `OPENAI_API_KEY` 即可。

## 故障排除

### 数据库连接失败
确保PostgreSQL容器已启动并健康：
```bash
docker-compose ps postgres
docker-compose logs postgres
```

### Milvus连接失败
Milvus需要一定时间启动，检查状态：
```bash
docker-compose ps milvus
docker-compose logs milvus
```

### LaTeX编译失败
- 确保LaTeX内容格式正确
- 查看后端日志获取详细错误信息
- 中文支持需要正确的字体配置

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 联系方式

如有问题，请提交Issue或联系维护者。
