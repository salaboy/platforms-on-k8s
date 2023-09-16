# Capítulo 1 :: (A ascensão das) Plataformas Baseadas em Kubernetes

Graças à fantástica comunidade cloud-native, você tem acesso a estes tutoriais nas seguintes línguas:
- [中国人 (Chinese) ](README.zh-cn.md)
- [Inglês (English)](README.md)
- [Português (Portuguese) ](README-pt.md) 

## Cenário: Aplicação de Conferência

A aplicação que modificaremos e usaremos ao longo dos capítulos do livro representa um "esqueleto em funcionamento", o que significa que é complexa o suficiente para nos permitir testar suposições, ferramentas e frameworks. No entanto, não é o produto final que nossos clientes usarão.

O "esqueleto em funcionamento" da Aplicação de Conferência implementa um caso de uso bem simples, permitindo que potenciais palestrantes submetam propostas que os organizadores da conferência avaliarão.

![home](imgs/homepage.png)

O fluxo é simples. Palestrantes em potencial podem enviar uma nova proposta indo à area de **Call for Papers** da aplicação.

![proposals](imgs/proposals.png)

Uma vez submetida, os organizadores da conferência podem revisar (aprovar ou rejeitar) as propostas submetidas na área de **Backoffice** da aplicação.

![backoffice](imgs/backoffice.png)

Se aceita, a proposta é automaticamente publicada na página de **Agenda** da conferência.

![agenda](imgs/agenda.png)

No **Backoffice**, você pode verificar a aba **Notificações** que mostra todas as notificações (e-mails) enviadas aos potenciais palestrantes. Você verá e-mails de aprovação e rejeição nesta aba.

![notifications](imgs/notifications-backoffice.png)

Cada ação na aplicação emite eventos. Portanto, quando uma nova proposta é submetida, quando a proposta é aceita ou rejeitada, e quando notificações são enviadas, eventos são enviados e capturados pela interface do usuário da aplicação. Você pode verificar esses eventos na aba **Eventos** na seção **Backoffice**.

![events](imgs/events-backoffice.png)

## Pré-requisitos para os outros capítulos

As seguintes ferramentas são necessárias para os tutoriais passo a passo vinculados no livro.

- [Docker](https://docs.docker.com/engine/install/)
  > Uma vez que não há nada nos tutoriais que seja exclusivo da tecnologia Docker, você pode considerar também o [Podman](https://podman.io/). No entanto, tenha em mente que não foram efetuados testes destes projetos com a utilização desta tecnologia.
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [KinD](https://kind.sigs.k8s.io/docs/user/quick-start/)
- [Helm](https://helm.sh/docs/intro/install/)

## Resumo e Contribuição

Quer melhorar este tutorial? Abre uma issue, envie uma mensagem via [Twitter](https://twitter.com/salaboy) ou submeta um Pull Request.
