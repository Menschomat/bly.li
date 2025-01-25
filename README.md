[![golang-pipeline](https://github.com/Menschomat/bly.li/actions/workflows/push.yml/badge.svg)](https://github.com/Menschomat/bly.li/actions/workflows/push.yml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Menschomat_bly.li&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Menschomat_bly.li)

# bly.li

### Welcome to bly.li: A calable Short URL Service in GoLang

Easily create and share memorable links to any webpage with bly.li! Simply input the URL you want to shorten, and our service will generate a concise link that you can share on social media, via email, or anywhere else.

### Tech-stack

<div>
<image style=" padding:1rem; margin:auto; border-radius:1rem; height: 6rem;" src="./etc/assets/logos/go_logo.png">
<image style=" padding:1rem; margin:auto; border-radius:1rem; height: 4rem;  padding-bottom:1.5rem;" src="./etc/assets/logos/redis_mark.svg">
<image style=" padding:1rem; margin:auto; border-radius:1rem; height: 5rem" src="./etc/assets/logos/mongo_db.svg">
<image style=" padding:1rem; margin:auto; border-radius:1rem; height: 5rem" src="./etc/assets/logos/traefik.svg">
<image style=" padding:1rem; margin:auto; border-radius:1rem; width: 5rem" src="./etc/assets/logos/angular_gradient.png">
<image style=" padding:1rem; margin:auto; border-radius:1rem; width: 5rem; padding-bottom:1.7rem; " src="./etc/assets/logos/tailwindcss.svg">
<image style=" padding:1rem; margin:auto; border-radius:1rem; width: 5rem" src="./etc/assets/logos/bun_logo.svg">
<image style=" padding:1rem; margin:auto; border-radius:1rem; width: 5rem" src="./etc/assets/logos/docker-mark.svg">
</div>

### Key Features:

- **Scalability**: Our architecture is designed for horizontal scaling, allowing services to be easily expanded to meet growing demands.
- **Statelessness**: All services are stateless, ensuring seamless performance even in the event of node failures. This also makes it easy to maintain and update individual components without affecting the overall system.
- **Flexible Cache Management**: Redis can be used as a shared cache cluster or per-node read cache, providing flexibility in scaling depending on your specific needs.
- **MongoDB Scalability**: MongoDB can be scaled using a shared cluster or via a cloud-based solution like MongoDB Atlas.

<image style="background:darkgray; padding:1rem; margin:auto; border-radius:1rem" src="./etc/assets/blyli.arch.svg">

| Service | Description                                                                                 |
| ------- | ------------------------------------------------------------------------------------------- |
| Traefik | Reverse-Proxy and Load-balancer                                                             |
| Shortn  | Shortener-Service - Handels new shortn-requests and stores them into redis                  |
| BlowUp  | Redirects to the saved URL - So it blows up the short-url to full and redirects via 302     |
| Front   | Serves the Angular-Frontend for the stack                                                   |
| Redis   | The data-platform for this app. Caches shortens and is used for messaging via redis Pub/Sub |
| MongoDB | Used for Persistency                                                                        |

### Project-Structur

```
└──src
    └──services
    |   ├─-front
    │   ├──shortn
    │   └──blowup
    └──shared
```

#### front - The Front-End

The frontend is built using Angular. It serves as a single-page-application.
For an easy styling the project uses tailwindcss.
To manage dependencies and build tasks, we're using Bun as an alternative to npm.
Bun performs great and speeds up package installations and building tasks.

### Deployment

Deployment works via docker.
Images are hosted on dockerhub. Feel free to use the docker-compose.demo.yml as a blueprint for your deployment.

| Service                  | URL                                                                                  | Tag    | Description                     |
| ------------------------ | ------------------------------------------------------------------------------------ | ------ | ------------------------------- |
| mensch0mat/bly.li.blowup | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.blowup/) | latest | latest stable version           |
| "                        | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.blowup/) | main   | on push builds from main-branch |
| mensch0mat/bly.li.shortn | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.shortn/) | latest | latest stable version           |
| "                        | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.shortn/) | main   | on push builds from main-branch |
