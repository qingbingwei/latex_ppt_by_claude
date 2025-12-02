# 快速开始指南

在 5 分钟内启动并运行 LaTeX PPT 生成器！

## 前置要求

- 已安装 Docker 和 Docker Compose
- 已安装 Node.js (用于 copilot-api)
- GitHub Copilot 订阅（推荐）或 OpenAI/Claude API 密钥
- 4GB+ 可用内存

## 🚀 快速启动 (使用 GitHub Copilot)

### 1. 安装并配置 copilot-api

```bash
# 安装 copilot-api
npm install -g copilot-api

# 进行 GitHub 认证 (会打开浏览器)
copilot-api auth

# 启动代理服务 (保持运行)
copilot-api start --port 4141
```

### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env` 文件：
```env
# 使用 copilot-api 代理 (推荐)
OPENAI_API_KEY=dummy-key
OPENAI_BASE_URL=http://host.docker.internal:4141/v1

# 注意: 使用 host.docker.internal 而不是 localhost
# 因为 Docker 容器内的 localhost 指向容器自己，不是主机
# host.docker.internal 是 Docker 提供的特殊域名，指向宿主机
```

### 3. 启动服务

```bash
docker-compose up -d
```

等待 1-2 分钟让所有服务启动。

### 4. 访问

在浏览器中打开 http://localhost:3000

## 🔧 其他 AI 配置方式

### 使用 OpenAI 官方 API
```env
OPENAI_API_KEY=sk-your-key-here
OPENAI_BASE_URL=https://api.openai.com/v1
```

### 使用 Claude API
```env
CLAUDE_API_KEY=your-claude-key-here
```

### 使用硅基流动 (SiliconFlow) 免费 API
```env
OPENAI_API_KEY=your-siliconflow-key
OPENAI_BASE_URL=https://api.siliconflow.cn/v1
```

## ⚠️ 关于 host.docker.internal

```
┌─────────────────────────────────────────────────────────┐
│                    你的 Mac/PC 主机                      │
│                                                         │
│   ┌─────────────────┐                                   │
│   │  copilot-api    │  ← 运行在主机上 (localhost:4141)  │
│   └────────▲────────┘                                   │
│            │                                            │
│            │ host.docker.internal:4141                  │
│   ┌────────┴────────────────────────────────────┐       │
│   │              Docker 容器网络                  │       │
│   │  ┌──────────┐  ┌──────────┐                 │       │
│   │  │ backend  │  │ frontend │                 │       │
│   │  │          │  │          │                 │       │
│   │  └──────────┘  └──────────┘                 │       │
│   │                                              │       │
│   │  容器内 localhost ≠ 主机 localhost           │       │
│   └──────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────┘
```

- Docker 容器有独立的网络命名空间
- 容器内的 `localhost` 指向容器自己，**不是主机**
- `host.docker.internal` 是 Docker 提供的特殊 DNS，用于访问主机服务

## 📋 快速命令

```bash
# 启动所有服务
make up

# 停止所有服务
make down

# 查看日志
make logs

# 检查状态
make status

# 重启服务
make restart

# 清理所有内容
make clean
```

## 🎯 首次使用流程

1. **注册** - 在 http://localhost:3000/login 创建您的账户
2. **上传** (可选) - 将文档添加到知识库
3. **生成** - 创建您的第一个 PPT：
   - 标题："我的第一个 PPT"
   - 提示词："创建一个关于 AI 的 5 页演示文稿"
   - 点击 "Generate PPT"
4. **下载** - 保存您的 PDF

## 🔧 故障排除

### 服务无法启动？
```bash
# 检查端口是否可用
docker-compose ps
docker-compose logs
```

### 后端无法连接到数据库？
```bash
# 等待 PostgreSQL 准备就绪
docker-compose logs postgres
# 应该看到 "database system is ready to accept connections"
```

### Milvus 连接问题？
```bash
# Milvus 启动需要较长时间 (1-2 分钟)
docker-compose logs milvus
# 等待健康检查通过
```

## 📱 服务 URL

| 服务 | URL | 描述 |
|---------|-----|-------------|
| 前端 | http://localhost:3000 | Web 界面 |
| 后端 API | http://localhost:8080 | REST API |
| 健康检查 | http://localhost:8080/api/v1/health | API 状态 |
| PostgreSQL | localhost:5432 | 数据库 |
| Milvus | localhost:19530 | 向量数据库 |
| MinIO 控制台 | http://localhost:9001 | 存储 |

## 🔑 默认凭据

### MinIO (对象存储)
- 用户名: `minioadmin`
- 密码: `minioadmin`
- 访问地址: http://localhost:9001

### PostgreSQL
- 用户: `postgres`
- 密码: `postgres`
- 数据库: `latex_ppt`

## 📖 下一步

- 阅读 [SETUP.md](SETUP.md) 了解详细设置
- 查看 [API.md](API.md) 获取 API 文档
- 参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 进行贡献

## ⚡ 性能提示

1. **首次启动较慢** - Docker 需要下载镜像
2. **Milvus 需要时间** - `docker-compose up` 后等待 1-2 分钟
3. **增加内存** - 如果需要，为 Docker 分配更多内存
4. **使用 SSD** - SSD 存储性能更好

## 🐛 常见问题

### "Connection refused" 错误
- 服务尚未完成启动
- 检查：`docker-compose ps`
- 所有服务应显示 "Up" 和 "healthy"

### "Out of memory" 错误
- 增加 Docker 内存限制
- 关闭其他应用程序
- 检查：`docker stats`

### LaTeX 编译失败
- 检查后端日志：`docker-compose logs backend`
- 确保容器中安装了 XeLaTeX

### 前端显示 "Network Error"
- 后端可能未运行
- 检查：`curl http://localhost:8080/api/v1/health`

## 💡 提示

- **使用知识库**：上传相关文档以获得更好的 PPT 内容
- **选择合适的模型**：OpenAI 速度快，Claude 质量高
- **编辑 LaTeX**：生成的代码在编译前可编辑
- **保存历史**：所有生成的 PPT 都会自动保存

## 🎓 了解更多

| 资源 | 用途 |
|----------|---------|
| [README.md](README.md) | 项目概述和功能 |
| [SETUP.md](SETUP.md) | 详细安装指南 |
| [API.md](API.md) | 完整 API 参考 |
| [latex.md](latex.md) | LaTeX Beamer 指南 |

## 💬 获取帮助

1. 检查日志：`make logs`
2. 查看本仓库中的文档
3. 搜索现有的 [GitHub Issues](https://github.com/qingbingwei/latex_ppt_by_claude/issues)
4. 如果需要，创建一个新 issue

---

**准备好创建精彩的演示文稿了吗？现在开始！🚀**
