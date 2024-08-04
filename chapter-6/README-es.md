# Capítulo 6: Construyamos una plataforma basada en Kubernetes

—
_🌍 Disponible
en_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md) | [Español](README-es.md)
> **Nota:** Presentado por la fantástica comunidad
> de [ 🌟 contribuidores](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) cloud-native!

---

En este tutorial paso a paso, crearemos las APIs de nuestra plataforma reutilizando el poder de las APIs de Kubernetes.
El primer caso de uso en el que la plataforma puede ayudar a los equipos de desarrollo es creando nuevos entornos de
desarrollo y proporcionando un enfoque de autoservicio.

Para construir este ejemplo, utilizaremos Crossplane y `vcluster`, dos proyectos de código abierto alojados en la
Cloud-Native Computing Foundation.

## Instalación

Para instalar Crossplane, necesitas tener un clúster de Kubernetes; puedes crear uno usando KinD como hicimos para
ti. [Capítulo 5](../chapter-5/README-es.md#instalación-de-crossplane)

Usaremos [`vcluster`](https://www.vcluster.com/) en este tutorial, pero no es necesario instalar nada en nuestro clúster
para que `vcluster` funcione. Necesitamos el CLI de `vcluster` para conectarnos a nuestros `vcluster`s. Puedes
instalarlo siguiendo las instrucciones en el sitio
oficial: [https://www.vcluster.com/docs/getting-started/setup](https://www.vcluster.com/docs/getting-started/setup)

## Definiendo nuestra API de Entorno

Un entorno representa un clúster de Kubernetes donde la Aplicación de Conferencias se instalará para desarrollo. La idea
es proporcionar a los equipos entornos de autoservicio para que realicen su trabajo.

Para este tutorial, definiremos una API de Entorno y una Composición de Crossplane que usa el Proveedor de Helm para
crear una nueva instancia de `vcluster`.

Consulta la Definición de Recursos Compuestos (XRD) de Crossplane para
nuestros [Entornos aquí](resources/env-resource-definition.yaml) y la Composición de
[Crossplane aquí](resources/composition-devenv.yaml). Este recurso configura el aprovisionamiento de un nuevo `vcluster`
usando el Proveedor de Helm de
Crossplane, [consulta esta configuración aquí](resources/compositions/composition-devenv.yaml). Cuando se crea un
nuevo `vcluster`, la composición instala nuestra Aplicación de Conferencias en él, nuevamente utilizando el Proveedor de
Helm de Crossplane, pero esta vez
configurado [para apuntar a las APIs del `vcluster` recién creado](resources/compositions/composition-devenv.yaml),
puedes [consultar esto aquí](resources/compositions/composition-devenv.yaml).

Vamos a instalar ambos XRD ejecutando:

```shell
kubectl apply -f resources/definitions
```

Ahora que el XRD está definido, instalemos la Composición ejecutando:

```shell
kubectl apply -f resources/compositions
```

Deberías ver:

```shell
composition.apiextensions.crossplane.io/dev.env.salaboy.com created
compositeresourcedefinition.apiextensions.crossplane.io/environments.salaboy.com created
```

Con el recurso de Entorno y la Composición de Crossplane usando `vcluster`, nuestros equipos ahora pueden solicitar sus
Entornos bajo demanda.

## Solicitar un nuevo Entorno

Para solicitar un nuevo Entorno, los equipos pueden crear nuevos recursos de entorno como este:

```yml
apiVersion: salaboy.com/v1alpha1
kind: Environment
metadata:
  name: team-a-dev-env
spec:
  compositionSelector:
    matchLabels:
      type: development
  parameters:
    installInfra: true

```

Una vez enviado al clúster, la Composición de Crossplane entrará en acción y creará un nuevo `vcluster` con una
instancia de la Aplicación de Conferencias dentro.

```shell
kubectl apply -f team-a-dev-env.yaml
```

Deberías ver:

```shell
environment.salaboy.com/team-a-dev-env created
```

Siempre puedes verificar el estado de tus Entornos ejecutando:

```shell
> kubectl get env
NAME             CONNECT-TO             TYPE          INFRA   DEBUG   SYNCED   READY   CONNECTION-SECRET   AGE
team-a-dev-env   team-a-dev-env-jp7j4   development   true    true    True     False   team-a-dev-env      1s

```

Puedes verificar que Crossplane está creando y gestionando los recursos relacionados con la composición ejecutando:

```shell
> kubectl get managed
NAME                            CHART            VERSION          SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
team-a-dev-env-jp7j4-8lbtj      conference-app   v1.0.0           True     True    deployed   1          Install complete   57s
team-a-dev-env-jp7j4-vcluster   vcluster         0.15.0-alpha.0   True     True    deployed   1          Install complete   57s
```

Estos recursos gestionados no son otros que los lanzamientos de Helm que se están creando:

```shell
kubectl get releases
NAME                            CHART            VERSION          SYNCED   READY   STATE      REVISION   DESCRIPTION        AGE
team-a-dev-env-jp7j4-8lbtj      conference-app   v1.0.0           True     True    deployed   1          Install complete   45s
team-a-dev-env-jp7j4-vcluster   vcluster         0.15.0-alpha.0   True     True    deployed   1          Install complete   45s
```

Luego, podemos conectarnos al entorno aprovisionado ejecutando (utiliza la columna CONNECT-TO para el nombre
del `vcluster`):

```shell
vcluster connect team-a-dev-env-jp7j4 --server https://localhost:8443 -- zsh
```

Una vez que estés conectado al `vcluster`, estarás en un clúster de Kubernetes diferente, por lo que si listas todos los
namespaces disponibles, deberías ver:

```shell
kubectl get ns
NAME              STATUS   AGE
default           Active   64s
kube-system       Active   64s
kube-public       Active   64s
kube-node-lease   Active   64s
```

Como puedes ver, Crossplane no está instalado aquí. Pero si listas todos los pods en este clúster, deberías ver todos
los pods de la aplicación en ejecución:

```shell
NAME                                                              READY   STATUS    RESTARTS      AGE
conference-app-kafka-0                                            1/1     Running   0             103s
conference-app-postgresql-0                                       1/1     Running   0             103s
conference-app-c4p-service-deployment-57d4ddcd68-45f6h            1/1     Running   2 (99s ago)   104s
conference-app-agenda-service-deployment-9bf7946c9-mmx8h          1/1     Running   2 (98s ago)   104s
conference-app-redis-master-0                                     1/1     Running   0             103s
conference-app-frontend-deployment-c8c64c54d-lntnw                1/1     Running   2 (98s ago)   104s
conference-app-notifications-service-deployment-64ff7bcdf8nbvhl   1/1     Running   3 (80s ago)   104s
```

También puedes hacer un port-forwarding a este clúster para acceder a la aplicación usando:

```shell
kubectl port-forward svc/frontend 8080:80
```

Ahora tu aplicación está disponible en [http://localhost:8080](http://localhost:8080)

Puedes salir del contexto de `vcluster` escribiendo `exit` en el terminal.

## Simplificando la superficie de nuestra plataforma

Podemos dar un paso más para simplificar la interacción con las APIs de la plataforma, evitando que los equipos se
conecten al Clúster de Plataforma y eliminando la necesidad de tener acceso a las APIs de Kubernetes.

En esta breve sección, desplegamos una Interfaz de Usuario de Administración que permite a los equipos solicitar nuevos
entornos a través de un sitio web o un conjunto de REST APIs simplificadas.

Antes de instalar la Interfaz de Usuario de Administración, debes asegurarte de que no estás dentro de una sesión de
`vcluster`. (Puedes salir del contexto de `vcluster` escribiendo `exit` en el terminal). Verifica que tienes los
namespaces `crossplane-system` en el clúster actual al que estás conectado.

Puedes instalar esta Interfaz de Usuario de Administración usando Helm:

```shell
helm install admin oci://docker.io/salaboy/conference-admin --version v1.0.0
```

Una vez instalada, puedes hacer un port-forwarding a la Interfaz de Usuario de Administración ejecutando:

```shell
kubectl port-forward svc/admin 8081:80
```

Ahora puedes crear y verificar tus entornos usando una interfaz simple
en [http://localhost:8081](http://localhost:8081). Si esperas a que el entorno esté
listo, recibirás el comando `vcluster` que debes usar para conectarte al entorno.

![imgs/admin-ui.png](imgs/admin-ui.png)

Al usar esta interfaz simple, los equipos de desarrollo no necesitarán acceder directamente a las APIs de Kubernetes
desde el clúster que tiene todas las herramientas de la plataforma (Por ejemplo: Crossplane y Argo CD).

Además de la interfaz de usuario, la aplicación de Administración de la Plataforma te ofrece un conjunto simplificado de
endpoints REST donde tienes total flexibilidad para definir cómo deben ser los recursos en lugar de seguir el Modelo de
Recursos de Kubernetes. Por ejemplo, en lugar de tener un Recurso de Kubernetes con todos los metadatos necesarios por
la API de Kubernetes, podemos usar el siguiente payload JSON para crear un nuevo Entorno:

```json
{
  "name": "team-curl-dev-env",
  "parameters": {
    "type": "development",
    "installInfra": true,
    "frontend": {
      "debug": true
    }
  }
}
```

Puedes crear este entorno ejecutando:

```shell
curl -X POST -H "Content-Type: application/json" -d @team-a-dev-env-simple.json http://localhost:8081/api/environments/
```

Luego, lista todos los entornos con:

```shell
curl localhost:8081/api/environments/
```

O elimina un entorno ejecutando:

```shell
curl -X DELETE http://localhost:8081/api/environments/team-curl-dev-env
```

Esta aplicación sirve como una fachada entre Kubernetes y el mundo exterior. Dependiendo de las necesidades de tu
organización, es posible que desees tener estas abstracciones (APIs) desde el principio, para que el equipo de
plataforma pueda ajustar sus decisiones sobre herramientas y flujos de trabajo por debajo de la superficie.

## Limpieza

Si deseas deshacerte del clúster de KinD creado para estos tutoriales, puedes ejecutar:

```shell
kind delete clusters dev
```

## Próximos Pasos

¿Puedes extender la Interfaz de Usuario de Administración para crear Bases de Datos y Brokers de Mensajes como hicimos
en el Capítulo 5? ¿Qué se necesitaría? Entender dónde deben hacerse los cambios te dará experiencia práctica en el
desarrollo de componentes que interactúan con las APIs de Kubernetes y proporcionan interfaces simplificadas para los
consumidores.

¿Puedes crear tus propias composiciones para usar Clústeres Reales en lugar de `vcluster`? ¿Para qué tipo de escenario
usarías un Clúster real y cuándo usarías un `vcluster`?

¿Qué pasos adicionales necesitarías seguir para ejecutar esto en un Clúster de Kubernetes real en lugar de ejecutarlo en
Kubernetes KinD?

## Resumen y Contribuciones

En este tutorial, hemos construido una nueva API de Plataforma reutilizando el modelo de Recursos de Kubernetes para
aprovisionar entornos de desarrollo bajo demanda. Además, con la aplicación de Administración de Plataforma hemos creado
una capa simplificada para exponer las mismas capacidades sin presionar a los equipos para que aprendan cómo funciona
Kubernetes o los detalles subyacentes, proyectos y tecnologías que hemos utilizado para construir nuestra Plataforma.

Al basarse en contratos (en este ejemplo, la definición de recurso de Entorno), el equipo de plataforma tiene la
flexibilidad para cambiar los mecanismos utilizados para aprovisionar entornos dependiendo de sus requisitos y
herramientas disponibles.

¿Quieres mejorar este tutorial? Crea un issue, envíame un mensaje en [Twitter](https://twitter.com/salaboy)  o envía un
Pull Request.
