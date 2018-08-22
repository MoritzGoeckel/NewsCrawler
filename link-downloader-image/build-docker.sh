docker rmi moritzgoeckel/news-service:downloader

docker build -t moritzgoeckel/news-service:downloader .
docker push moritzgoeckel/news-service:downloader 
