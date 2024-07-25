# Tekton en Acción

Este breve tutorial cubre cómo instalar Tekton y cómo crear una tarea y un pipeline muy simples.

[Tekton](https://tekton.dev) es un motor de pipelines sin opiniones construido para la nube (específicamente para
Kubernetes). Puedes construir cualquier tipo de pipelines que desees, ya que el motor no impone restricciones sobre el
tipo de tareas que puede ejecutar. Esto lo hace perfecto para construir pipelines de servicios donde podrías tener
requisitos especiales que no pueden ser satisfechos por un servicio gestionado.

Después de ejecutar nuestro primer pipeline de Tekton, este tutorial también incluye enlaces a pipelines de servicios
más complejos utilizados para construir los servicios de la aplicación de conferencia.

## Instalar Tekton

Sigue los siguientes pasos para instalar y configurar Tekton en tu clúster de Kubernetes. Si no tienes un clúster de
Kubernetes, puedes crear uno
con [KinD, como hicimos para el Capítulo 2](../../chapter-2/README-es.md#creando-un-clúster-local-con-kubernetes-kind)

1. **Instalar Tekton Pipelines**

```shell
  kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.45.0/release.yaml
```

1. **Instalar el Dashboard de Tekton (opcional)**

```shell
kubectl apply -f https://github.com/tektoncd/dashboard/releases/download/v0.33.0/release.yaml
```

Puedes acceder al dashboard mediante el port-forwarding utilizando `kubectl`:

```shell
kubectl port-forward svc/tekton-dashboard  -n tekton-pipelines 9097:9097
```

![Tekton Dashboard](imgs/tekton-dashboard.png)

Luego, puedes acceder apuntando tu navegador a [http://localhost:9097](http://localhost:9097)

1. **Instalar el CLI de Tekton (opcional)**:

También puedes instalar [Tekton `tkn` CLI tool](https://github.com/tektoncd/cli).

Si estás en Mac OSX, puedes ejecutar:

```shell
brew install tektoncd-cli
```

## Comenzando con las Tareas de Tekton

Esta sección tiene como objetivo ayudarte a comenzar a crear Tareas y un Pipeline Simple, para que luego puedas explorar
los Pipelines de Servicios utilizados para construir los artefactos para la Aplicación de Conferencia.

Con Tekton, podemos definir lo que hacen nuestras tareas creando definiciones de Tareas de Tekton. A continuación se
muestra el ejemplo más simple de una tarea:

```yml
apiVersion: tekton.dev/v1
kind: Task
metadata:
  name: hello-world-task
spec:
  params:
    - name: name
      type: string
      description: who do you want to welcome?
      default: tekton user
  steps:
    - name: echo
      image: ubuntu
      command:
        - echo
      args:
        - "Hello World: $(params.name)" 
```

Esta `Task` de Tekton utiliza la imagen `ubuntu` y el comando `echo` ubicado dentro de esa imagen. Esta `Task` también
acepta un parámetro llamado `name` que se usará para imprimir un mensaje. Apliquemos esta definición de `Task` a nuestro
clúster ejecutando:

```shell
kubectl apply -f hello-world/hello-world-task.yaml
```

Cuando aplicamos este recurso a Kubernetes, no estamos ejecutando la tarea; solo estamos haciendo que la definición de
la tarea esté disponible para que otros la utilicen. Esta tarea ahora puede ser referenciada en múltiples pipelines o
ejecutada de forma independiente por diferentes usuarios.

Ahora puedes listar las tareas disponibles en el clúster ejecutando:

```shell
> kubectl get tasks
NAME               AGE
hello-world-task   88s
```

Ahora ejecutemos nuestra tarea. Lo hacemos creando un recurso `TaskRun`, que representa una ejecución individual de
nuestra tarea. Ten en cuenta que esta ejecución concreta tendrá un nombre de recurso fijo (`hello-world-task-run-1`) y
un valor concreto para el parámetro de la tarea llamado `name`.

```yml
apiVersion: tekton.dev/v1
kind: TaskRun
metadata:
  name: hello-world-task-run-1
spec:
  params:
    - name: name
      value: "Building Platforms on top of Kubernetes reader!"
  taskRef:
    name: hello-world-task
```

Apliquemos el recurso `TaskRun` a nuestro clúster para crear su primera ejecución de tareas

```shell
kubectl apply -f hello-world/task-run.yaml
taskrun.tekton.dev/hello-world-task-run-1 created
```

Tan pronto se crea `TaskRun`, el motor de pipelines de Tekton se encarga de programar las tareas y crear el Pod de
Kubernetes necesario para ejecutarlas. Si enumeras los pods en el namespace predeterminado, deberías ver algo como esto:

```shell
kubectl get pods
NAME                         READY   STATUS     RESTARTS   AGE
hello-world-task-run-1-pod   0/1     Init:0/1   0          2s
```

También puedes enumerar los `TaskRun` para verificar su estado:

```shell
kubectl get taskrun
NAME                     SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-task-run-1   True        Succeeded   66s         7s
```

Finalmente, dado que estábamos ejecutando una sola tarea, puedes ver los registros de la ejecución del `TaskRun`
observando los registros del pod que se creó:

```shell
kubectl logs -f hello-world-task-run-1-pod 
Defaulted container "step-echo" out of: step-echo, prepare (init)
Hello World: Building Platforms on top of Kubernetes reader!
```

Ahora veamos cómo secuenciar múltiples tareas juntas utilizando un Pipeline de Tekton.

## Comenzando con los Pipelines de Tekton

Ahora podemos usar Pipelines para coordinar múltiples tareas, como la que definimos anteriormente. También podemos
reutilizar las definiciones de tareas creadas por la comunidad de Tekton desde el [Tekton Hub](https://hub.tekton.dev/).

![Tekton Hub](imgs/tekton-hub.png)

Antes de crear el Pipeline, instalaremos la tarea `wget` de Tekton Hub ejecutando:

```shell
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/wget/0.1/wget.yaml
```

Deberías ver:

```
task.tekton.dev/wget created
```

Ahora utilicemos nuestra tarea `Hello World` y la tarea `wget` que acabamos de instalar juntas en un pipeline simple.

Crearemos esta definición de Pipeline simple, que descargará un archivo, leerá su contenido y luego usará la
tarea `Hello World` definida anteriormente.

![Hello World Pipeline](imgs/hello-world-pipeline.png)

Vamos a crear la siguiente definición de pipeline:

```yaml
apiVersion: tekton.dev/v1
kind: Pipeline
metadata:
  name: hello-world-pipeline
  annotations:
    description: |
      Fetch resource from internet, cat content and then say hello
spec:
  results:
    - name: message
      type: string
      value: $(tasks.cat.results.messageFromFile)
  params:
    - name: url
      description: resource that we want to fetch
      type: string
      default: ""
  workspaces:
    - name: files
  tasks:
    - name: wget
      taskRef:
        name: wget
      params:
        - name: url
          value: "$(params.url)"
        - name: diroptions
          value:
            - "-P"
      workspaces:
        - name: wget-workspace
          workspace: files
    - name: cat
      runAfter: [ wget ]
      workspaces:
        - name: wget-workspace
          workspace: files
      taskSpec:
        workspaces:
          - name: wget-workspace
        results:
          - name: messageFromFile
            description: the message obtained from the file
        steps:
          - name: cat
            image: bash:latest
            script: |
              #!/usr/bin/env bash
              cat $(workspaces.wget-workspace.path)/welcome.md | tee /tekton/results/messageFromFile
    - name: hello-world
      runAfter: [ cat ]
      taskRef:
        name: hello-world-task
      params:
        - name: name
          value: "$(tasks.cat.results.messageFromFile)"
```

Al final no es tan fácil recuperar un archivo, leer su contenido y luego usar nuestra tarea
`hello-world` previamente definida para imprimir el contenido del archivo que hemos descargado.

Con los pipelines, tenemos la flexibilidad de agregar nuevas tareas si es necesario para realizar transformaciones o
procesamiento adicional de las entradas y salidas de cada tarea individual.

Para este ejemplo, estamos usando la tarea `wget` que instalamos desde Tekton Hub, y una tarea definida en línea llamada
`cat` que básicamente obtiene el contenido del archivo descargado y lo almacena en un resultado de Tekton que puede ser
referenciado posteriormente en nuestra `hello-world-task`.

Ahora, instala esta definición de pipeline ejecutando:

```shell
kubectl apply -f hello-world/hello-world-pipeline.yaml
```

Luego, podemos crear un nuevo `PipelineRun` cada vez que queramos ejecutar este pipeline:

```yaml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: hello-world-pipeline-run-1
spec:
  workspaces:
    - name: files
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 1M
  params:
    - name: url
      value: "https://raw.githubusercontent.com/salaboy/salaboy/main/welcome.md"
  pipelineRef:
    name: hello-world-pipeline

```

Debido a que nuestras tareas necesitan descargar y almacenar archivos en el sistema de archivos, estamos utilizando
espacios de trabajo de Tekton como abstracciones para proporcionar almacenamiento para nuestros `PipelineRun`. Al igual
que hicimos antes con nuestro `TaskRun`, también podemos proporcionar parámetros para el `PipelineRun`, lo que nos
permite parametrizar cada ejecución para usar diferentes configuraciones, o en este caso, diferentes archivos.

Tanto con `PipelineRuns` como con `TaskRuns`, necesitarás generar un nuevo nombre de recurso para cada ejecución. Si
intentas volver a aplicar el mismo recurso dos veces, el servidor de API de Kubernetes no te permitirá modificar el
recurso existente con el mismo nombre.

Ejecuta este pipeline ejecutando:

```shell
kubectl apply -f hello-world/pipeline-run.yaml
```

Verifica los pods que se crean:

```shell
> kubectl get pods
NAME                                         READY   STATUS        RESTARTS   AGE
affinity-assistant-ca1de9eb35-0              1/1     Terminating   0          19s
hello-world-pipeline-run-1-cat-pod           0/1     Completed     0          11s
hello-world-pipeline-run-1-hello-world-pod   0/1     Completed     0          5s
hello-world-pipeline-run-1-wget-pod          0/1     Completed     0          19s
```

Nota que hay un Pod por cada tarea y un pod llamado `affinity-assistant-ca1de9eb35-0`, que se encarga de asegurar que
los Pods se creen en el nodo correcto (donde se vinculó el volumen).

También verifica los `TaskRuns`:

```shell
> kubectl get taskrun
NAME                                     SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-pipeline-run-1-cat           True        Succeeded   109s        104s
hello-world-pipeline-run-1-hello-world   True        Succeeded   103s        98s
hello-world-pipeline-run-1-wget          True        Succeeded   117s        109s

```

Y, por supuesto, si todas las tareas son exitosas, el `PipelineRun` también lo será:

```shell
kubectl get pipelinerun
NAME                         SUCCEEDED   REASON      STARTTIME   COMPLETIONTIME
hello-world-pipeline-run-1   True        Succeeded   2m13s       114s
```

Asegúrate de revisar las ejecuciones del pipeline y de las tareas en el Dashboard de Tekton si lo has instalado.
![Tekton Dashboard](imgs/tekton-dashboard-hello-world-pipeline.png)

## Tekton para Pipelines de Servicios

Los Pipelines de Servicios en la vida real son mucho más complejos que los ejemplos simples anteriores. Esto se debe
principalmente a que las tareas del pipeline necesitarán tener configuraciones y credenciales especiales para acceder a
sistemas externos.

Un ejemplo de la definición de Pipeline de Servicios se puede encontrar en este directorio en un archivo
llamado [service-pipeline.yaml](service-pipeline.yaml).

![Service Pipeline](imgs/service-pipeline.png)

El Pipeline de Servicios de ejemplo utiliza [`ko`] para construir y publicar la imagen del contenedor para nuestro
servicio. Este pipeline es muy específico para nuestros servicios en Go; si estuviéramos construyendo servicios
utilizando un lenguaje de programación diferente, necesitaríamos usar otras herramientas. El pipeline de servicio de
ejemplo puede ser parametrizado para construir diferentes servicios.

To be able to run this Service Pipeline you need to set up credentials to a Container Registry, this means allowing the
pipelines to push containers to a container registry such as Docker Hub. To authenticate with a container registry from
a Tekton
Task/Pipeline

Para poder ejecutar este Pipeline de Servicios, necesitas configurar credenciales para un Registro de Contenedores, lo
que significa permitir que los pipelines envíen contenedores a un registro de contenedores, como Docker Hub. Para
autenticarte con un registro de contenedores desde una tarea/pipeline de
Tekton [check the official documentation](https://tekton.dev/docs/how-to-guides/kaniko-build-push/#container-registry-authentication).

Para este ejemplo, crearemos un `Secret` de Kubernetes con nuestras credenciales de Docker Hub:

```shell
kubectl create secret docker-registry docker-credentials --docker-server=https://index.docker.io/v1/ --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
```

Luego, instalaremos las tareas `Git Clone` y `ko` de Tekton:

```shell
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/git-clone/0.9/git-clone.yaml
kubectl apply -f https://raw.githubusercontent.com/tektoncd/catalog/main/task/ko/0.1/ko.yaml
```

Vamos a instalar la definición de nuestro Pipeline de Servicios en el clúster:

```shell
kubectl apply -f service-pipeline.yaml
```

Ahora podemos crear nuevas instancias del pipeline para construir y publicar las imágenes de contenedores de nuestros
servicios. El siguiente recurso `PipelineRun` configura nuestro Pipeline de Servicios para construir el Servicio de
Notificaciones.

```yml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: service-pipeline-run-1
  annotations:
    kubernetes.io/ssh-auth: kubernetes.io/dockerconfigjson
spec:
  params:
    - name: target-registry
      value: docker.io/salaboy
    - name: target-service
      value: notifications-service
    - name: target-version
      value: 1.0.0-from-pipeline-run
  workspaces:
    - name: sources
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 100Mi
    - name: docker-credentials
      secret:
        secretName: docker-credentials
  pipelineRef:
    name: service-pipeline
```

Aplica esta definición de `PipelineRun` al clúster para crear una nueva instancia del Pipeline de Servicios:

```shell
kubectl apply -f service-pipeline-run.yaml
```

Observa la sección `spec.params`, que necesitarás modificar para que el pipeline envíe la imagen de contenedor
resultante
a tu propio registro. En otras palabras, reemplaza `docker.io/salaboy` con tu registro + nombre de usuario. El parámetro
`target-service` te permite elegir de qué servicio de la aplicación de conferencia deseas construir (de los servicios
disponibles: `notifications-service`, `agenda-service`, `c4p-service`, `frontend`).

Hay un pipeline separado ([app-helm-chart-pipeline.yaml](app-helm-chart-pipeline.yaml)) que empaqueta y publica el Helm
Chart, que incluye todos los servicios de la aplicación. Cuando el equipo decide la combinación de servicios y la
versión que desean incluir en el Helm Chart, pueden ejecutar otro pipeline para empaquetar y publicar el chart en el
mismo registro de contenedores donde se publican las imágenes de los servicios.

![Helm Chart Application Pipeline](imgs/app-helm-pipeline.png)

Puede instalar `Application Helm Chart Pipeline` ejecutando:

```shell
kubectl apply -f app-helm-chart-pipeline.yaml
```

Luego, puedes crear nuevas instancias creando nuevos recursos `PipelineRun`:

```yml
apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: app-helm-chart-pipeline-run-1
  annotations:
    kubernetes.io/ssh-auth: kubernetes.io/dockerconfigjson
spec:
  params:
    - name: target-registry
      value: docker.io/salaboy
    - name: target-version
      value: v0.9.9
  workspaces:
    - name: sources
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 100Mi
    - name: dockerconfig
      secret:
        secretName: docker-credentials
  pipelineRef:
    name: app-helm-chart-pipeline
```

Aplica esta definición de `PipelineRun` al clúster para crear una nueva instancia del Pipeline del Helm Chart de la
aplicación:

```shell
kubectl apply -f app-helm-chart-pipeline-run.yaml
```

Nota que el pipeline `Application Helm Char` también utiliza las mismas credenciales de `docker-credentials` para
enviar el Helm Chart como una imagen de contenedor OCI. El pipeline acepta el parámetro `target-version`, que se utiliza
para parchear el archivo `Chart.yaml` antes de empaquetar y enviar el Helm Chart al registro de contenedores OCI.

Ten en cuenta que este pipeline no actualiza las versiones de los contenedores referenciados por el chart, lo que
significa que corresponde al usuario adaptar el pipeline para aceptar como parámetros las versiones de cada servicio y
validar que las imágenes de contenedor referenciadas existan en el registro de contenedores mencionado.

**Nota**: Estos pipelines son solo ejemplos para ilustrar el trabajo necesario para configurar Tekton para construir
contenedores
y charts. Por ejemplo, el Pipeline del Helm Chart de la aplicación no cambia la versión del chart ni la versión de las
imágenes de contenedor referenciadas dentro del chart. Si realmente queremos automatizar todo el proceso, podemos
obtener las versiones de las imágenes y la versión del chart a partir de una etiqueta de Git que represente la versión
que deseamos liberar.

## Limpiar

Si deseas eliminar el clúster KinD creado para estos tutoriales, puedes ejecutar:

```shell
kind delete clusters dev
```

