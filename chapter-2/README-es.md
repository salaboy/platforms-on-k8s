# Capitulo 2 :: Desafios de una Aplicación Cloud-Native

—
_🌍 Disponible en_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [Português (Portuguese)](README-pt.md) | [Español](README-es.md) | [日本語 (Japanese)](README-ja.md)
> **Nota:** Presentado por la fantástica comunidad de [ 🌟 contribuidores](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) cloud-native!

En este breve tutorial, instalaremos la `Aplicación para Conferencia` usando Helm en un Clúster Kubernetes local de KinD. 

> [!NOTA]
> Los Helm Charts se pueden publicar en repositorios de Helm Charts o también, desde Helm 3.7, como contenedores OCI en registros de contenedores.

## Creando un clúster local con Kubernetes KinD

> [!IMPORTANTE]
> Asegúrese de tener los pre-requisitos para todos los tutoriales. Puedes encontrarlos [aquí](../chapter-1/README.md#pre-requisites-for-the-tutorials).

Usa el siguiente comando para crear un clúster KinD con tres nodos trabajadores y 1 nodo de Control Plane.

```shell
cat <<EOF | kind create cluster --name dev --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
- role: worker
- role: worker
- role: worker
EOF

```

![3 Nodos esclavos](imgs/cluster-topology.png)

## Cargando algunas imágenes de contenedores antes de instalar la aplicación y otros componentes

El script `kind-load.sh` pre-carga, en otras palabras, descarga y carga las imágenes de contenedores que utilizaremos para nuestra aplicación en nuestro clúster KinD.

La idea aquí es optimizar el proceso para nuestro clúster, de modo que cuando instalemos la aplicación, no tengamos que esperar más de 10 minutos mientras se descargan todas las imágenes de contenedores necesarias. Con todas las imágenes ya precargadas en nuestro clúster KinD, la aplicación debería comenzar en alrededor de 1 minuto, que es el tiempo necesario para que PostgreSQL, Redis y Kafka arranquen.

Ahora, carguemos las imágenes necesarias en nuestro clúster KinD.

> [!Importante]
> Al ejecutar el script mencionado en el siguiente paso, obtendrá todas las imágenes requeridas y luego las cargará en cada nodo de su clúster KinD. Si está ejecutando los ejemplos en un proveedor de nube, esto podría no valer la pena, ya que los proveedores de nube con conexiones de Gigabyte a registros de contenedores podrían obtener estas imágenes en cuestión de segundos.

En su terminal, acceda al directorio `chapter-2`, y desde allí, ejecute el script:

```shell
./kind-load.sh
```

> [!Nota]
> Si está ejecutando Docker Desktop en macOS y ha establecido un tamaño más pequeño para el disco virtual, puede encontrar el siguiente error:
>
>
> ```shell
> $ ./kind-load.sh
> ...
> Command Output: Error response from daemon: write /var/lib/docker/...
> /layer.tar: no space left on device
> ```
>
> Puede modificar el valor del límite del disco virtual en el menú ``Settings -> Resources``.
>   ![Límites del disco virtual de Docker Desktop de MacOS](imgs/macos-docker-desktop-virtual-disk-setting.png)

### Instalando NGINX Ingress Controller
Necesitamos NGINX Ingress Controller para enrutar el tráfico desde nuestra computadora portátil a los servicios que se ejecutan dentro del clúster. NGINX Ingress Controller actúa como un enrutador que se ejecuta dentro del clúster pero también está expuesto al mundo exterior.

```shell
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/release-1.8/deploy/static/provider/kind/deploy.yaml
```
Verifique que los pods dentro del `ingress-nginx` se hayan iniciado correctamente antes de continuar:

```shell
> kubectl get pods -n ingress-nginx
NAME                                        READY   STATUS      RESTARTS   AGE
ingress-nginx-admission-create-cflcl        0/1     Completed   0          62s
ingress-nginx-admission-patch-sb64q         0/1     Completed   0          62s
ingress-nginx-controller-5bb6b499dc-7chfm   0/1     Running     0          62s
```
Esto debería permitirte dirigir el tráfico desde `http://localhost` a los servicios dentro del clúster. Observa que para que KinD funcione de esta manera, proporcionamos parámetros y etiquetas adicionales para el nodo Control Plane cuando creamos el clúster:

```yaml
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true" #Esto permite que el Ingress Controller se instale en el nodo Control Plane
  extraPortMappings:
  - containerPort: 80 # Esto nos permite vincular el puerto 80 en el host local al Ingress Controller, para que pueda dirigir el tráfico a los servicios que se ejecutan dentro del clúster.
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
```
Una vez que tenemos nuestro clúster y nuestro Ingress Controller instalado y configurado, podemos seguir adelante para instalar nuestra aplicación.

## Instalando la Aplicación para Conferencia
Desde Helm 3.7+, podemos usar imágenes OCI para publicar, descargar e instalar Helm Charts. Este enfoque utiliza Docker Hub como un registro de Helm Charts.

Para instalar la Aplicación para Conferencia, solo necesitas ejecutar el siguiente comando:

```shell
helm install conference oci://docker.io/salaboy/conference-app --version v1.0.0
```
También puedes ejecutar el siguiente comando para ver los detalles del Helm Chart:

```shell
helm show all oci://docker.io/salaboy/conference-app --version v1.0.0
```
Verifica que todas los pods de la aplicación estén en funcionamiento.

> [!NOTA]
> Ten en cuenta que si tu conexión a Internet es lenta, puede llevar un tiempo que la aplicación se inicie. Dado que los servicios de la aplicación dependen de algunos componentes de infraestructura (Redis, Kafka, PostgreSQL), estos componentes deben iniciarse y estar listos para que los servicios se conecten.
>
>Componentes como Kafka son bastante pesados, con alrededor de 335+ MB, PostgreSQL 88+ MB y Redis 35+ MB.

Eventualmente, deberías ver algo como esto. Puede tomar unos minutos:

```shell
kubectl get pods
NAME                                                           READY   STATUS    RESTARTS      AGE
conference-agenda-service-deployment-7cc9f58875-k7s2x          1/1     Running   4 (45s ago)   2m2s
conference-c4p-service-deployment-54f754b67c-br9dg             1/1     Running   4 (65s ago)   2m2s
conference-frontend-deployment-74cf86495-jthgr                 1/1     Running   4 (56s ago)   2m2s
conference-kafka-0                                             1/1     Running   0             2m2s
conference-notifications-service-deployment-7cbcb8677b-rz8bf   1/1     Running   4 (47s ago)   2m2s
conference-postgresql-0                                        1/1     Running   0             2m2s
conference-redis-master-0                                      1/1     Running   0             2m2s
```

La columna RESTARTS del pod muestra que quizás Kafka fue lento y el servicio se inicio primero por Kubernetes, por lo tanto, se reinicio para esperar a que Kafka estuviera listo.

Ahora puedes dirigir tu navegador a `http://localhost` para ver la aplicación.

![conference app](imgs/conference-app-homepage.png)

## [Importante] Limpieza - ¡¡¡Debes LEER!!!
Dado que la Aplicación para Conferencia está instalando PostgreSQL, Redis y Kafka, si deseas eliminar e instalar la aplicación nuevamente (lo cual haremos a medida que avancemos en las guías), debes asegurarte de eliminar los correspondientes PersistenceVolumeClaims (PVCs).

Estos PVCs son los volúmenes utilizados para almacenar los datos de las bases de datos y Kafka. No eliminar estos PVCs entre instalaciones hará que los servicios utilicen credenciales antiguas para conectarse a las nuevas bases de datos provisionadas.

Puedes eliminar todos los PVCs enumerándolos con:

```shell
kubectl get pvc
```

Deberías ver:

```shell
NAME                                   STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-conference-kafka-0                Bound    pvc-2c3ccdbe-a3a5-4ef1-a69a-2b1022818278   8Gi        RWO            standard       8m13s
data-conference-postgresql-0           Bound    pvc-efd1a785-e363-462d-8447-3e48c768ae33   8Gi        RWO            standard       8m13s
redis-data-conference-redis-master-0   Bound    pvc-5c2a96b1-b545-426d-b800-b8c71d073ca0   8Gi        RWO            standard       8m13s
```

Y luego eliminar con:
```shell
kubectl delete pvc  data-conference-kafka-0 data-conference-postgresql-0 redis-data-conference-redis-master-0
```
El nombre de los PVCs cambiará según el nombre de la versión de Helm que usaste al instalar la gráfica.

Finalmente, si deseas deshacerte completamente del clúster KinD, puedes ejecutar:

```shell
kind delete clusters dev
```
-------

## Siguientes Pasos
Te recomiendo encarecidamente que te ensucies las manos con un clúster de Kubernetes real alojado en un proveedor de servicios en la nube. Puedes probar la mayoría de los proveedores de servicios en la nube, ya que ofrecen una prueba gratuita donde puedes crear clústeres de Kubernetes y ejecutar todos estos ejemplos [consulta este repositorio](https://github.com/learnk8s/free-kubernetes) para obtener más información.

Si puedes crear un clúster en un proveedor de servicios en la nube y poner en marcha la aplicación, obtendrás experiencia práctica en todos los temas tratados en el Capítulo 2.

## Resumen y Contribución
En este breve tutorial, logramos instalar el esqueleto de la Aplicación para Conferencia. Utilizaremos esta aplicación como ejemplo a lo largo del resto de los capítulos. Asegúrate de que esta aplicación funcione para ti, ya que cubre lo básico de usar e interactuar con un clúster de Kubernetes.

¿Quieres mejorar este tutorial? Crea un [issue](https://github.com/salaboy/platforms-on-k8s/issues/new), envíame un mensaje en [Twitter](https://twitter.com/salaboy) o envía un [Pull Request](https://github.com/salaboy/platforms-on-k8s/compare).
