pipeline {
    agent any

    tools {
        jdk 'jdk17'
        nodejs 'node18'
    }

    environment {
        DOCKER_HUB_CREDS = credentials('dockerhub-id')
        DOCKER_IMAGE     = "gjisus/curso-gitops"
        SCANNER_HOME     = tool('sonar-scanner')
        GITHUB_USER      = "GutsNet"
        INFRA_REPO       = "gitops-infra"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: 'git rev-parse --short HEAD',
                        returnStdout: true
                    ).trim()
                    env.BUILD_TAG = "${env.BUILD_NUMBER}-${env.GIT_COMMIT_SHORT}"
                    echo "BUILD_TAG: ${env.BUILD_TAG}"
                }
            }
        }

        stage('SonarQube Analysis') {
            steps {
                withSonarQubeEnv('sonarqube-server') {
                    sh """
                        ${SCANNER_HOME}/bin/sonar-scanner \
                        -Dsonar.projectKey=curso-gitops \
                        -Dsonar.projectName=curso-gitops \
                        -Dsonar.sources=. \
                        -Dsonar.exclusions=**/vendor/**,**/node_modules/**,**/frontend/**
                    """
                }
            }
        }

        stage('Quality Gate') {
            steps {
                timeout(time: 10, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }

        stage('Docker Build') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE}:${BUILD_TAG} ."
                sh "docker tag ${DOCKER_IMAGE}:${BUILD_TAG} ${DOCKER_IMAGE}:latest"
                echo "Imagen construida: ${DOCKER_IMAGE}:${BUILD_TAG}"
            }
        }

        // stage('Trivy Scan') {
        //     steps {
        //         sh """
        //             trivy image \
        //               --exit-code 0 \
        //               --severity HIGH,CRITICAL \
        //               --format table \
        //               ${DOCKER_IMAGE}:${BUILD_TAG}
        //         """
        //     }
        // }

        stage('Docker Push') {
            steps {
                sh "echo ${DOCKER_HUB_CREDS_PSW} | docker login -u ${DOCKER_HUB_CREDS_USR} --password-stdin"
                sh "docker push ${DOCKER_IMAGE}:${BUILD_TAG}"
                sh "docker push ${DOCKER_IMAGE}:latest"
                echo "Imagen subida a Docker Hub: ${DOCKER_IMAGE}:${BUILD_TAG}"
            }
        }

        stage('Deploy to GitOps Repo') {
            steps {
                withCredentials([string(credentialsId: 'github-token-id', variable: 'GITHUB_TOKEN')]) {
                    sh """
                        rm -rf infra-repo

                        git clone https://${GITHUB_TOKEN}@github.com/${GITHUB_USER}/${INFRA_REPO}.git infra-repo

                        cd infra-repo
                        git config user.email "jenkins@local.com"
                        git config user.name "Jenkins CI"

                        # Actualizar el tag de la imagen en deployment.yaml
                        sed -i "s|image: ${DOCKER_IMAGE}:.*|image: ${DOCKER_IMAGE}:${BUILD_TAG}|" \
                            infrastructure/kubernetes/app/deployment.yaml

                        git add infrastructure/kubernetes/app/deployment.yaml
                        git commit -m "ci: deploy version ${BUILD_TAG} from Jenkins"
                        git push origin main
                    """
                }
            }
        }

        stage('Cleanup') {
            steps {
                sh "docker rmi ${DOCKER_IMAGE}:${BUILD_TAG} || true"
                sh "docker rmi ${DOCKER_IMAGE}:latest || true"
                sh "docker image prune -f || true"
            }
        }
    }

    post {
        success {
            echo "✅ Pipeline completado — imagen: ${DOCKER_IMAGE}:${BUILD_TAG}"
        }
        failure {
            echo "❌ Pipeline fallido en stage: ${env.STAGE_NAME}"
        }
        always {
            cleanWs()
        }
    }
}
