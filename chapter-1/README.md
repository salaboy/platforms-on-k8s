# Chapter 1 :: (The rise of) Platforms on Top of Kubernetes

---
_🌍 Available in_: [English](README.md) | [中文 (Chinese)](README.zh-cn.md)
> **Note:** Brought to you by the fantastic cloud-native community's [ 🌟 contributors](https://github.com/salaboy/platforms-on-k8s/graphs/contributors)!

---

## Conference Application Scenario

The application that we will modify and use throughout the book chapters represents a simple "walking skeleton", meaning that it is complex enough to allow us to test assumptions, tools, and frameworks. Still, it is not the final product that our customers will use. 

The "Conference Application" walking skeleton implements a straightforward use case, allowing potential speakers to submit proposals that the conference organizers will evaluate. 

![home](imgs/homepage.png)


The flow is simple. Potential speakers can submit a new proposal by going to the application's **Call for Proposals** section.

![proposals](imgs/proposals.png)

Once submitted, the conference organizers can review (approve or reject) submitted proposals in the **Backoffice** section of the application.

![backoffice](imgs/backoffice.png)

If accepted, the proposal is automatically published on the conference **Agenda** page.

![agenda](imgs/agenda.png)

In the **Backoffice**, you can check the **Notifications** tab that shows all the notifications (emails) sent to the potential speakers. You will see both approval and rejection emails in this tab. 

![notifications](imgs/notifications-backoffice.png)

Every action in the application emits events. Hence, when a new proposal is submitted, when the proposal is accepted or rejected, and when notifications are sent, events are sent and captured by the application frontend. You can check these events in the **Events** tab in the **Backoffice** section.

![events](imgs/events-backoffice.png)


## Pre-requisites for the other chapters

The following tools are required for the step-by-step tutorials linked in the book. 


- [Docker](https://docs.docker.com/engine/install/)
  - Note: You can try to use [Podman](https://podman.io/) as well, as there is nothing specific to Docker, but all the tutorials had been tested with Docker.
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Helm](https://helm.sh/docs/intro/install/) 


## Sum up and Contribute

Do you want to improve this tutorial? Create an issue, message me on [Twitter](https://twitter.com/salaboy), or send a Pull Request.

