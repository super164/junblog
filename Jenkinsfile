pipeline {
    agent any

    environment {
        DOCKER_COMPOSE_FILE = 'deploy/docker-compose.yml'
        SERVER_IP = '39.96.91.255'
        IMAGE_TAG = "${env.BUILD_NUMBER}"
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

        stage('Build Backend') {
            steps {
                echo '构建后端镜像...'
                sh '''
                    cd blog_backend
                    docker build -t junblog-backend:${IMAGE_TAG} .
                    docker tag junblog-backend:${IMAGE_TAG} junblog-backend:latest
                '''
            }
        }

        stage('Build Frontend') {
            steps {
                echo '构建前端镜像...'
                sh '''
                    cd bolg_forntend
                    docker build -t junblog-frontend:${IMAGE_TAG} .
                    docker tag junblog-frontend:${IMAGE_TAG} junblog-frontend:latest
                '''
            }
        }

        stage('Stop Old Containers') {
            steps {
                echo '停止旧容器...'
                sh '''
                    cd deploy || true
                    docker-compose down || true
                '''
            }
        }

        stage('Deploy') {
            steps {
                echo '部署新版本...'
                sh '''
                    cd deploy
                    docker-compose up -d --build
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
                        if curl -sf http://${SERVER_IP}/api/v1/health > /dev/null; then
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
                sh '''
                    # 清除悬空镜像
                    docker image prune -f || true
                '''
            }
        }
    }

    post {
        success {
            echo '========================================='
            echo '部署成功！'
            echo "前端地址: http://${SERVER_IP}"
            echo "后端 API: http://${SERVER_IP}/api/v1"
            echo '========================================='
        }
        failure {
            echo '========================================='
            echo '部署失败！请检查日志。'
            echo '========================================>'
            sh 'docker-compose logs --tail=50'
        }
        always {
            cleanWs()
        }
    }
}
