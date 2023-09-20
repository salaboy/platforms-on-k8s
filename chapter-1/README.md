# Chapter 1 :: (The rise of) Platforms on Top of Kubernetes

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README-zh.md) | [Português (Portuguese)](README-pt.md) | [Español](README-es.md)
> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

## Pre-requisites for the tutorials

You'll need the tools below to follow the step-by-step tutorials linked in the book:
- [Docker](https://docs.docker.com/engine/install/), v24.0.2
- [kubectl](https://kubernetes.io/docs/tasks/tools/), Client v1.27.3
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/), v0.20.0
- [Helm](https://helm.sh/docs/intro/install/), v3.12.3

These are technologies and versions used when testing the tutorials.

> [!Warning]
> If you want to use other technologies, like [Podman](https://podman.io/) instead of Docker, it should be possible as there is nothing specific to Docker.

## Conference Application Scenario

The application that we will modify and use throughout the book's chapters represents a simple "walking skeleton", meaning that it is complex enough to allow us to test assumptions, tools, and frameworks. Still, it is not the final product that our customers will use.

The "Conference Application" walking skeleton implements a straightforward use case, allowing potential _speakers_ to submit proposals that the conference _organizers_ will evaluate. See below the app's home page:

![home](imgs/homepage.png)

Check below how the application is commonly used:
1. **C4P:** Potential _speakers_ can submit a new proposal by going to the application's **Call for Proposals** (C4P) section.
   ![proposals](imgs/proposals.png)
2. **Review & Approval**: Once a proposal is submitted, the conference _organizers_ can review (approve or reject) them by using the **Backoffice** section of the application.
   ![backoffice](imgs/backoffice.png)
3. **Announcement**: If accepted by the _organizers_, the proposal is automatically published on the conference **Agenda** page.
   ![agenda](imgs/agenda.png)
4. **Speaker's Notification**: In the **Backoffice**, a _speaker_ can check the **Notifications** tab. There, potential _speakers_ can find all the notifications (emails) sent them. A speaker will see both approval and rejection emails in this tab.
   ![notifications](imgs/notifications-backoffice.png)

### An event-driven application

**Every action in the application results in new events being emitted.** For instance, events are emitted:
- when a new proposal is submitted;
- when the proposal is accepted or rejected;
- when notifications are sent.

These events are sent and then captured by a frontend application. Luckily, you, the reader, can see these details in the app by accessing **Events** tab in the **Backoffice** section.

![events](imgs/events-backoffice.png)

## Sum up and Contribute

Do you want to improve this tutorial? Create an [issue](https://github.com/salaboy/platforms-on-k8s/issues/new), message me on [Twitter](https://twitter.com/salaboy), or send a [Pull Request](https://github.com/salaboy/platforms-on-k8s/compare).
