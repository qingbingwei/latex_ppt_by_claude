# 安装指南

本指南将帮助您设置和运行 LaTeX PPT 生成器应用程序。

## 前置要求

### Docker 部署（推荐）
- Docker (版本 20.10 或更高)
- Docker Compose (版本 2.0 或更高)
- 至少 4GB 可用内存
- OpenAI API Key 或 Claude API Key

### 本地开发
- Go 1.21 或更高
- Node.js 18 或更高
- PostgreSQL 15
- Milvus 2.3
- XeLaTeX (用于 PDF 编译)

## Docker 快速开始

### 1. 克隆仓库
```bash
git clone https://github.com/qingbingwei/latex_ppt_by_claude.git
cd latex_ppt_by_claude
```

### 2. 配置环境变量
```bash
cp .env.example .env
```

编辑 `.env` 并添加您的 API 密钥：
```env
OPENAI_API_KEY=your-openai-api-key-here
CLAUDE_API_KEY=your-claude-api-key-here
JWT_SECRET=your-secure-random-secret
```

### 3. 启动应用程序
```bash
make up
# 或者
docker-compose up -d
```

这将启动所有服务：
- PostgreSQL (端口 5432)
- Milvus (端口 19530)
- 后端 API (端口 8080)
- 前端 (端口 3000)

### 4. 访问应用程序
- 打开浏览器并访问：http://localhost:3000
- 注册新账号或登录
- 开始生成 PPT！

### 5. 检查服务状态
```bash
make status
# 或者
docker-compose ps
```

### 6. 查看日志
```bash
make logs
# 或者
docker-compose logs -f
```

## 本地开发设置

### 后端开发

1. **安装 Go 依赖**
```bash
cd backend
go mod download
```

2. **启动所需服务**
```bash
# 在项目根目录
make dev
```

3. **配置环境**
在 backend 目录中创建 `.env` 文件，内容与根目录 `.env` 相同。

4. **运行后端**
```bash
cd backend
go run cmd/server/main.go
```

后端 API 将在 http://localhost:8080 上可用

### 前端开发

1. **安装依赖**
```bash
cd frontend
npm install
```

2. **运行开发服务器**
```bash
npm run dev
```

前端将在 http://localhost:3000 上可用

### 构建后端

```bash
cd backend
go build -o server ./cmd/server
./server
```

### 构建前端

```bash
cd frontend
npm run build
```

构建的文件将在 `frontend/dist/` 中

## 验证安装

### 检查后端健康状态
```bash
curl http://localhost:8080/api/v1/health
```

预期响应：
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

### 检查前端
在浏览器中打开 http://localhost:3000。您应该能看到主页。

## 常见问题与解决方案

### 问题：后端无法连接到 PostgreSQL
**解决方案**：确保 PostgreSQL 容器正在运行且健康
```bash
docker-compose ps postgres
docker-compose logs postgres
```

### 问题：Milvus 连接超时
**解决方案**：Milvus 启动需要时间。等待几分钟并检查日志：
```bash
docker-compose logs milvus
```

### 问题：LaTeX 编译失败
**解决方案**：确保后端容器中安装了 XeLaTeX。Dockerfile 默认包含它。

### 问题：前端无法访问后端 API
**解决方案**：检查后端是否在端口 8080 上运行，并且已启用 CORS。

### 问题：构建时出现 "No space left on device"
**解决方案**：清理 Docker 资源：
```bash
make clean
docker system prune -a
```

## 数据库管理

### 初始化/重置数据库
```bash
make init-db
```

### 访问 PostgreSQL
```bash
docker-compose exec postgres psql -U postgres -d latex_ppt
```

### 备份数据库
```bash
docker-compose exec postgres pg_dump -U postgres latex_ppt > backup.sql
```

### 恢复数据库
```bash
cat backup.sql | docker-compose exec -T postgres psql -U postgres -d latex_ppt
```

## 更新应用程序

### 拉取最新更改
```bash
git pull origin main
```

### 重建并重启
```bash
make down
make build
make up
```

## 性能调优

### 生产部署

1. **增加 Milvus 内存**
编辑 `docker-compose.yml`：
```yaml
milvus:
  deploy:
    resources:
      limits:
        memory: 8G
```

2. **为生产环境配置 PostgreSQL**
添加到 docker-compose.yml postgres 服务：
```yaml
command: postgres -c shared_buffers=256MB -c max_connections=200
```

3. **启用 Redis 缓存（可选）**
添加 Redis 服务并配置后端使用它。

## 安全注意事项

1. **更改默认 JWT 密钥**
   - 永远不要在生产环境中使用默认的 JWT_SECRET
   - 生成强随机密钥：`openssl rand -base64 32`

2. **使用特定环境的配置**
   - 为开发和生产使用不同的 `.env` 文件
   - 永远不要将 `.env` 文件提交到版本控制

3. **启用 HTTPS**
   - 使用反向代理 (nginx/traefik) 和 SSL 证书
   - 配置 Let's Encrypt 自动续期证书

4. **保护 API 密钥**
   - 将 API 密钥存储在环境变量或密钥管理系统中
   - 定期轮换 API 密钥

## 监控

### 查看实时日志
```bash
# 所有服务
docker-compose logs -f

# 特定服务
docker-compose logs -f backend
docker-compose logs -f frontend
```

### 检查资源使用情况
```bash
docker stats
```

## 故障排除命令

```bash
# 重启所有服务
make restart

# 停止所有服务
make down

# 删除所有数据并重新开始
make clean
make up

# 查看后端日志
make backend-logs

# 查看前端日志
make frontend-logs

# 检查服务健康状态
docker-compose ps
```

## 支持

如有问题和疑问：
- 检查 [GitHub Issues](https://github.com/qingbingwei/latex_ppt_by_claude/issues)
- 查看 [README.md](README.md) 获取一般信息
- 检查 Docker 和服务日志中的错误消息

## 下一步

设置完成后：
1. 注册新账号
2. (可选) 上传文档到知识库
3. 导航到生成页面
4. 创建您的第一个 AI 驱动的 LaTeX PPT！

有关功能和用法的更多详细信息，请参阅 [README.md](README.md)。
