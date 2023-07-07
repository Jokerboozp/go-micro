
![e923936c265ad2a3568e567774564b06.jpeg](https://i2.mjj.rip/2023/07/05/e923936c265ad2a3568e567774564b06.jpeg)

- chi
```shell
go get github.com/go-chi/chi/v5
go get github.com/go-chi/chi/v5/middleware
go get github.com/go-chi/cors
```
- grpc
```shell
go get google.golang.org/grpc
go get google.golang.org/protobuf
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto 
```

- docker
```shell
docker build -f logger-service.dockerfile -t jokerboozp/logger-service:1.0.0 .
docker push jokerboozp/logger-service:1.0.0
docker swarm init
docker stack deploy -c swarm.yml myapp
docker service ls
docker service scale myapp_logger-service=2
docker service update --image jokerboozp/logger-service:1.0.1 myapp_logger-service
docker stack rm myapp
docker swarm leave --force
make build_front_linux
docker build -f front-end.dockerfile -t jokerboozp/front-end:1.0.0 .
docker push jokerboozp/front-end:1.0.0
```

- minikube
```shell
brew install minikube
brew install kubectl
minikube start --nodes=2
minikube start
minikube stop
minikube status
minikube dashboard
kubectl get pods
kubectl get svc
kubectl apply -f k8s
kubectl get deployments
kubectl delete deployments broker-service mongo rabbitmq
kubectl delete svc broker-service mongo rabbitmq
kubectl logs broker-service-b45b98fdb-lq7hp
kubectl expose deployment broker-service --type=LoadBalancer --port=8080 --target-port=8080
minikube tunnel
minikube addons enable ingress
```