# Setup Guide

This guide will help you set up and run the LaTeX PPT Generator application.

## Prerequisites

### For Docker Deployment (Recommended)
- Docker (version 20.10 or higher)
- Docker Compose (version 2.0 or higher)
- At least 4GB of available RAM
- OpenAI API Key or Claude API Key

### For Local Development
- Go 1.21 or higher
- Node.js 18 or higher
- PostgreSQL 15
- Milvus 2.3
- XeLaTeX (for PDF compilation)

## Quick Start with Docker

### 1. Clone the Repository
```bash
git clone https://github.com/qingbingwei/latex_ppt_by_claude.git
cd latex_ppt_by_claude
```

### 2. Configure Environment Variables
```bash
cp .env.example .env
```

Edit `.env` and add your API keys:
```env
OPENAI_API_KEY=your-openai-api-key-here
CLAUDE_API_KEY=your-claude-api-key-here
JWT_SECRET=your-secure-random-secret
```

### 3. Start the Application
```bash
make up
# or
docker-compose up -d
```

This will start all services:
- PostgreSQL (port 5432)
- Milvus (port 19530)
- Backend API (port 8080)
- Frontend (port 3000)

### 4. Access the Application
- Open your browser and navigate to: http://localhost:3000
- Register a new account or login
- Start generating PPTs!

### 5. Check Service Status
```bash
make status
# or
docker-compose ps
```

### 6. View Logs
```bash
make logs
# or
docker-compose logs -f
```

## Local Development Setup

### Backend Development

1. **Install Go Dependencies**
```bash
cd backend
go mod download
```

2. **Start Required Services**
```bash
# In project root
make dev
```

3. **Configure Environment**
Create `.env` file in backend directory with same content as root `.env`.

4. **Run Backend**
```bash
cd backend
go run cmd/server/main.go
```

The backend API will be available at http://localhost:8080

### Frontend Development

1. **Install Dependencies**
```bash
cd frontend
npm install
```

2. **Run Development Server**
```bash
npm run dev
```

The frontend will be available at http://localhost:3000

### Building the Backend

```bash
cd backend
go build -o server ./cmd/server
./server
```

### Building the Frontend

```bash
cd frontend
npm run build
```

The built files will be in `frontend/dist/`

## Verifying Installation

### Check Backend Health
```bash
curl http://localhost:8080/api/v1/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "Service is running"
}
```

### Check Frontend
Open http://localhost:3000 in your browser. You should see the home page.

## Common Issues and Solutions

### Issue: Backend fails to connect to PostgreSQL
**Solution**: Ensure PostgreSQL container is running and healthy
```bash
docker-compose ps postgres
docker-compose logs postgres
```

### Issue: Milvus connection timeout
**Solution**: Milvus takes time to start. Wait a few minutes and check logs:
```bash
docker-compose logs milvus
```

### Issue: LaTeX compilation fails
**Solution**: Ensure XeLaTeX is installed in the backend container. The Dockerfile includes it by default.

### Issue: Frontend can't reach backend API
**Solution**: Check that backend is running on port 8080 and CORS is enabled.

### Issue: "No space left on device" when building
**Solution**: Clean up Docker resources:
```bash
make clean
docker system prune -a
```

## Database Management

### Initialize/Reset Database
```bash
make init-db
```

### Access PostgreSQL
```bash
docker-compose exec postgres psql -U postgres -d latex_ppt
```

### Backup Database
```bash
docker-compose exec postgres pg_dump -U postgres latex_ppt > backup.sql
```

### Restore Database
```bash
cat backup.sql | docker-compose exec -T postgres psql -U postgres -d latex_ppt
```

## Updating the Application

### Pull Latest Changes
```bash
git pull origin main
```

### Rebuild and Restart
```bash
make down
make build
make up
```

## Performance Tuning

### For Production Deployment

1. **Increase Memory for Milvus**
Edit `docker-compose.yml`:
```yaml
milvus:
  deploy:
    resources:
      limits:
        memory: 8G
```

2. **Configure PostgreSQL for Production**
Add to docker-compose.yml postgres service:
```yaml
command: postgres -c shared_buffers=256MB -c max_connections=200
```

3. **Enable Redis for Caching (Optional)**
Add Redis service and configure backend to use it.

## Security Considerations

1. **Change Default JWT Secret**
   - Never use the default JWT_SECRET in production
   - Generate a strong random secret: `openssl rand -base64 32`

2. **Use Environment-Specific Configs**
   - Use different `.env` files for development and production
   - Never commit `.env` files to version control

3. **Enable HTTPS**
   - Use a reverse proxy (nginx/traefik) with SSL certificates
   - Configure Let's Encrypt for automatic certificate renewal

4. **Secure API Keys**
   - Store API keys in environment variables or secret management system
   - Rotate API keys regularly

## Monitoring

### View Real-Time Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
```

### Check Resource Usage
```bash
docker stats
```

## Troubleshooting Commands

```bash
# Restart all services
make restart

# Stop all services
make down

# Remove all data and start fresh
make clean
make up

# View backend logs
make backend-logs

# View frontend logs
make frontend-logs

# Check service health
docker-compose ps
```

## Support

For issues and questions:
- Check the [GitHub Issues](https://github.com/qingbingwei/latex_ppt_by_claude/issues)
- Review the [README.md](README.md) for general information
- Check Docker and service logs for error messages

## Next Steps

After setup:
1. Register a new account
2. (Optional) Upload documents to the knowledge base
3. Navigate to the Generate page
4. Create your first AI-powered LaTeX PPT!

For more detailed information about features and usage, see [README.md](README.md).
