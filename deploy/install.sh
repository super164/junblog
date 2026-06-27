#!/bin/bash

# ========================================
# JunBlog 服务器环境安装脚本
# 适用于 Ubuntu/Debian 系统
# ========================================

set -e

echo "========================================="
echo "JunBlog 服务器环境安装脚本"
echo "========================================="

# 更新系统
echo "[1/5] 更新系统包..."
sudo apt update
sudo apt upgrade -y

# 安装必要的工具
echo "[2/5] 安装必要的工具..."
sudo apt install -y curl wget git vim

# 安装 Docker
echo "[3/5] 安装 Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    rm get-docker.sh
    echo "Docker 安装完成"
else
    echo "Docker 已安装"
fi

# 将当前用户添加到 docker 组
sudo usermod -aG docker $USER
echo "当前用户已添加到 docker 组"

# 安装 Docker Compose
echo "[4/5] 安装 Docker Compose..."
if ! command -v docker-compose &> /dev/null; then
    COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep -o '"tag_name":.*' | cut -d'"' -f4)
    sudo curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo "Docker Compose 安装完成"
else
    echo "Docker Compose 已安装"
fi

# 配置 Docker 镜像加速（可选）
echo "[5/5] 配置 Docker..."
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
  "registry-mirrors": [
    "https://mirror.ccs.tencentyun.com",
    "https://registry.docker-cn.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

sudo systemctl daemon-reload
sudo systemctl enable docker
sudo systemctl restart docker

echo ""
echo "========================================="
echo "安装完成！"
echo "========================================="
echo ""
echo "Docker 版本："
docker --version
echo ""
echo "Docker Compose 版本："
docker-compose --version
echo ""
echo "请重新登录以使 docker 组权限生效，或者运行："
echo "  newgrp docker"
echo ""
echo "部署步骤："
echo "  1. 克隆代码：git clone <your-repo-url>"
echo "  2. 进入部署目录：cd junblog/deploy"
echo "  3. 修改 .env 配置文件"
echo "  4. 启动服务：docker-compose up -d"
echo ""
