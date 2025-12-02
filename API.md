# API 文档

基础 URL: `http://localhost:8080/api/v1`

## 认证

大多数端点需要 JWT 认证。请在 Authorization 头中包含令牌：

```
Authorization: Bearer <your_jwt_token>
```

## 端点

### 健康检查

#### GET /health

检查 API 是否正在运行。

**响应:**
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

---

## 认证

### POST /auth/register

注册新用户账户。

**请求体:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password"
}
```

**响应:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-12-02T00:00:00Z",
    "updated_at": "2024-12-02T00:00:00Z"
  }
}
```

**状态码:**
- 201: 创建成功
- 400: 请求数据无效
- 409: 用户名或邮箱已存在

---

### POST /auth/login

使用现有凭据登录。

**请求体:**
```json
{
  "username": "john_doe",
  "password": "secure_password"
}
```

**响应:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "created_at": "2024-12-02T00:00:00Z",
    "updated_at": "2024-12-02T00:00:00Z"
  }
}
```

**状态码:**
- 200: 登录成功
- 400: 请求数据无效
- 401: 凭据无效

---

### GET /auth/profile

获取当前用户个人资料。需要认证。

**响应:**
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**状态码:**
- 200: 成功
- 401: 未授权

---

## 知识库管理

### POST /knowledge/upload

上传文档到知识库。需要认证。

**请求:**
- Content-Type: `multipart/form-data`
- 字段名: `file`
- 支持的格式: PDF, DOCX, TXT, MD

**响应:**
```json
{
  "id": 1,
  "user_id": 1,
  "filename": "document.pdf",
  "file_type": ".pdf",
  "file_size": 1024000,
  "file_path": "/uploads/1_1234567890_document.pdf",
  "status": "pending",
  "chunk_count": 0,
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**状态码:**
- 201: 上传成功，开始处理
- 400: 未提供文件或文件类型无效
- 401: 未授权
- 500: 上传失败

**注意:** 文档处理是异步进行的。使用 GET /knowledge/:id 轮询文档状态

---

### GET /knowledge/list

获取当前用户的所有文档。需要认证。

**响应:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "filename": "document.pdf",
    "file_type": ".pdf",
    "file_size": 1024000,
    "file_path": "/uploads/1_1234567890_document.pdf",
    "status": "completed",
    "chunk_count": 15,
    "created_at": "2024-12-02T00:00:00Z",
    "updated_at": "2024-12-02T00:00:00Z"
  }
]
```

**状态码:**
- 200: 成功
- 401: 未授权

**文档状态值:**
- `pending`: 上传完成，等待处理
- `processing`: 正在解析和索引
- `completed`: 准备就绪
- `failed`: 处理失败

---

### GET /knowledge/:id

获取特定文档的详细信息。需要认证。

**响应:**
```json
{
  "id": 1,
  "user_id": 1,
  "filename": "document.pdf",
  "file_type": ".pdf",
  "file_size": 1024000,
  "file_path": "/uploads/1_1234567890_document.pdf",
  "status": "completed",
  "chunk_count": 15,
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**状态码:**
- 200: 成功
- 401: 未授权
- 404: 文档未找到

---

### DELETE /knowledge/:id

从知识库中删除文档。需要认证。

**响应:**
```json
{
  "message": "Document deleted successfully"
}
```

**状态码:**
- 200: 删除成功
- 401: 未授权
- 404: 文档未找到
- 500: 删除失败

---

### POST /knowledge/search

在知识库中搜索相似内容。需要认证。

**请求体:**
```json
{
  "query": "machine learning algorithms",
  "top_k": 5
}
```

**响应:**
```json
[
  {
    "ChunkID": 123,
    "DocumentID": 1,
    "Content": "Machine learning algorithms are...",
    "Score": 0.95
  }
]
```

**状态码:**
- 200: 成功
- 400: 请求无效
- 401: 未授权
- 500: 搜索失败

---

## PPT 生成

### POST /ppt/generate

生成 LaTeX PPT。需要认证。

**请求体:**
```json
{
  "title": "Introduction to AI",
  "prompt": "Create a presentation about artificial intelligence, covering history, applications, and future trends. Include 5-7 slides.",
  "template": "default",
  "document_ids": [1, 2],
  "use_openai": true
}
```

**字段:**
- `title` (必填): PPT 标题
- `prompt` (必填): 详细要求
- `template` (可选): 模板名称 (default, madrid, modern)。默认: "default"
- `document_ids` (可选): 知识库中要使用的文档 ID 数组
- `use_openai` (可选): 使用 OpenAI (true) 或 Claude (false)。默认: true

**响应:**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Introduction to AI",
  "prompt": "Create a presentation...",
  "latex_content": "\\documentclass[aspectratio=169,11pt]{beamer}...",
  "pdf_path": "/outputs/ppt_1_1234567890.pdf",
  "template": "default",
  "status": "completed",
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**状态码:**
- 200: 生成成功
- 400: 请求无效
- 401: 未授权
- 500: 生成或编译失败

**PPT 状态值:**
- `pending`: 已收到请求
- `generating`: AI 正在生成 LaTeX 代码
- `completed`: LaTeX 已生成并编译为 PDF
- `failed`: 生成或编译失败

---

### GET /ppt/templates

获取可用 LaTeX Beamer 模板列表。

**响应:**
```json
{
  "templates": ["default", "madrid", "modern"]
}
```

**状态码:**
- 200: 成功

---

### POST /ppt/compile

将 LaTeX 代码编译为 PDF。需要认证。

**请求体:**
```json
{
  "latex_content": "\\documentclass[aspectratio=169,11pt]{beamer}..."
}
```

**响应:**
```json
{
  "id": 2,
  "user_id": 1,
  "title": "Manual Compile",
  "prompt": "Manual compilation",
  "latex_content": "\\documentclass...",
  "pdf_path": "/outputs/ppt_2_1234567890.pdf",
  "template": "default",
  "status": "completed",
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**状态码:**
- 200: 编译成功
- 400: LaTeX 内容无效
- 401: 未授权
- 500: 编译失败

---

### GET /ppt/history

获取当前用户的 PPT 生成历史。需要认证。

**响应:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "title": "Introduction to AI",
    "prompt": "Create a presentation...",
    "latex_content": "\\documentclass...",
    "pdf_path": "/outputs/ppt_1_1234567890.pdf",
    "template": "default",
    "status": "completed",
    "created_at": "2024-12-02T00:00:00Z",
    "updated_at": "2024-12-02T00:00:00Z"
  }
]
```

**状态码:**
- 200: 成功
- 401: 未授权

---

### GET /ppt/:id

获取特定 PPT 记录的详细信息。需要认证。

**响应:**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Introduction to AI",
  "prompt": "Create a presentation...",
  "latex_content": "\\documentclass...",
  "pdf_path": "/outputs/ppt_1_1234567890.pdf",
  "template": "default",
  "status": "completed",
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**状态码:**
- 200: 成功
- 401: 未授权
- 404: PPT 未找到

---

### GET /ppt/:id/download

下载生成的 PDF。需要认证。

**响应:** 二进制 PDF 文件

**状态码:**
- 200: 成功 (返回 PDF 文件)
- 401: 未授权
- 404: PPT 或 PDF 未找到

---

### DELETE /ppt/:id

删除 PPT 记录。需要认证。

**响应:**
```json
{
  "message": "PPT deleted successfully"
}
```

**状态码:**
- 200: 删除成功
- 401: 未授权
- 404: PPT 未找到
- 500: 删除失败

---

## 错误响应格式

所有错误响应遵循此格式：

```json
{
  "error": "Error message description"
}
```

## 速率限制

目前未实施速率限制。对于生产使用，请考虑在以下方面实施速率限制：
- 认证端点 (防止暴力破解)
- PPT 生成 (控制 API 成本)
- 文件上传 (防止滥用)

## curl 使用示例

### 注册用户
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### 登录
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### 上传文档
```bash
curl -X POST http://localhost:8080/api/v1/knowledge/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/document.pdf"
```

### 生成 PPT
```bash
curl -X POST http://localhost:8080/api/v1/ppt/generate \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My Presentation",
    "prompt": "Create a 5-slide presentation about climate change",
    "template": "default",
    "use_openai": true
  }'
```

### 下载 PPT
```bash
curl -X GET http://localhost:8080/api/v1/ppt/1/download \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o presentation.pdf
```

## 最佳实践

1. **令牌管理**: 在客户端安全地存储 JWT 令牌
2. **错误处理**: 始终检查状态码并适当处理错误
3. **大文件**: 上传大文件时使用适当的超时值
4. **轮询**: 上传文档时，轮询状态端点而不是等待
5. **重试逻辑**: 为瞬态故障 (5xx 错误) 实现重试逻辑
