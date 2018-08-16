minikube delete
./scripts/start_minkube.sh
kubectl create -f agt-redis.yml
kubectl create -f pq-redis.yml
kubectl create -f downloader.yml
