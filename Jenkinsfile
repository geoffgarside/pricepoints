pipeline {
  agent any

  stages {
    stage('Build') {
      steps {
        sh 'make build'
      }
    }

    stage('Test') {
      steps {
        sh 'make reports/tests.xml reports/coverage.xml'
        junit 'reports/*.xml'
      }
    }

    stage('Lint') {
      steps {
        sh 'make check'
      }
    }
  }
}
