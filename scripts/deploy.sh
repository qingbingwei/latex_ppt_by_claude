#!/bin/bash
#
# LaTeX PPT Generator 一键部署脚本
# 
# 使用方法:
#   chmod +x scripts/deploy.sh
#   ./scripts/deploy.sh
#
# 或者带参数运行:
#   ./scripts/deploy.sh --skip-token    # 跳过 Token 配置
#   ./scripts/deploy.sh --dev           # 开发模式
#   ./scripts/deploy.sh --clean         # 清理后重新部署
#

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║           LaTeX PPT Generator 一键部署脚本                   ║"
    echo "║                   Powered by AI                              ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# 检查依赖
check_dependencies() {
    print_info "检查系统依赖..."
    
    local missing_deps=()
    
    # 检查 Docker
    if ! command -v docker &> /dev/null; then
        missing_deps+=("docker")
    fi
    
    # 检查 Docker Compose
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        missing_deps+=("docker-compose")
    fi
    
    # 检查 GitHub CLI (可选，用于获取 token)
    if ! command -v gh &> /dev/null; then
        print_warning "GitHub CLI (gh) 未安装，如需使用 GitHub Copilot API，请先安装: brew install gh"
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "缺少以下依赖: ${missing_deps[*]}"
        echo ""
        echo "请先安装依赖:"
        echo "  macOS:   brew install docker docker-compose"
        echo "  Ubuntu:  sudo apt install docker.io docker-compose"
        exit 1
    fi
    
    # 检查 Docker 是否运行
    if ! docker info &> /dev/null; then
        print_error "Docker 未运行，请先启动 Docker Desktop"
        exit 1
    fi
    
    print_success "依赖检查通过"
}

# 配置环境变量
setup_env() {
    print_info "配置环境变量..."
    
    if [ -f ".env" ]; then
        print_info "检测到已有 .env 文件"
        read -p "是否覆盖现有配置? (y/N): " overwrite
        if [[ ! "$overwrite" =~ ^[Yy]$ ]]; then
            print_info "保留现有配置"
            return
        fi
    fi
    
    # 复制示例配置
    # cp .env.example .env
    
    # 配置 API Token
    if [[ "$SKIP_TOKEN" != "true" ]]; then
        configure_api_token
    fi
    
    # 生成 JWT Secret
    JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)
    sed -i.bak "s/your-jwt-secret-key-change-this-in-production/$JWT_SECRET/" .env
    rm -f .env.bak
    
    print_success "环境变量配置完成"
}

# 配置 API Token
configure_api_token() {
    echo ""
    echo "请选择 AI API 类型:"
    echo "  1) GitHub Copilot API (推荐，需要 Copilot 订阅)"
    echo "  2) OpenAI 官方 API"
    echo "  3) 其他 OpenAI 兼容 API"
    echo "  4) 跳过配置 (稍后手动配置)"
    echo ""
    read -p "请输入选项 [1-4]: " api_choice
    
    case $api_choice in
        1)
            configure_github_copilot
            ;;
        2)
            configure_openai
            ;;
        3)
            configure_custom_api
            ;;
        4)
            print_warning "跳过 API 配置，请稍后手动编辑 .env 文件"
            ;;
        *)
            print_warning "无效选项，跳过 API 配置"
            ;;
    esac
}

# 配置 GitHub Copilot
configure_github_copilot() {
    print_info "配置 GitHub Copilot API..."
    
    # 检查是否安装了 copilot-api
    if ! command -v copilot-api &> /dev/null; then
        print_warning "copilot-api 未安装"
        echo "安装命令: npm install -g copilot-api"
        read -p "是否现在安装? (Y/n): " install_choice
        if [[ ! "$install_choice" =~ ^[Nn]$ ]]; then
            sudo npm install -g copilot-api
        else
            print_warning "跳过安装，请稍后手动安装并配置"
            return
        fi
    fi
    
    # 检查是否已认证
    if [ ! -f "$HOME/.local/share/copilot-api/github_token" ]; then
        print_info "需要进行 GitHub 认证..."
        copilot-api auth
    fi
    
    # 配置 .env 使用 copilot-api
    sed -i.bak "s|OPENAI_API_KEY=.*|OPENAI_API_KEY=dummy-key|" .env
    sed -i.bak "s|OPENAI_BASE_URL=.*|OPENAI_BASE_URL=http://host.docker.internal:4141/v1|" .env
    rm -f .env.bak
    
    # 启动 copilot-api
    print_info "启动 copilot-api 代理服务..."
    # 先杀掉可能存在的旧进程
    pkill -f "copilot-api" 2>/dev/null || true
    sleep 1
    
    # 后台启动
    nohup copilot-api start --port 4141 > /tmp/copilot-api.log 2>&1 &
    sleep 3
    
    # 验证服务是否启动
    if curl -s http://localhost:4141/v1/models > /dev/null 2>&1; then
        print_success "copilot-api 代理服务启动成功"
        print_info "注意: Docker 容器通过 host.docker.internal:4141 访问此服务"
    else
        print_error "copilot-api 启动失败，请检查日志: cat /tmp/copilot-api.log"
    fi
}

# 配置 OpenAI
configure_openai() {
    print_info "配置 OpenAI API..."
    read -p "请输入 OpenAI API Key (sk-xxx): " api_key
    
    if [[ -n "$api_key" ]]; then
        sed -i.bak "s|OPENAI_API_KEY=.*|OPENAI_API_KEY=$api_key|" .env
        sed -i.bak "s|OPENAI_BASE_URL=.*|OPENAI_BASE_URL=https://api.openai.com/v1|" .env
        rm -f .env.bak
        print_success "OpenAI API 配置成功"
    else
        print_warning "未输入 API Key，跳过配置"
    fi
}

# 配置自定义 API
configure_custom_api() {
    print_info "配置自定义 OpenAI 兼容 API..."
    read -p "请输入 API Base URL: " base_url
    read -p "请输入 API Key: " api_key
    
    if [[ -n "$base_url" && -n "$api_key" ]]; then
        sed -i.bak "s|OPENAI_API_KEY=.*|OPENAI_API_KEY=$api_key|" .env
        sed -i.bak "s|OPENAI_BASE_URL=.*|OPENAI_BASE_URL=$base_url|" .env
        rm -f .env.bak
        print_success "自定义 API 配置成功"
    else
        print_warning "配置不完整，跳过"
    fi
}

# 构建镜像
build_images() {
    print_info "构建 Docker 镜像..."
    docker-compose build --parallel
    print_success "镜像构建完成"
}

# 启动服务
start_services() {
    print_info "启动服务..."
    docker-compose up -d
    print_success "服务启动中..."
}

# 等待服务就绪
wait_for_services() {
    print_info "等待服务就绪..."
    
    local max_attempts=60
    local attempt=0
    
    # 等待后端服务
    echo -n "等待后端服务"
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
            echo ""
            print_success "后端服务就绪"
            break
        fi
        echo -n "."
        sleep 2
        ((attempt++))
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo ""
        print_warning "后端服务启动超时，请检查日志: docker-compose logs backend"
    fi
    
    # 等待前端服务
    attempt=0
    echo -n "等待前端服务"
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:3000 > /dev/null 2>&1; then
            echo ""
            print_success "前端服务就绪"
            break
        fi
        echo -n "."
        sleep 2
        ((attempt++))
    done
    
    if [ $attempt -eq $max_attempts ]; then
        echo ""
        print_warning "前端服务启动超时，请检查日志: docker-compose logs frontend"
    fi
}

# 显示服务状态
show_status() {
    echo ""
    print_info "服务状态:"
    docker-compose ps
    echo ""
}

# 显示访问信息
show_access_info() {
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║                    🎉 部署完成！                              ║${NC}"
    echo -e "${GREEN}╠══════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${GREEN}║${NC}                                                              ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}  前端地址:     ${BLUE}http://localhost:3000${NC}                        ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}  后端 API:     ${BLUE}http://localhost:8080${NC}                        ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}  健康检查:     ${BLUE}http://localhost:8080/api/v1/health${NC}          ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}  MinIO 控制台: ${BLUE}http://localhost:9001${NC}                        ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}                                                              ${GREEN}║${NC}"
    echo -e "${GREEN}╠══════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${GREEN}║${NC}  常用命令:                                                    ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}    查看日志:   ${YELLOW}make logs${NC}                                    ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}    停止服务:   ${YELLOW}make down${NC}                                    ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}    重启服务:   ${YELLOW}make restart${NC}                                 ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}    查看状态:   ${YELLOW}make status${NC}                                  ${GREEN}║${NC}"
    echo -e "${GREEN}║${NC}                                                              ${GREEN}║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# 清理部署
clean_deploy() {
    print_warning "清理现有部署..."
    docker-compose down -v
    print_success "清理完成"
}

# 开发模式
dev_mode() {
    print_info "启动开发模式..."
    docker-compose up postgres milvus etcd minio -d
    
    echo ""
    print_info "基础服务已启动，请在新终端中运行:"
    echo "  后端: make backend-dev"
    echo "  前端: make frontend-dev"
    echo ""
}

# 主函数
main() {
    # 解析参数
    SKIP_TOKEN="false"
    DEV_MODE="false"
    CLEAN_MODE="false"
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-token)
                SKIP_TOKEN="true"
                shift
                ;;
            --dev)
                DEV_MODE="true"
                shift
                ;;
            --clean)
                CLEAN_MODE="true"
                shift
                ;;
            --help|-h)
                echo "使用方法: $0 [选项]"
                echo ""
                echo "选项:"
                echo "  --skip-token    跳过 API Token 配置"
                echo "  --dev           开发模式 (只启动基础服务)"
                echo "  --clean         清理后重新部署"
                echo "  --help, -h      显示帮助信息"
                exit 0
                ;;
            *)
                print_error "未知参数: $1"
                exit 1
                ;;
        esac
    done
    
    print_banner
    
    # 检查依赖
    check_dependencies
    
    # 清理模式
    if [[ "$CLEAN_MODE" == "true" ]]; then
        clean_deploy
    fi
    
    # 配置环境
    setup_env
    
    # 开发模式
    if [[ "$DEV_MODE" == "true" ]]; then
        dev_mode
        exit 0
    fi
    
    # 构建镜像
    build_images
    
    # 启动服务
    start_services
    
    # 等待服务就绪
    wait_for_services
    
    # 显示状态
    show_status
    
    # 显示访问信息
    show_access_info
}

# 运行主函数
main "$@"
