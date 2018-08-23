kubectl create -f agt-article-redis.yml
kubectl create -f agt-link-redis.yml

kubectl create -f pq-redis.yml
kubectl create -f lq-redis.yml

kubectl create -f cache-redis.yml
kubectl create -f mongodb.yml
kubectl create -f elastic.yml

kubectl create -f article-downloader.yml
kubectl create -f link-downloader.yml
kubectl create -f processor.yml
kubectl create -f static-content-server.yml
kubectl create -f api-server.yml