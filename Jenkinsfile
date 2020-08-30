pipeline {
    agent armonly
    environment { 
        GOROOT=/usr/local/go
	GOBIN=$GOROOT/bin
        PATH="$PATH:$GOBIN"
    }
    stages {
        stage('Prepare') {
            steps {
		checkout scm
		sh 'make deps'
            }
        }
        stage('Build') {
            steps {
                sh 'make build'
            }
        }
        stage('Test') {
            steps {
                sh 'make test'
            }
        }
        stage('Deploy') {
            steps {
                sh 'make docker'
		sh './run-docker-dev.sh'
            }
        }
    }
}
