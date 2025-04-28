pipeline {
    agent any
    
    environment {
        // AWS Configuration
        586794481217.dkr.ecr.us-east-1.amazonaws.com/eks-demo/golang-app
        AWS_REGION = 'us-east-1'                                                          // Replace with your AWS region
        AWS_ACCOUNT_ID = '586794481217'                                                   // Replace with your AWS account ID
        ECR_REPOSITORY = "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/eks-demo/golang-app"   // Replace with your ECR repo name
        
        // Kubernetes Configuration
        EKS_CLUSTER_NAME = 'my-eks-cluster'                                               // Replace with your EKS cluster name
        K8S_NAMESPACE = 'default'                                                // Replace with your Kubernetes namespace
        K8S_DEPLOYMENT_NAME = 'mygolang-app'                                                    // Replace with your deployment name
        K8S_CONTAINER_NAME = 'mygolang-app-container'                                           // Replace with your container name
        
        // Application Configuration
        GIT_COMMIT_SHORT = sh(script: "git rev-parse --short HEAD", returnStdout: true).trim()
        GIT_BRANCH = sh(script: "git rev-parse --abbrev-ref HEAD", returnStdout: true).trim()
        APP_VERSION = "${GIT_BRANCH}-${GIT_COMMIT_SHORT}"
        DOCKER_IMAGE_TAG = "${ECR_REPOSITORY}:${APP_VERSION}"
        // Also create a latest tag for the current branch
        DOCKER_LATEST_TAG = "${ECR_REPOSITORY}:${GIT_BRANCH}-latest"
    
        
        // Credentials
        // DOCKER_CREDENTIALS = credentials('docker-credentials')                            // Optional: if needed for private base images
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo 'Source code checkout complete.'
            }
        }
        
        stage('Build and Push Docker Image') {
            steps {
                script {
                    echo 'Building and pushing Docker image...'
                    
                    // Authenticate with AWS ECR
                    sh "aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
                    
                    // Build the Docker image
                    sh "docker build -t ${DOCKER_IMAGE_TAG} ."
                    
                    // Push to ECR
                    sh "docker push ${DOCKER_IMAGE_TAG}"
                    
                    echo "Docker image pushed: ${DOCKER_IMAGE_TAG}"
                }
            }
        }
        
        stage('Deploy to EKS') {
            steps {
                script {
                    echo 'Deploying to EKS...'
                    
                    // Update kubectl context for the specified EKS cluster
                    sh "aws eks update-kubeconfig --name ${EKS_CLUSTER_NAME} --region ${AWS_REGION}"
                    
                    // Check if namespace exists, create if not
                    sh """
                    if ! kubectl get namespace ${K8S_NAMESPACE} &> /dev/null; then
                        kubectl create namespace ${K8S_NAMESPACE}
                        echo "Created namespace: ${K8S_NAMESPACE}"
                    fi
                    """
                    
                    // Method 1: Update deployment with new image
                    sh "kubectl set image deployment/${K8S_DEPLOYMENT_NAME} ${K8S_CONTAINER_NAME}=${DOCKER_IMAGE_TAG} -n ${K8S_NAMESPACE}"
                    
                    // Method 2: Alternative - Apply Kubernetes manifests from files
                    // Uncomment the below block if using Kubernetes YAML manifests
                    /*
                    // Replace placeholders in Kubernetes manifests
                    sh """
                    sed -i 's|{{IMAGE}}|${DOCKER_IMAGE_TAG}|g' kubernetes/deployment.yaml
                    sed -i 's|{{NAMESPACE}}|${K8S_NAMESPACE}|g' kubernetes/deployment.yaml
                    """
                    
                    // Apply Kubernetes manifests
                    sh "kubectl apply -f kubernetes/deployment.yaml -n ${K8S_NAMESPACE}"
                    sh "kubectl apply -f kubernetes/service.yaml -n ${K8S_NAMESPACE}"
                    */
                    
                    // Wait for deployment to be ready
                    sh "kubectl rollout status deployment/${K8S_DEPLOYMENT_NAME} -n ${K8S_NAMESPACE} --timeout=300s"
                    
                    echo "Deployment to EKS completed successfully."
                }
            }
        }
        
        stage('Verify Deployment') {
            steps {
                script {
                    echo 'Verifying deployment...'
                    
                    // Check pods are running
                    sh "kubectl get pods -l app=${K8S_DEPLOYMENT_NAME} -n ${K8S_NAMESPACE}"
                    
                    // Get service details (if applicable)
                    sh "kubectl get svc -l app=${K8S_DEPLOYMENT_NAME} -n ${K8S_NAMESPACE}"
                    
                    // Optional: Run smoke tests against the deployed application
                    // sh "./scripts/smoke-tests.sh"
                }
            }
        }
    }
    
    post {
        success {
            echo 'Pipeline completed successfully!'
            // Send success notifications
            // slackSend channel: '#deployments', color: 'good', message: "Deployment of ${K8S_DEPLOYMENT_NAME} v${APP_VERSION} to EKS cluster ${EKS_CLUSTER_NAME} was successful!"
        }
        
        failure {
            echo 'Pipeline failed!'
            // Send failure notifications
            // slackSend channel: '#deployments', color: 'danger', message: "Deployment of ${K8S_DEPLOYMENT_NAME} v${APP_VERSION} to EKS cluster ${EKS_CLUSTER_NAME} failed!"
        }
        
        always {
            echo 'Cleaning up workspace...'
            // Clean up any temporary files
            sh 'rm -rf *.tmp'
            
            // Optionally clean Docker images to save space
            // sh "docker rmi ${DOCKER_IMAGE_TAG} || true"
            
            // Archive artifacts
            archiveArtifacts artifacts: 'target/*.jar, build/libs/*.jar', allowEmptyArchive: true
        }
    }
}