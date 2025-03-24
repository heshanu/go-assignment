##first install minikube

minikube start

//I changed dockerhub image,because without public,this cannot create 
//deployments and services

kubectl create deployment golangassignment --image=hecha/goassignment
kubectl expose deployment golangassignment --type=NodePort --port=8081

##And run project ,start goassignment service

minikube service golangassignment

##For dashboard
minikube dashboard


############
I create services,deployments using yml files,they are inside k8s folder