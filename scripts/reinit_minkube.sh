minikube delete

minikube start --network-plugin=cni

minikube ssh <<'ENDSSH'
sudo mkdir /data/elastic-volume
sudo chmod 777 /data/elastic-volume
logout
ENDSSH
