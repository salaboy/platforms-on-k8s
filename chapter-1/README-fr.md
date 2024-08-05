# Chapitre 1 :: (L'Ascension des) Plateformes sur Kubernetes

---
_🌍 Disponible en_: [Anglais](README.md) | [中文 (Chinois)](README-zh.md) | [Português (Portugais)](README-pt.md) | [Español](README-es.md) | [日本語 (Japonais)](README-ja.md) | [Français](README-fr.md)
> **Remarque:** Proposé par les [ 🌟 contributeurs](https://github.com/salaboy/platforms-on-k8s/graphs/contributors) de la fantastique communauté cloud-native!

---

## Prérequis pour les tutoriels

Vous aurez besoin des outils ci-dessous pour suivre les tutoriels pas à pas mentionnés dans le livre:
- [Docker](https://docs.docker.com/engine/install/), v24.0.2
- [kubectl](https://kubernetes.io/docs/tasks/tools/), Client v1.27.3
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/), v0.20.0
- [Helm](https://helm.sh/docs/intro/install/), v3.12.3

Ce sont les technologies et versions utilisées lors des tests des tutoriels.

> [!Warning]
> Si vous souhaitez utiliser d'autres technologies, comme [Podman](https://podman.io/) au lieu de Docker, c'est possible en activant l'exécution de conteneurs root avec cette commande:
```shell
podman machine set --rootful
```

## Scénario de l'Application de Conférence

L'application que nous allons modifier et utiliser tout au long des chapitres du livre représente un simple "squelette de base", c'est-à-dire qu'elle est suffisamment complexe pour nous permettre de tester des hypothèses, des outils et des frameworks, mais ce n'est pas le produit final que nos clients utiliseront.

Le "squelette de base" de l'Application de Conférence implémente un cas d'utilisation simple, permettant aux _conférenciers_ potentiels de soumettre des propositions que les _organisateurs_ de la conférence évalueront. Voir ci-dessous la page d'accueil de l'application:

![home](imgs/homepage.png)

Voici comment l'application est généralement utilisée:
1. **C4P:** Les _conférenciers_ potentiels peuvent soumettre une nouvelle proposition en se rendant dans la section **Call for Proposals** (C4P) de l'application.
   ![proposals](imgs/proposals.png)
2. **Révision & Approbation:** Une fois une proposition soumise, les _organisateurs_ de la conférence peuvent la réviser (approuver ou rejeter) en utilisant la section **Backoffice** de l'application.
   ![backoffice](imgs/backoffice.png)
3. **Annonce:** Si elle est acceptée par les _organisateurs_, la proposition est automatiquement publiée sur la page **Agenda** de la conférence.
   ![agenda](imgs/agenda.png)
4. **Notification du Conférencier:** Dans le **Backoffice**, un _conférencier_ peut vérifier l'onglet **Notifications**. Là, les _conférenciers_ potentiels peuvent trouver toutes les notifications (emails) qui leur ont été envoyées. Un conférencier verra à la fois les emails d'approbation et de rejet dans cet onglet.
   ![notifications](imgs/notifications-backoffice.png)

### Une application événementielle

**Chaque action dans l'application génère de nouveaux événements.** Par exemple, des événements sont émis:
- lorsqu'une nouvelle proposition est soumise;
- lorsque la proposition est acceptée ou rejetée;
- lorsque des notifications sont envoyées.

Ces événements sont envoyés puis capturés par une application frontend. Heureusement, vous, le lecteur, pouvez voir ces détails dans l'application en accédant à l'onglet **Events** dans la section **Backoffice**.

![events](imgs/events-backoffice.png)

## Résumé et Contribuez

Vous voulez améliorer ce tutoriel? Créez une [issue](https://github.com/salaboy/platforms-on-k8s/issues/new), envoyez-moi un message sur [Twitter](https://twitter.com/salaboy), ou soumettez une [Pull Request](https://github.com/salaboy/platforms-on-k8s/compare).

