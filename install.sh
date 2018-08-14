# Create the static content map
kubectl create configmap config --from-file=config/

kubectl create -f frontend-deployment.yml
kubectl create -f frontend-service.yml

#kubectl apply -f .
