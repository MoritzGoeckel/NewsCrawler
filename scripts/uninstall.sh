kubectl delete deployment agt-article-redis 
kubectl delete deployment agt-link-redis 

kubectl delete deployment pq-redis 
kubectl delete deployment lq-redis 

kubectl delete deployment cache-redis 
kubectl delete deployment mongodb 
kubectl delete deployment elastic 

kubectl delete deployment article-downloader 
kubectl delete cronjob link-downloader 
kubectl delete deployment processor 
kubectl delete deployment static-content-server 
kubectl delete deployment api-server 
kubectl delete cronjob lang-model
kubectl delete cronjob headline-analyzer