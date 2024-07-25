# GitHub Actions en Acción

Utilizando GitHub Actions, podemos automatizar nuestros Pipelines de Servicio sin necesidad de ejecutar ninguna
infraestructura. Todos los pipelines/workflows se ejecutarán en la infraestructura de GitHub.com, lo que, a gran
escala (cuando comiencen a cobrar), puede volverse muy costoso.

Puedes encontrar los Pipelines de Servicio para la Aplicación de Conferencia definidos aquí:

- [Pipeline de Servicio para el Servicio de Agenda en GH](../../.github/workflows/agenda-service-service-pipeline.yaml)
- [Pipeline de Servicio para el Servicio C4P en GH](../../.github/workflows/c4p-service-service-pipeline.yaml)
- [Pipeline de Servicio para el Servicio de Notificaciones en GH](../../.github/workflows/notifications-service-service-pipeline.yaml)

Estos pipelines utilizan un filtro para activarse únicamente si el código fuente del servicio
cambia [paths-filter](https://github.com/dorny/paths-filter).

Para construir y publicar los contenedores de cada servicio, estos pipelines
utilizan `ko-build` [setup-ko](https://github.com/ko-build/setup-ko) para generar
contenedores multiplataforma para nuestras aplicaciones Go.

Para publicar en Docker Hub, se utiliza la acción `docker/login-action@v2`, la cual requiere la configuración de dos
`Secrets` (`secrets.DOCKERHUB_USERNAME`, `secrets.DOCKERHUB_TOKEN`), que permiten a los pipelines realizar el push
de los contenedores a mi cuenta de Docker Hub.