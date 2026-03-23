#!/bin/bash
# ML_Elec Deployment Script
# Simplifies Docker deployment with a single command

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."

    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    # Check Docker Compose
    if ! command -v docker compose &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose v2+."
        exit 1
    fi

    # Check Docker is running
    if ! docker info &> /dev/null; then
        log_error "Docker daemon is not running. Please start Docker."
        exit 1
    fi

    log_success "Prerequisites check passed"
}

# Setup environment file
setup_env() {
    if [ ! -f .env ]; then
        log_info "Creating .env file from template..."
        cp .env.example .env
        
        # Generate JWT secret if on Linux/macOS
        if command -v openssl &> /dev/null; then
            JWT_SECRET=$(openssl rand -hex 32)
            sed -i "s/JWT_SECRET=change-me-in-production/JWT_SECRET=$JWT_SECRET/" .env
            log_success "JWT secret generated"
        else
            log_warning "OpenSSL not found. Please manually set JWT_SECRET in .env file"
        fi
    else
        log_info ".env file already exists"
    fi
}

# Create required directories
setup_directories() {
    log_info "Creating required directories..."
    mkdir -p mqtt/config
    mkdir -p packages/edge-agent/pki
}

# Main deployment
deploy() {
    log_info "Starting ML_Elec deployment..."
    
    # Pull and build images
    log_info "Building Docker images (this may take a few minutes)..."
    docker compose build
    
    # Start services
    log_info "Starting services..."
    docker compose up -d
    
    # Wait for services to be healthy
    log_info "Waiting for services to be healthy..."
    sleep 10
    
    # Check health
    if docker compose ps | grep -q "healthy"; then
        log_success "ML_Elec deployed successfully!"
        echo ""
        echo "=========================================="
        echo "  🚀 ML_Elec is running!"
        echo "=========================================="
        echo ""
        echo "  📊 Dashboard: http://localhost:3000"
        echo "  🔧 API:       http://localhost:3001"
        echo "  📨 MQTT:      localhost:1883"
        echo ""
        echo "  Useful commands:"
        echo "    - View logs:      docker compose logs -f"
        echo "    - Stop:           docker compose down"
        echo "    - Restart:        docker compose restart"
        echo "    - Check status:   docker compose ps"
        echo ""
        echo "=========================================="
    else
        log_warning "Some services may still be starting..."
        log_info "Check status with: docker compose ps"
    fi
}

# Show help
show_help() {
    echo "ML_Elec Deployment Script"
    echo ""
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  deploy    Deploy ML_Elec (default)"
    echo "  start     Start existing deployment"
    echo "  stop      Stop ML_Elec"
    echo "  restart   Restart ML_Elec"
    echo "  logs      View logs"
    echo "  status    Check service status"
    echo "  clean     Remove all containers and volumes"
    echo "  help      Show this help message"
    echo ""
}

# Clean deployment
clean() {
    log_warning "This will remove all containers and volumes!"
    read -p "Are you sure? (y/N) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker compose down -v
        log_success "Cleanup complete"
    else
        log_info "Cleanup cancelled"
    fi
}

# Parse arguments
case "${1:-deploy}" in
    deploy)
        check_prerequisites
        setup_directories
        setup_env
        deploy
        ;;
    start)
        docker compose up -d
        ;;
    stop)
        docker compose down
        ;;
    restart)
        docker compose restart
        ;;
    logs)
        docker compose logs -f
        ;;
    status)
        docker compose ps
        ;;
    clean)
        clean
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        log_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac
