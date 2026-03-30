pipeline {
    agent any

    tools {
        jdk 'jdk17'
        nodejs 'node18'
    }

    environment {
        NOMBRE_APP = "mi-primera-app"
        VERSION    = "1.0.0"
    }

    stages {

        stage('Checkout') {
            steps {
                checkout scm
                sh 'echo "Código clonado en: $(pwd)"'
                sh 'ls -la'
            }
        }

        stage('Build') {
            steps {
                sh 'echo "Construyendo ${NOMBRE_APP} versión ${VERSION}"'
                sh 'java -version'
            }
        }

        stage('Test') {
            steps {
                sh 'echo "Ejecutando tests..."'
                sh 'echo "Tests completados"'
            }
        }

        stage('Deploy') {
            steps {
                sh 'echo "Desplegando ${NOMBRE_APP}:${VERSION}"'
            }
        }

    }

    post {
        success {
            echo "✅ Pipeline completado — ${NOMBRE_APP}:${VERSION} desplegado"
        }
        failure {
            echo "❌ Pipeline fallido en el stage: ${env.STAGE_NAME}"
        }
        always {
            cleanWs()
        }
    }
}