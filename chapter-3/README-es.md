# Capitulo 3: Pipelines de servicios: Construcción de aplicaciones nativas en la nube

—
_🌍 Disponible
en_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [日本語 (Japanese)](README-ja.md) | [Español](README-es.md)
> **Nota:** Presentado por la fantástica comunidad
> de [ 🌟 contribuidores](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) cloud-native!

---

Estos breves tutoriales cubren tanto Tekton como Dagger para pipelines de servicios. Con Tekton, extendemos las API de
Kubernetes para definir nuestros Pipelines y tareas. Con Dagger, definimos programáticamente los Pipelines que pueden
ejecutarse de forma remota en Kubernetes o localmente en nuestras laptops de desarrollo. Finalmente, se proporciona un
enlace a un conjunto de GitHub Actions para poder comparar entre estos diferentes enfoques.

- [Tutorial de Tekton para pipelines de servicios](tekton/README-es.md)
- [Tutorial de Dagger para pipelines de servicios](dagger/README-es.md)
- [GitHub Actions](github-actions/README-es.md)

## Limpiar

Si desea deshacerse del KinD Cluster creado para estos tutoriales, puede ejecutar:

```shell
kind delete clusters dev
```

## Próximos pasos

Recomiendo encarecidamente seguir los tutoriales listados para Tekton y Dagger en tus entornos locales. Si tienes
experiencia en desarrollo, puedes extender los pipelines de Dagger con tus propios pasos personalizados.

Si no eres un desarrollador de Go, ¿te atreverías a construir un pipeline para tu stack tecnológico usando Tekton y
Dagger?

Si tienes una cuenta en un registro de contenedores, como una cuenta de Docker Hub, puedes intentar configurar las
credenciales para que los pipelines puedan enviar las imágenes de contenedores al registro. Luego puedes usar el Helm
Chart `values.yaml` para consumir las imágenes desde tu cuenta en lugar de las oficiales alojadas
en `docker.io/salaboy`.

Finalmente, puedes hacer un fork del repositorio `salaboy/platforms-on-k8s` en tu propio usuario (organización) de
GitHub para probar con la ejecución de los GitHub Actions ubicadas en este
directorio. [../../.github/workflows/](../../.github/workflows/).

## Resumir y contribuir

En estos tutoriales, experimentamos con dos enfoques completamente diferentes para pipelines de servicios. Comenzamos
con [Tekton](https://tekton.dev), un motor de pipelines sin opiniones que fue diseñado para ser nativo de Kubernetes,
aprovechando el
poder declarativo de las API de Kubernetes y los bucles de reconciliación de Kubernetes. Luego
probamos [Dagger](https://dagger.io), un motor
de pipelines diseñado para orquestar contenedores y que se puede configurar utilizando tu stack tecnológico favorito a
través de sus SDKs.

Una cosa es segura, no importa qué motor de pipelines elija tu organización, los equipos de desarrollo se beneficiarán
enormemente al poder consumir pipelines de servicios sin la necesidad de definir cada paso, credenciales y detalles de
configuración. Si el equipo de plataforma puede crear pasos/tareas compartidos o incluso pipelines predeterminados para
tus servicios, los desarrolladores pueden centrarse en escribir código de aplicación.

¿Quieres mejorar este tutorial? Crea un issue, envíame un mensaje en [Twitter](https://twitter.com/salaboy)  o envía un
Pull Request.