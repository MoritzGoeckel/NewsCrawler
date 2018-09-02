docker rmi moritzgoeckel/news-service:api-server

docker build -t moritzgoeckel/news-service:api-server .
docker push moritzgoeckel/news-service:api-server 

kubectl delete deployment api-server
kubectl apply -f ../api-server.yml
