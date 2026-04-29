# Go Assignment Kubernetes Deployment#

The application is deployed using Minikube

Minikube is installed and running.
kubectl installed.
Docker is installed and configured to use Minikube's Docker daemon.

## minikube start

kubectl create deployment goassignment-deployment --image=hecha/goassignment:latest

kubectl expose deployment  goassignment-deployment --type=NodePort --port=8081

kubectl port-forward service/goassignment-deployment 7080:8081

Tada! Your application is now available at http://localhost:7080/.


## Prerequisites
I deployed using manifest files,same container with a different port number 8082

### create pod using file
kubectl create -f bookapi-pod.yml

### create deployment using file
kubectl create -f bookapi-service.yml

### create service using file
kubectl create -f book-api-service.yml

after that
minikube service book-api-service

Alternatively, you can access the service using the Minikube IP and the NodePort. For example:

http://<static IP>:30007/books?page=1&limit=3

minikube ip
5. Access the Kubernetes Dashboard
To access the Kubernetes dashboard, run:

minikube dashboard
Manifest Files
The k8s folder contains the Kubernetes manifest files for the deployment and service.

## POSTMAN DOCUMENTS LINK

https://documenter.getpostman.com/view/7287191/2sAYkKGcCv


https://documenter.getpostman.com/view/7287191/2sAYkKGcCx

### screen shots of minikube 
<img src="https://github.com/heshanu/go-assignment/blob/dev/Screenshot%202026-04-29%20233225.png"/>
<img src="https://github.com/heshanu/go-assignment/blob/dev/Screenshot%202026-04-29%20233236.png"/>
<img src="https://github.com/heshanu/go-assignment/blob/dev/Screenshot%202026-04-29%20233253.png"/>

