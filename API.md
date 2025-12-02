# API Documentation

Base URL: `http://localhost:8080/api/v1`

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Health Check

#### GET /health

Check if the API is running.

**Response:**
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

---

## Authentication

### POST /auth/register

Register a new user account.

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password"
}
```

**Response:**
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

**Status Codes:**
- 201: Created successfully
- 400: Invalid request data
- 409: Username or email already exists

---

### POST /auth/login

Login with existing credentials.

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "secure_password"
}
```

**Response:**
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

**Status Codes:**
- 200: Login successful
- 400: Invalid request data
- 401: Invalid credentials

---

### GET /auth/profile

Get current user profile. Requires authentication.

**Response:**
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "created_at": "2024-12-02T00:00:00Z",
  "updated_at": "2024-12-02T00:00:00Z"
}
```

**Status Codes:**
- 200: Success
- 401: Unauthorized

---

## Knowledge Base Management

### POST /knowledge/upload

Upload a document to the knowledge base. Requires authentication.

**Request:**
- Content-Type: `multipart/form-data`
- Field name: `file`
- Supported formats: PDF, DOCX, TXT, MD

**Response:**
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

**Status Codes:**
- 201: Upload successful, processing started
- 400: No file provided or invalid file type
- 401: Unauthorized
- 500: Upload failed

**Note:** Document processing happens asynchronously. Poll the document status using GET /knowledge/:id

---

### GET /knowledge/list

Get all documents for the current user. Requires authentication.

**Response:**
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

**Status Codes:**
- 200: Success
- 401: Unauthorized

**Document Status Values:**
- `pending`: Upload complete, waiting to be processed
- `processing`: Currently being parsed and indexed
- `completed`: Ready for use
- `failed`: Processing failed

---

### GET /knowledge/:id

Get details of a specific document. Requires authentication.

**Response:**
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

**Status Codes:**
- 200: Success
- 401: Unauthorized
- 404: Document not found

---

### DELETE /knowledge/:id

Delete a document from the knowledge base. Requires authentication.

**Response:**
```json
{
  "message": "Document deleted successfully"
}
```

**Status Codes:**
- 200: Deleted successfully
- 401: Unauthorized
- 404: Document not found
- 500: Deletion failed

---

### POST /knowledge/search

Search for similar content in the knowledge base. Requires authentication.

**Request Body:**
```json
{
  "query": "machine learning algorithms",
  "top_k": 5
}
```

**Response:**
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

**Status Codes:**
- 200: Success
- 400: Invalid request
- 401: Unauthorized
- 500: Search failed

---

## PPT Generation

### POST /ppt/generate

Generate a LaTeX PPT. Requires authentication.

**Request Body:**
```json
{
  "title": "Introduction to AI",
  "prompt": "Create a presentation about artificial intelligence, covering history, applications, and future trends. Include 5-7 slides.",
  "template": "default",
  "document_ids": [1, 2],
  "use_openai": true
}
```

**Fields:**
- `title` (required): PPT title
- `prompt` (required): Detailed requirements
- `template` (optional): Template name (default, madrid, modern). Default: "default"
- `document_ids` (optional): Array of document IDs to use from knowledge base
- `use_openai` (optional): Use OpenAI (true) or Claude (false). Default: true

**Response:**
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

**Status Codes:**
- 200: Generation successful
- 400: Invalid request
- 401: Unauthorized
- 500: Generation or compilation failed

**PPT Status Values:**
- `pending`: Request received
- `generating`: AI is generating LaTeX code
- `completed`: LaTeX generated and compiled to PDF
- `failed`: Generation or compilation failed

---

### GET /ppt/templates

Get list of available LaTeX Beamer templates.

**Response:**
```json
{
  "templates": ["default", "madrid", "modern"]
}
```

**Status Codes:**
- 200: Success

---

### POST /ppt/compile

Compile LaTeX code to PDF. Requires authentication.

**Request Body:**
```json
{
  "latex_content": "\\documentclass[aspectratio=169,11pt]{beamer}..."
}
```

**Response:**
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

**Status Codes:**
- 200: Compilation successful
- 400: Invalid LaTeX content
- 401: Unauthorized
- 500: Compilation failed

---

### GET /ppt/history

Get PPT generation history for the current user. Requires authentication.

**Response:**
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

**Status Codes:**
- 200: Success
- 401: Unauthorized

---

### GET /ppt/:id

Get details of a specific PPT record. Requires authentication.

**Response:**
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

**Status Codes:**
- 200: Success
- 401: Unauthorized
- 404: PPT not found

---

### GET /ppt/:id/download

Download the generated PDF. Requires authentication.

**Response:** Binary PDF file

**Status Codes:**
- 200: Success (returns PDF file)
- 401: Unauthorized
- 404: PPT or PDF not found

---

### DELETE /ppt/:id

Delete a PPT record. Requires authentication.

**Response:**
```json
{
  "message": "PPT deleted successfully"
}
```

**Status Codes:**
- 200: Deleted successfully
- 401: Unauthorized
- 404: PPT not found
- 500: Deletion failed

---

## Error Response Format

All error responses follow this format:

```json
{
  "error": "Error message description"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. For production use, consider implementing rate limiting on:
- Authentication endpoints (to prevent brute force)
- PPT generation (to control API costs)
- File uploads (to prevent abuse)

## Example Usage with curl

### Register a User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Upload Document
```bash
curl -X POST http://localhost:8080/api/v1/knowledge/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@/path/to/document.pdf"
```

### Generate PPT
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

### Download PPT
```bash
curl -X GET http://localhost:8080/api/v1/ppt/1/download \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -o presentation.pdf
```

## Best Practices

1. **Token Management**: Store JWT tokens securely on the client side
2. **Error Handling**: Always check status codes and handle errors appropriately
3. **Large Files**: Use appropriate timeout values when uploading large documents
4. **Polling**: When uploading documents, poll the status endpoint rather than waiting
5. **Retry Logic**: Implement retry logic for transient failures (5xx errors)
