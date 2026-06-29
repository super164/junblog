pipeline {
    agent any

    environment {
        SERVER_IP = '39.96.91.255'
        GOPROXY = 'https://goproxy.cn,direct'
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
                echo '编译后端...'
                bat '''
                    cd blog_backend
                    set GOPROXY=https://goproxy.cn,direct
                    set CGO_ENABLED=0
                    set GOOS=linux
                    set GOARCH=amd64
                    go build -o server.exe ./cmd/server
                    dir server.exe
                '''
            }
        }

        stage('Upload Files') {
            steps {
                echo '上传文件到服务器...'
                sshPublisher(publishers: [
                    sshPublisherDesc(
                        configName: 'deploy-server',
                        transfers: [
                            sshTransfer(
                                sourceFiles: 'blog_backend/server.exe',
                                remoteDirectory: 'blog_backend/',
                                flatten: true
                            ),
                            sshTransfer(
                                sourceFiles: 'blog_backend/Dockerfile',
                                remoteDirectory: 'blog_backend/'
                            ),
                            sshTransfer(
                                sourceFiles: 'blog_backend/configs/**',
                                remoteDirectory: 'blog_backend/configs/'
                            ),
                            sshTransfer(
                                sourceFiles: 'bolg_forntend/**',
                                remoteDirectory: 'bolg_forntend/'
                            ),
                            sshTransfer(
                                sourceFiles: 'deploy/docker-compose.yml',
                                remoteDirectory: 'deploy/'
                            )
                        ]
                    )
                ])
            }
        }

        stage('Deploy') {
            steps {
                echo '部署服务...'
                sshPublisher(publishers: [
                    sshPublisherDesc(
                        configName: 'deploy-server',
                        transfers: [
                            sshTransfer(
                                remoteDirectory: '',
                                execCommand: 'bash /opt/junblog/deploy.sh'
                            )
                        ]
                    )
                ])
            }
        }
    }

    post {
        success {
            echo '========================================='
            echo '部署成功！'
            echo "前端: http://${SERVER_IP}"
            echo "后端: http://${SERVER_IP}:8080"
            echo '========================================='
        }
        failure {
            echo '部署失败！'
        }
    }
}
