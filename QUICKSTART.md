# Quick Start Guide

Get up and running with LaTeX PPT Generator in under 5 minutes!

## Prerequisites

- Docker & Docker Compose installed
- OpenAI or Claude API key
- 4GB+ available RAM

## 🚀 Start in 3 Steps

### 1. Configure
```bash
cp .env.example .env
```

Edit `.env` and add your API key:
```env
OPENAI_API_KEY=sk-your-key-here
# OR
CLAUDE_API_KEY=your-claude-key-here
```

### 2. Launch
```bash
docker-compose up -d
```

Wait 1-2 minutes for all services to start.

### 3. Access
Open http://localhost:3000 in your browser

## 📋 Quick Commands

```bash
# Start all services
make up

# Stop all services
make down

# View logs
make logs

# Check status
make status

# Restart services
make restart

# Clean everything
make clean
```

## 🎯 First Use Workflow

1. **Register** - Create your account at http://localhost:3000/login
2. **Upload** (optional) - Add documents to knowledge base
3. **Generate** - Create your first PPT:
   - Title: "My First PPT"
   - Prompt: "Create a 5-slide presentation about AI"
   - Click "Generate PPT"
4. **Download** - Save your PDF

## 🔧 Troubleshooting

### Services won't start?
```bash
# Check if ports are available
docker-compose ps
docker-compose logs
```

### Backend can't connect to database?
```bash
# Wait for PostgreSQL to be ready
docker-compose logs postgres
# Should see "database system is ready to accept connections"
```

### Milvus connection issues?
```bash
# Milvus takes longer to start (1-2 minutes)
docker-compose logs milvus
# Wait for health check to pass
```

## 📱 Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost:3000 | Web interface |
| Backend API | http://localhost:8080 | REST API |
| Health Check | http://localhost:8080/api/v1/health | API status |
| PostgreSQL | localhost:5432 | Database |
| Milvus | localhost:19530 | Vector DB |
| MinIO Console | http://localhost:9001 | Storage |

## 🔑 Default Credentials

### MinIO (Object Storage)
- Username: `minioadmin`
- Password: `minioadmin`
- Access at: http://localhost:9001

### PostgreSQL
- User: `postgres`
- Password: `postgres`
- Database: `latex_ppt`

## 📖 Next Steps

- Read [SETUP.md](SETUP.md) for detailed setup
- Check [API.md](API.md) for API documentation
- See [CONTRIBUTING.md](CONTRIBUTING.md) to contribute

## ⚡ Performance Tips

1. **First startup is slow** - Docker needs to download images
2. **Milvus needs time** - Wait 1-2 minutes after `docker-compose up`
3. **Increase memory** - Allocate more RAM to Docker if needed
4. **Use SSD** - Better performance with SSD storage

## 🐛 Common Issues

### "Connection refused" error
- Services haven't finished starting yet
- Check with: `docker-compose ps`
- All services should show "Up" and "healthy"

### "Out of memory" error
- Increase Docker memory limit
- Close other applications
- Check: `docker stats`

### LaTeX compilation fails
- Check backend logs: `docker-compose logs backend`
- Ensure XeLaTeX is installed in container

### Frontend shows "Network Error"
- Backend might not be running
- Check: `curl http://localhost:8080/api/v1/health`

## 💡 Tips

- **Use knowledge base**: Upload relevant documents for better PPT content
- **Choose right model**: OpenAI for speed, Claude for quality
- **Edit LaTeX**: Generated code is editable before compilation
- **Save history**: All generated PPTs are saved automatically

## 🎓 Learn More

| Resource | Purpose |
|----------|---------|
| [README.md](README.md) | Project overview and features |
| [SETUP.md](SETUP.md) | Detailed installation guide |
| [API.md](API.md) | Complete API reference |
| [latex.md](latex.md) | LaTeX Beamer guidelines |

## 💬 Get Help

1. Check the logs: `make logs`
2. Review documentation in this repo
3. Search existing [GitHub Issues](https://github.com/qingbingwei/latex_ppt_by_claude/issues)
4. Create a new issue if needed

---

**Ready to create amazing presentations? Start now! 🚀**
