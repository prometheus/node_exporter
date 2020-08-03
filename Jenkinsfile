library "shared-library@master"

pipeline {
  agent {
    kubernetes {
      idleMinutes 10
      label uniqueLabel('dbadmin')
      // ci container serves as a base image in a CI context to build, test & containerize projects. 
      yaml libraryResource('podtemplates/template.yaml').replace("{container}",'ci').replace("{tag}",'1.2.1')
      defaultContainer 'ci'
    }
  }
  
  stages {
    stage('Setup') {
      steps {
        script {
            stash name: 'buildArtifacts', includes: '**'
            }
          }
        }

    stage('Snyk scan') {
      steps {
        snykscan()
      }
    }

   }
  }
