minikube ssh <<'ENDSSH'
su
rm data -r

mkdir data
mkdir data/elastic-volume 
chmod 777 data/elastic-volume

logout
ENDSSH
