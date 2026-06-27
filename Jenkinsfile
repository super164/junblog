pipeline {
    agent any

    environment {
        SERVER_IP = '39.96.91.255'
        IMAGE_TAG = "${env.BUILD_NUMBER}"
        NETWORK = 'junblog-network'
    }

    options {
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
    }

    stages {
        stage('Checkout') {
            steps {
                echo '检出代码...'
                checkout scm
            }
        }

        stage('Ensure Network') {
            steps {
                echo '确保 Docker 网络存在...'
                sh "docker network create ${NETWORK} || true"
            }
        }

        stage('Start Infrastructure') {
            steps {
                echo '启动 MySQL 和 Redis...'
                // 加载环境变量
                script {
                    env.MYSQL_ROOT_PASSWORD = sh(script: 'grep MYSQL_ROOT_PASSWORD deploy/.env | cut -d= -f2', returnStdout: true).trim()
                    env.MYSQL_PASSWORD = sh(script: 'grep MYSQL_PASSWORD deploy/.env | grep -v ROOT | cut -d= -f2', returnStdout: true).trim()
                }
                sh '''
                    # MySQL - 如果没运行就启动
                    if [ "$(docker inspect -f '{{.State.Running}}' junblog-mysql 2>/dev/null)" != "true" ]; then
                        docker rm -f junblog-mysql 2>/dev/null || true
                        docker run -d \
                            --name junblog-mysql \
                            --network ${NETWORK} \
                            --restart always \
                            -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
                            -e MYSQL_DATABASE=junblog \
                            -e MYSQL_USER=junblog \
                            -e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
                            -e TZ=Asia/Shanghai \
                            -v mysql_data:/var/lib/mysql \
                            -v $(pwd)/deploy/init.sql:/docker-entrypoint-initdb.d/init.sql \
                            -p 3306:3306 \
                            mysql:8.0 \
                            --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci
                    fi

                    # Redis - 如果没运行就启动
                    if [ "$(docker inspect -f '{{.State.Running}}' junblog-redis 2>/dev/null)" != "true" ]; then
                        docker rm -f junblog-redis 2>/dev/null || true
                        docker run -d \
                            --name junblog-redis \
                            --network ${NETWORK} \
                            --restart always \
                            -v redis_data:/data \
                            -p 6379:6379 \
                            redis:7-alpine \
                            redis-server --appendonly yes
                    fi

                    # 等待 MySQL 就绪
                    echo "等待 MySQL 就绪..."
                    for i in $(seq 1 30); do
                        if docker exec junblog-mysql mysqladmin ping -h localhost 2>/dev/null; then
                            echo "MySQL 已就绪"
                            break
                        fi
                        echo "等待 MySQL... ($i/30)"
                        sleep 2
                    done
                '''
            }
        }

        stage('Build Backend') {
            steps {
                echo '构建后端镜像...'
                sh """
                    cd blog_backend
                    docker build -t junblog-backend:${IMAGE_TAG} .
                    docker tag junblog-backend:${IMAGE_TAG} junblog-backend:latest
                """
            }
        }

        stage('Build Frontend') {
            steps {
                echo '构建前端镜像...'
                sh """
                    cd bolg_forntend
                    docker build -t junblog-frontend:${IMAGE_TAG} .
                    docker tag junblog-frontend:${IMAGE_TAG} junblog-frontend:latest
                """
            }
        }

        stage('Deploy Backend') {
            steps {
                echo '部署后端...'
                sh '''
                    docker rm -f junblog-backend 2>/dev/null || true
                    docker run -d \
                        --name junblog-backend \
                        --network ${NETWORK} \
                        --restart always \
                        -v $(pwd)/deploy/uploads:/app/uploads \
                        -v $(pwd)/deploy/logs:/app/logs \
                        -v $(pwd)/blog_backend/configs/config.docker.yaml:/app/configs/config.yaml \
                        -p 8080:8080 \
                        -e TZ=Asia/Shanghai \
                        junblog-backend:latest
                '''
            }
        }

        stage('Deploy Frontend') {
            steps {
                echo '部署前端...'
                sh '''
                    docker rm -f junblog-frontend 2>/dev/null || true
                    docker run -d \
                        --name junblog-frontend \
                        --network ${NETWORK} \
                        --restart always \
                        -p 80:80 \
                        junblog-frontend:latest
                '''
            }
        }

        stage('Health Check') {
            steps {
                echo '健康检查...'
                sh '''
                    echo "等待服务启动..."
                    sleep 15

                    # 检查后端健康状态
                    for i in $(seq 1 5); do
                        if curl -sf http://${SERVER_IP}:8080/api/v1/health > /dev/null; then
                            echo "后端服务启动成功！"
                            break
                        else
                            echo "等待后端服务启动... ($i/5)"
                            sleep 5
                        fi
                    done

                    # 检查前端是否可访问
                    if curl -sf http://${SERVER_IP} > /dev/null; then
                        echo "前端服务启动成功！"
                    else
                        echo "前端服务启动失败！"
                        exit 1
                    fi
                '''
            }
        }

        stage('Clean Old Images') {
            steps {
                echo '清理旧镜像...'
                sh 'docker image prune -f || true'
            }
        }
    }

    post {
        success {
            echo '========================================='
            echo '部署成功！'
            echo "前端地址: http://${SERVER_IP}"
            echo "后端 API: http://${SERVER_IP}:8080/api/v1"
            echo '========================================='
        }
        failure {
            echo '========================================='
            echo '部署失败！请检查日志。'
            echo '========================================='
            sh 'docker logs junblog-backend --tail=50 || true'
            sh 'docker logs junblog-frontend --tail=50 || true'
        }
        always {
            cleanWs()
        }
    }
}
