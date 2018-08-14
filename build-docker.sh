docker rmi moritzgoeckel/news-service:static-content-server

docker build -t moritzgoeckel/news-service:static-content-server ./static-content-server-image
docker push moritzgoeckel/news-service:static-content-server 
