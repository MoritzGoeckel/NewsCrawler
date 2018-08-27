docker rmi moritzgoeckel/news-service:wordcloud-generator

docker build -t moritzgoeckel/news-service:wordcloud-generator .
docker push moritzgoeckel/news-service:wordcloud-generator 
