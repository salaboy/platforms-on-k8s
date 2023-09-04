docker pull bitnami/redis:7.0.11-debian-11-r12
docker pull bitnami/postgresql:15.3.0-debian-11-r17
docker pull bitnami/kafka:3.4.1-debian-11-r0
docker pull registry.k8s.io/ingress-nginx/controller:v1.8.1
docker pull registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20230407
docker pull salaboy/frontend-go-1739aa83b5e69d4ccb8a5615830ae66c:v1.0.0
docker pull salaboy/agenda-service-0967b907d9920c99918e2b91b91937b3:v1.0.0
docker pull salaboy/c4p-service-a3dc0474cbfa348afcdf47a8eee70ba9:v1.0.0
docker pull salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
kind load docker-image -n dev bitnami/redis:7.0.11-debian-11-r12
kind load docker-image -n dev bitnami/postgresql:15.3.0-debian-11-r17
kind load docker-image -n dev bitnami/kafka:3.4.1-debian-11-r0
kind load docker-image -n dev registry.k8s.io/ingress-nginx/controller:v1.8.1
kind load docker-image -n dev registry.k8s.io/ingress-nginx/kube-webhook-certgen:v20230407
kind load docker-image -n dev salaboy/frontend-go-1739aa83b5e69d4ccb8a5615830ae66c:v1.0.0
kind load docker-image -n dev salaboy/agenda-service-0967b907d9920c99918e2b91b91937b3:v1.0.0
kind load docker-image -n dev salaboy/c4p-service-a3dc0474cbfa348afcdf47a8eee70ba9:v1.0.0
kind load docker-image -n dev salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0
