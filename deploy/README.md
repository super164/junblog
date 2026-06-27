# JunBlog Docker 部署指南

## 前置要求

- Linux 服务器（Ubuntu/Debian 推荐）
- 公网 IP 或域名
- Git

## 快速开始

### 1. 服务器环境安装

```bash
# 克隆代码
git clone <your-repo-url>
cd junblog

# 运行安装脚本
chmod +x deploy/install.sh
./deploy/install.sh

# 重新登录或执行以下命令使 docker 组生效
newgrp docker
```

### 2. 配置环境变量

```bash
cd deploy

# 编辑 .env 文件
vim .env
```

修改以下配置：

```bash
# MySQL root 密码
MYSQL_ROOT_PASSWORD=your_secure_root_password

# MySQL 应用密码
MYSQL_PASSWORD=your_secure_app_password
```

### 3. 配置后端

编辑 `blog_backend/configs/config.docker.yaml`：

```bash
vim ../blog_backend/configs/config.docker.yaml
```

修改以下配置：
- MySQL DSN 中的密码
- JWT Secret
- GitHub OAuth 配置
- CORS 配置（添加你的域名/IP）

### 4. 启动服务

```bash
cd deploy

# 构建并启动所有服务
docker-compose up -d --build

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 5. 验证部署

```bash
# 检查后端健康状态
curl http://39.96.91.255/api/v1/health

# 访问前端
# 浏览器打开 http://39.96.91.255
```

## 服务管理

### 常用命令

```bash
# 查看服务状态
docker-compose ps

# 停止所有服务
docker-compose down

# 重启所有服务
docker-compose restart

# 查看日志
docker-compose logs -f [service_name]

# 进入容器
docker exec -it junblog-backend sh
docker exec -it junblog-frontend sh
```

### 数据库管理

```bash
# 连接 MySQL
docker exec -it junblog-mysql mysql -u junblog -p

# 备份数据库
docker exec junblog-mysql mysqldump -u root -p junblog > backup.sql

# 恢复数据库
docker exec -i junblog-mysql mysql -u root -p junblog < backup.sql
```

## Jenkins 自动化部署

### 安装 Jenkins

```bash
# 使用 Docker 安装 Jenkins
docker run -d \
  --name jenkins \
  -p 8080:8080 \
  -p 50000:50000 \
  -v jenkins_home:/var/jenkins_home \
  -v /var/run/docker.sock:/var/run/docker.sock \
  jenkins/jenkins:lts

# 获取初始管理员密码
docker exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
```

### 配置 Jenkins

1. 访问 `http://39.96.91.255:8080`
2. 安装推荐的插件
3. 创建管理员账号
4. 安装以下插件：
   - Docker Pipeline
   - Git
   - Pipeline

### 创建 Pipeline

1. 新建任务 → Pipeline
2. 配置 Git 仓库地址
3. 选择 Pipeline script from SCM
4. 指定 Jenkinsfile 路径
5. 保存并构建

## 目录结构

```
junblog/
├── blog_backend/          # Go 后端
│   ├── Dockerfile         # 后端 Docker 镜像
│   └── configs/
│       └── config.docker.yaml  # Docker 环境配置
├── bolg_forntend/         # Vue 前端
│   ├── Dockerfile         # 前端 Docker 镜像
│   └── nginx.conf         # Nginx 配置
├── deploy/                # 部署配置
│   ├── docker-compose.yml # Docker Compose 编排
│   ├── .env              # 环境变量
│   ├── init.sql          # 数据库初始化
│   └── install.sh        # 服务器安装脚本
├── Jenkinsfile           # Jenkins Pipeline
└── README.md
```

## 故障排查

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs

# 查看特定服务日志
docker-compose logs backend
docker-compose logs frontend
docker-compose logs mysql
```

### 常见问题

1. **MySQL 连接失败**
   - 检查 MySQL 容器是否正常运行
   - 检查 DSN 配置是否正确

2. **前端无法访问后端 API**
   - 检查 Nginx 配置
   - 检查后端容器是否正常运行

3. **端口被占用**
   ```bash
   # 查看端口占用
   netstat -tlnp | grep :80
   # 修改 docker-compose.yml 中的端口映射
   ```

## 安全建议

1. 修改所有默认密码
2. 配置防火墙，只开放必要端口
3. 使用 HTTPS（配置 SSL 证书）
4. 定期备份数据
5. 定期更新 Docker 镜像
