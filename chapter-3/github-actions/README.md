# GitHub Actions in Action 

Using GitHub actions we can automate our Service Pipelines without the need of running any infrastructure. All the pipelines / workflows will run on GitHub.com infrastructure, which at scale (when they start charging) this  becomes very expensive.

You can find the Service Pipelines for the Conference Application defined here: 
- [Agenda Service GH Service Pipeline](../../.github/workflows/agenda-service-service-pipeline.yaml)
- [C4P Service GH Service Pipeline](../../.github/workflows/c4p-service-service-pipeline.yaml)
- [Notifications Service GH Service Pipeline](../../.github/workflows/notifications-service-service-pipeline.yaml)

These Pipelines use a filter to be only triggered if the source code for the service changes (#https://github.com/dorny/paths-filter). To build and publish the containers for each service, these pipelines use `ko-build` (https://github.com/ko-build/setup-ko) to generate multi platform containers for our Go applications. 
To publish to Docker Hub the `docker/login-action@v2` action is used and required the configuration of two environment Secrets (`secrets.DOCKERHUB_USERNAME`, `secrets.DOCKERHUB_TOKEN`) which enable the Pipeline runs to push containers to my Docker Hub account.