docker rmi moritzgoeckel/news-service:api-server

docker build -t moritzgoeckel/news-service:api-server .
docker push moritzgoeckel/news-service:api-server