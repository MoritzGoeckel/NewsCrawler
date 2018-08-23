minikube start --network-plugin=cni --vm-driver=virtualbox

minikube ssh <<'ENDSSH'
mkdir data
mkdir data/elastic-volume 
chmod 777 data/elastic-volume
ENDSSH