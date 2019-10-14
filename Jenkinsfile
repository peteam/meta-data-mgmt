pipeline {
  agent {
    label "jenkins-go"
  }
  environment {
    ORG = 'peteam-opentest'
    APP_NAME = 'metadata-mgmt-services'
    CHARTMUSEUM_CREDS = credentials('jenkins-x-chartmuseum')
    DOCKER_REGISTRY_ORG = 'platform-engineering'
  }
  stages {
    stage('CI Build and push snapshot') {
      when {
        branch 'PR-*'
      }
      environment {
        PREVIEW_VERSION = "0.0.0-SNAPSHOT-$BRANCH_NAME-$BUILD_NUMBER"
        PREVIEW_NAMESPACE = "$APP_NAME-$BRANCH_NAME".toLowerCase()
        HELM_RELEASE = "$PREVIEW_NAMESPACE".toLowerCase()
      }
      steps {
        container('go') {
          dir('/home/jenkins/go/src/cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services') {
            checkout scm
            sh "go get github.com/securego/gosec/cmd/gosec/..."
            sh "go get github.com/google/go-github/github"
            sh "go get github.com/securego/gosec"
            sh "go get golang.org/x/oauth2"
            sh "go get -u github.com/cweill/gotests/..."
            sh "go get github.com/tebeka/go2xunit/..."

            // Download sonar scanner -- this should be hosted in QP env.
            //sh "wget https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.0.0.1744-linux.zip ; unzip sonar-scanner-cli-4.0.0.1744-linux.zip; chmod +x sonar-scanner-4.0.0.1744-linux/bin/sonar-scanner"
            sh "make linux"

            // Generate Go Unit Test
            sh "/home/jenkins/go/bin/gotests -all . || echo ok"

            // Run Go Unit Test
            sh "GO111MODULE=on go test -v |/home/jenkins/go/bin/go2xunit > test_output.xml"

            // Run Go Security Test
            sh "/home/jenkins/go/bin/gosec -fmt=json -out gosec.json . || echo ok"

            // Run Code Quality Test
            //sh "sonar-scanner-4.0.0.1744-linux/bin/sonar-scanner -Dsonar.projectKey=go_code_check -Dsonar.sources=. -Dsonar.host.url=http://sonarq.devops.quickplay.com -Dsonar.login=11eeaba5aabb5cfa164d7b180a56625c5eed7ba4"
            sh "export VERSION=$PREVIEW_VERSION && skaffold build -f skaffold.yaml"
            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:$PREVIEW_VERSION"
          }
          dir('/home/jenkins/go/src/cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/charts/preview') {
            sh "make preview"
            sh "jx preview --app $APP_NAME --dir ../.."
          }
        }
      }
    }
    stage('Build Release') {
      when {
        branch 'master'
      }
      steps {
        container('go') {
          dir('/home/jenkins/go/src/cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services') {
            checkout scm

            // ensure we're not on a detached head
            sh "git checkout master"
            sh "git config --global credential.helper store"
            sh "jx step git credentials"
            sh "go get github.com/securego/gosec/cmd/gosec/..."
            sh "go get github.com/google/go-github/github"
            sh "go get github.com/securego/gosec"
            sh "go get golang.org/x/oauth2"
            sh "go get -u github.com/cweill/gotests/..."
            sh "go get github.com/tebeka/go2xunit/..."

            // Download sonar scanner -- this should be hosted in QP env.
            //sh "wget https://binaries.sonarsource.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-4.0.0.1744-linux.zip; unzip sonar-scanner-cli-4.0.0.1744-linux.zip; chmod +x sonar-scanner-4.0.0.1744-linux/bin/sonar-scanner"

            // so we can retrieve the version in later steps
            sh "echo \$(jx-release-version) > VERSION"
            sh "jx step tag --version \$(cat VERSION)"
            sh "make build"

            // Generate Go Unit Test
            sh "/home/jenkins/go/bin/gotests -all . || echo ok"

            // Run Go Unit Test
            sh "GO111MODULE=on go test -v |/home/jenkins/go/bin/go2xunit > test_output.xml"

            // Run Go Security Test
            sh "/home/jenkins/go/bin/gosec -fmt=json -out gosec.json . || echo ok"
            //sh "sonar-scanner-4.0.0.1744-linux/bin/sonar-scanner -Dsonar.projectKey=go_code_check -Dsonar.sources=. -Dsonar.host.url=http://sdd-isk8wrkrn-c2-02.quickplay.local:31687 -Dsonar.login=11eeaba5aabb5cfa164d7b180a56625c5eed7ba4"
            sh "export VERSION=`cat VERSION` && skaffold build -f skaffold.yaml"
            sh "jx step post build --image $DOCKER_REGISTRY/$ORG/$APP_NAME:\$(cat VERSION)"
          }
        }
      }
    }
    stage('Promote to Environments') {
      when {
        branch 'master'
      }
      steps {
        container('go') {
          dir('/home/jenkins/go/src/cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/charts/metadata-mgmt-services') {
            sh "jx step changelog --version v\$(cat ../../VERSION)"

            // release the helm chart
            sh "jx step helm release"

            // promote through all 'Auto' promotion Environments
            sh "jx promote -b --all-auto --timeout 1h --version \$(cat ../../VERSION)"
          }
        }
      }
    }
  }
}
