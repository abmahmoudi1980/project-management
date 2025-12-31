#!/bin/bash

# Project Management - Docker Quick Start Script
# This script helps initialize and run the application with Docker

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$SCRIPT_DIR"

echo "================================================"
echo "Project Management - Docker Setup"
echo "================================================"
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "Visit: https://docs.docker.com/engine/install/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    echo "Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✓ Docker version: $(docker --version)"
echo "✓ Docker Compose version: $(docker compose version)"
echo ""

# Create .env if it doesn't exist
if [ ! -f "$PROJECT_ROOT/.env" ]; then
    echo "📝 Creating .env file from template..."
    cp "$PROJECT_ROOT/.env.docker" "$PROJECT_ROOT/.env"
    
    # Generate random values
    DB_PASSWORD=$(openssl rand -base64 16)
    JWT_SECRET=$(openssl rand -base64 32)
    
    # Update .env with generated values
    sed -i "s/DB_PASSWORD=.*/DB_PASSWORD=$DB_PASSWORD/" "$PROJECT_ROOT/.env"
    sed -i "s/JWT_SECRET=.*/JWT_SECRET=$JWT_SECRET/" "$PROJECT_ROOT/.env"
    
    echo "✓ .env file created with secure defaults"
    echo "⚠️  Please review and update .env as needed"
    echo ""
fi

# Menu
echo "Select an action:"
echo "1) Start services (docker compose up -d)"
echo "2) Build images (docker compose build)"
echo "3) View logs (docker compose logs -f)"
echo "4) Stop services (docker compose down)"
echo "5) Status (docker compose ps)"
echo "6) Reset database (docker compose down -v && docker compose up -d)"
echo "7) Backend shell"
echo "8) Database shell"
echo "9) Exit"
echo ""

read -p "Enter your choice [1-9]: " choice

case $choice in
    1)
        echo "🚀 Starting services..."
        docker compose up -d
        echo ""
        echo "✓ Services started!"
        echo ""
        docker compose ps
        echo ""
        echo "Access the application:"
        echo "  Frontend: http://localhost"
        echo "  API: http://localhost:3000"
        ;;
    2)
        echo "🔨 Building images..."
        docker compose build
        echo "✓ Images built successfully!"
        ;;
    3)
        echo "📋 Showing logs (Ctrl+C to exit)..."
        docker compose logs -f
        ;;
    4)
        echo "🛑 Stopping services..."
        docker compose down
        echo "✓ Services stopped"
        ;;
    5)
        echo "📊 Service status:"
        docker compose ps
        ;;
    6)
        echo "⚠️  This will delete all data!"
        read -p "Are you sure? (yes/no): " confirm
        if [ "$confirm" = "yes" ]; then
            echo "🔄 Resetting database..."
            docker compose down -v
            docker compose up -d
            echo "✓ Database reset completed!"
        else
            echo "Cancelled"
        fi
        ;;
    7)
        echo "📦 Entering backend container shell..."
        docker compose exec backend sh
        ;;
    8)
        echo "🗄️  Entering database shell..."
        docker compose exec postgres psql -U postgres -d project_management
        ;;
    9)
        echo "Goodbye!"
        exit 0
        ;;
    *)
        echo "Invalid choice"
        exit 1
        ;;
esac
