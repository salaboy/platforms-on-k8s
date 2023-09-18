# Chapter 1 :: (The rise of) Platforms on Top of Kubernetes

Thanks to the fantastic cloud-native community, you can read these tutorials in the following languages:

- [Chinese `zh-cn`](README.zh-cn.md) 🇨🇳


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