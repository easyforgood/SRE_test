#!groovy

node {
    def TIDB_TEST_BRANCH = "master"
    def TIKV_BRANCH = "master"
    def PD_BRANCH = "master"

    checkout scm
    dir("SRE"){
        git url: 'https://github.com/easyforgood/SRE_test.git'
        if (!fileExists("bin") ){
            sh "mkdir bin"
        }
        sh "ls"
    }
    
    stage("build"){
        stage("build tidb"){
            dir("tidb"){
                checkout scm
                sh "make server"
                sh "cp bin/*  ../SRE/bin/"

            }
        }

        stage("build tikv"){
            dir("tikv"){
             git url: 'https://github.com/pingcap/tikv'
                 sh "make build"
                 sh "cp bin/*  ../SRE/bin/"
            }
        }

        stage("build pd"){
            dir("pd"){
                 git url: 'https://github.com/pingcap/pd'
                 sh "make build"
                 sh "cp bin/*  ../SRE/bin/"
            }
        }
    }

    stage("test"){
            dir("tidb"){
                sh "make test"
            }

            dir("pd"){
                sh "make test"
            }

            dir("tikv"){
                sh "make test"
            }
    }
    stage("create docker images"){
        dir("SRE"){
            docker.build("tidb_test", "-f tidb/Dockerfile .")
            docker.build("tikv_test", "-f tikv/Dockerfile .")
            docker.build("pd_test", "-f pd/Dockerfile .")
        }
    }
    stage("integration test"){
        dir("SRE"){
            sh "docker-compose up -d tidb"
            // sleep
            def integration_test_result = sh (
                script: "go run integration/main.go",
                returnStdout: true
                ).trim()
            echo integration_test_result
            sh "docker-compose down"
        }
    }
}
