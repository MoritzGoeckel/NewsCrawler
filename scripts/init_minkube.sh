minikube delete

minikube start --network-plugin=cni

minikube ssh <<'ENDSSH'
mkdir /data/elastic-volume
chmod 777 /data/elastic-volume
logout
ENDSSH
