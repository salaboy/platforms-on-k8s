# Capítulo 5: Infraestructura multi-nube (App)

—
_🌍 Disponible
en_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md) | [Español](README-es.md)
> **Nota:** Presentado por la fantástica comunidad
> de [ 🌟 contribuidores](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) cloud-native!

---

Este tutorial paso a paso usa Crossplane para aprovisionar las instancias de Redis, PostgreSQL y Kafka para nuestros
servicios de aplicación.

Usando Crossplane y las Composiciones de Crossplane, nuestro objetivo es unificar cómo se aprovisionan estos
componentes, ocultando su ubicación a los usuarios finales (equipos de aplicaciones).

Los equipos de aplicaciones deberían poder solicitar estos recursos utilizando un enfoque declarativo, al igual que con
cualquier otro recurso de Kubernetes. Esto permite a los equipos usar Pipelines de Entorno para configurar tanto los
servicios de aplicación como los componentes de infraestructura de aplicación necesarios para la aplicación.

## Instalación de Crossplane

Para instalar Crossplane, necesitas tener un clúster de Kubernetes; puedes crear uno usando KinD como hicimos para
ti. [Capítulo 2](../chapter-2/README-es.md#creando-un-clúster-local-con-kubernetes-kind).

Instalemos [Crossplane](https://crossplane.io) en su propio namespace utilizando Helm

```shell
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update

helm install crossplane --namespace crossplane-system --create-namespace crossplane-stable/crossplane --version 1.15.0 --wait 
```

Luego, instala el proveedor de Helm de Crossplane, junto con un nuevo `ClusterRoleBinding` para que el Proveedor de Helm
pueda instalar Charts en nuestro nombre.

```shell
kubectl apply -f crossplane/helm-provider.yaml
```

Después de unos segundos, si verificas los proveedores configurados, deberías ver que Helm está `INSTALLED` y `HEALTHY`:

```shell
❯ kubectl get providers.pkg.crossplane.io
NAME            INSTALLED   HEALTHY   PACKAGE                                                    AGE
provider-helm   True        True      xpkg.upbound.io/crossplane-contrib/provider-helm:v0.17.0   49s
```

Luego, crea un `ProviderConfig` que instruya al Proveedor de Helm a usar la configuración del clúster para instalar
Charts
dentro del clúster.

```shell
kubectl apply -f crossplane/helm-provider-config.yaml
```

Ahora estamos listos para instalar nuestras composiciones de Crossplane para Bases de Datos y Brokers de Mensajes para
proporcionar todos los componentes que nuestra aplicación necesita.

## Infraestructura de Aplicaciones bajo demanda usando Composiciones de Crossplane

Necesitamos instalar nuestras Definiciones de Recursos Compuestos de Crossplane (XRDs) para nuestra Base de Datos de
Clave-Valor (Redis), nuestra Base de Datos SQL (PostgreSQL) y nuestro Broker de Mensajes (Kafka).

```shell
kubectl apply -f resources/definitions
```

Ahora instala las composiciones de Crossplane correspondientes y los datos de inicialización:

```shell
kubectl apply -f resources/compositions
kubectl apply -f resources/config
```

El recurso de Composición de Crossplane (`app-database-redis.yaml`) define qué recursos en la nube necesitan ser creados
y cómo deben ser configurados juntos. La Definición de Recursos Compuestos de Crossplane (
XRD) (`app-database-resource.yaml`) define una interfaz simplificada que permite a los equipos de desarrollo de
aplicaciones solicitar rápidamente nuevas bases de datos creando recursos de este tipo.

Consulta el directorio de [recursos/](resources/) para las Composiciones y las Definiciones de Recursos Compuestos (
XRDs).

### Vamos a aprovisionar la Infraestructura de Aplicaciones

Podemos aprovisionar una nueva Base de Datos de Clave-Valor para que nuestro equipo la use ejecutando el siguiente
comando:

```shell
kubectl apply -f my-db-keyvalue.yaml
```

El recurso `my-db-keyvalue.yaml` se ve así:

```yml
apiVersion: salaboy.com/v1alpha1
kind: Database
metadata:
  name: my-db-keyvalue
spec:
  compositionSelector:
    matchLabels:
      provider: local
      type: dev
      kind: keyvalue
  parameters:
    size: small
```

Observa que estamos usando las etiquetas `provider: local`, `type: dev`, y `kind: keyvalue`. Esto permite a Crossplane
encontrar la composición correcta basada en las etiquetas. En este caso, el Proveedor de Helm creó una instancia local
de Redis.

Puedes verificar el estado de la base de datos usando:

```shell
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   97s
```

Puedes verificar que se creó una nueva instancia de Redis en el namespace `default`.

Puedes seguir los mismos pasos para aprovisionar una base de datos PostgreSQL ejecutando:

```shell
kubectl apply -f my-db-sql.yaml
```

Ahora deberías ver dos `dbs`.

```shell
> kubectl get dbs
NAME              SIZE    MOCKDATA   KIND       SYNCED   READY   COMPOSITION                     AGE
my-db-keyavalue   small   false      keyvalue   True     True    keyvalue.db.local.salaboy.com   2m
my-db-sql         small   false      sql        True     False   sql.db.local.salaboy.com        5s
```

Ahora puedes verificar que hay dos Pods en ejecución, uno para cada base de datos:

```shell
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          3m40s
my-db-sql-postgresql-0           1/1     Running   0          104s
```

Debería haber 4 Secrets de Kubernetes (dos para nuestros dos lanzamientos de Helm y dos que contienen las credenciales
para conectarse a las instancias recién creadas):

```shell
> kubectl get secret
NAME                                    TYPE                 DATA   AGE
my-db-keyavalue-redis                   Opaque               1      2m32s
my-db-sql-postgresql                    Opaque               1      36s
sh.helm.release.v1.my-db-keyavalue.v1   helm.sh/release.v1   1      2m32s
sh.helm.release.v1.my-db-sql.v1         helm.sh/release.v1   1      36s
```

Podemos hacer lo mismo para aprovisionar una nueva instancia de nuestro Broker de Mensajes Kafka:

```shell
kubectl apply -f my-messagebroker-kafka.yaml
```

Y luego listar con:

```shell
> kubectl get mbs
NAME          SIZE    KIND    SYNCED   READY   COMPOSITION                  AGE
my-mb-kafka   small   kafka   True     True    kafka.mb.local.salaboy.com   2m51s
```

Kafka no requiere crear ningún secreto al usar su configuración predeterminada.

Deberías ver tres Pods en ejecución (uno para Kafka, uno para Redis y uno para PostgreSQL).

```shell
> kubectl get pods
NAME                             READY   STATUS    RESTARTS   AGE
my-db-keyavalue-redis-master-0   1/1     Running   0          113s
my-db-sql-postgresql-0           1/1     Running   0          108s
my-mb-kafka-0                    1/1     Running   0          100s
```

**Nota**: si estás eliminando y recreando bases de datos o brokers de mensajes usando el mismo nombre de recurso,
recuerda eliminar los `PersistentVolumeClaims`, ya que estos recursos no se eliminan cuando eliminas los recursos de
Database o MessageBroker.

¡Ahora puedes crear tantas instancias de bases de datos o brokers de mensajes como pueda manejar los recursos de tu
clúster!

## Vamos a desplegar nuestra Aplicación de Conferencias

Bien, ahora que tenemos nuestras dos bases de datos y nuestro broker de mensajes en funcionamiento, necesitamos
asegurarnos de que nuestros servicios de aplicación se conecten a estas instancias. Lo primero que debemos hacer es
deshabilitar las dependencias de Helm definidas en el gráfico de la Aplicación de Conferencias para que, cuando se
instale la aplicación, no se instalen las bases de datos ni el broker de mensajes. Podemos hacer esto configurando el
flag `install.infrastructure` en `false`.

Para ello, utilizaremos el archivo `app-values.yaml` que contiene las configuraciones para que los servicios se conecten
a nuestras bases de datos recién creadas:

```shell
helm install conference oci://registry-1.docker.io/salaboy/conference-app --version v1.0.0 -f app-values.yaml
```

El contenido del archivo `app-values.yaml` se ve así:

```yml
install:
  infrastructure: false
frontend:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
agenda:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
  redis:
    host: my-db-keyavalue-redis-master.default.svc.cluster.local
    secretName: my-db-keyavalue-redis
c4p:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
  postgresql:
    host: my-db-sql-postgresql.default.svc.cluster.local
    secretName: my-db-sql-postgresql
notifications:
  kafka:
    url: my-mb-kafka.default.svc.cluster.local
```

Observa que el archivo `app-values.yaml` depende de los nombres que especificamos para nuestras bases de
datos (`my-db-keyvalue` y `my-db-sql`) y nuestros brokers de mensajes (`my-mb-kafka`) en los archivos de ejemplo. Si
solicitas otras bases de datos y brokers de mensajes con nombres diferentes, necesitarás adaptar este archivo con los
nuevos nombres.

Una vez que los pods de la aplicación se inicien, deberías tener acceso a la aplicación al apuntar tu navegador
a [http://localhost:8080](http://localhost:8080). Si has llegado hasta aquí, ahora puedes aprovisionar infraestructura
multi-nube utilizando las Composiciones de Crossplane. Consulta
el [Tutorial de Composiciones de Crossplane para AWS](aws/), que fue aportado
por [@asarenkansah](https://github.com/asarenkansah). Al separar
el aprovisionamiento de la infraestructura de la aplicación del código de la aplicación, no solo habilitas la
portabilidad entre proveedores de nube, sino que también permites a los equipos conectar los servicios de la aplicación
con infraestructura que puede ser gestionada por el equipo de plataforma.

## Limpieza

Si deseas deshacerte del clúster de KinD creado para este tutorial, puedes ejecutar:

```shell
kind delete clusters dev
```

## Próximos Pasos

Si tienes acceso a un Proveedor de Nube como Google Cloud Platform, Microsoft Azure o Amazon AWS, te recomiendo
encarecidamente que revises los **Proveedores de Crossplane** para estas plataformas. Instalar estos proveedores y
aprovisionar Recursos en la Nube, en lugar de usar el Proveedor de Helm de Crossplane, te dará una experiencia real
sobre cómo funcionan estas herramientas.

Como se mencionó en el Capítulo 5, ¿cómo manejarías los servicios que necesitan componentes de infraestructura que no se
ofrecen como servicios gestionados? En el caso de Google Cloud Platform, no ofrecen un Servicio de Kafka Gestionado que
puedas aprovisionar. ¿Instalarías Kafka usando Charts de Helm o VMs, o cambiarías Kafka por un servicio gestionado como
Google PubSub? ¿Mantendrás dos versiones del mismo servicio?

## Resumen y Contribuciones

En este tutorial, hemos logrado separar el aprovisionamiento de la infraestructura de la aplicación del despliegue de la
aplicación. Esto permite a diferentes equipos solicitar recursos bajo demanda (usando composiciones de Crossplane) y
servicios de aplicación que pueden evolucionar de manera independiente.

Usar dependencias de Helm Charts para fines de desarrollo y obtener rápidamente una instancia completamente funcional de
la aplicación en funcionamiento es excelente. Para entornos más sensibles, es posible que desees seguir un enfoque como
el mostrado aquí, donde tienes múltiples maneras de conectar tu aplicación con los componentes requeridos por cada
servicio.

¿Quieres mejorar este tutorial? Crea un issue, envíame un mensaje en [Twitter](https://twitter.com/salaboy)  o envía un
Pull Request.
