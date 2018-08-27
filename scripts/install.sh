kubectl apply -f agt-article-redis.yml
kubectl apply -f agt-link-redis.yml

kubectl apply -f pq-redis.yml
kubectl apply -f lq-redis.yml

kubectl apply -f cache-redis.yml
kubectl apply -f mongodb.yml
kubectl apply -f elastic.yml

kubectl apply -f article-downloader.yml
kubectl apply -f link-downloader.yml
kubectl apply -f processor.yml
kubectl apply -f static-content-server.yml
kubectl apply -f api-server.yml
