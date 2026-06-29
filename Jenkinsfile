pipeline {
    agent any

    environment {
        SERVER_IP = '39.96.91.255'
        SERVER_USER = 'root'
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
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: '*/master']],
                    userRemoteConfigs: [[
                        url: 'https://github.com/super164/junblog.git',
                        credentialsId: 'github-token'
                    ]]
                ])
            }
        }

        stage('Build Backend') {
            steps {
                echo '编译后端...'
                bat '''
                    cd blog_backend
                    set GOPROXY=https://goproxy.cn,direct
                    set CGO_ENABLED=0
                    set GOOS=linux
                    set GOARCH=amd64
                    go build -o server.exe ./cmd/server
                '''
            }
        }

        stage('Deploy') {
            steps {
                echo '部署到服务器...'
                sshPublisher(publishers: [
                    sshPublisherDesc(
                        configName: 'deploy-server',
                        transfers: [
                            sshTransfer(
                                sourceFiles: 'blog_backend/server.exe',
                                remoteDirectory: '/opt/junblog/blog_backend/',
                                flatten: true
                            ),
                            sshTransfer(
                                sourceFiles: 'blog_backend/Dockerfile, blog_backend/configs/**',
                                remoteDirectory: '/opt/junblog/blog_backend/'
                            ),
                            sshTransfer(
                                sourceFiles: 'bolg_forntend/**',
                                remoteDirectory: '/opt/junblog/bolg_forntend/'
                            ),
                            sshTransfer(
                                sourceFiles: 'deploy/docker-compose.yml, deploy/init.sql',
                                remoteDirectory: '/opt/junblog/deploy/'
                            ),
                            sshTransfer(
                                remoteDirectory: '/opt/junblog/deploy/',
                                execCommand: '''
                                    cd /opt/junblog/blog_backend
                                    mv server.exe server
                                    chmod +x server

                                    cd /opt/junblog/blog_backend
                                    docker build -t junblog-backend:latest .

                                    cd /opt/junblog/bolg_forntend
                                    docker build -t junblog-frontend:latest .

                                    cd /opt/junblog/deploy
                                    docker compose down
                                    docker compose up -d
                                '''
                            )
                        ]
                    )
                ])
            }
        }

        stage('Health Check') {
            steps {
                echo '健康检查...'
                bat '''
                    timeout /t 15 /nobreak
                    curl -sf http://%SERVER_IP%:8080/api/v1/health && echo "后端OK" || echo "后端启动中..."
                    curl -sf http://%SERVER_IP% && echo "前端OK" || echo "前端启动中..."
                '''
            }
        }
    }

    post {
        success {
            echo '========================================='
            echo '部署成功！'
            echo "前端: http://${SERVER_IP}"
            echo "后端: http://${SERVER_IP}:8080/api/v1"
            echo '========================================='
        }
        failure {
            echo '========================================='
            echo '部署失败！请检查日志。'
            echo '========================================='
        }
    }
}
