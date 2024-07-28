# bly.li

### Welcome to bly.li, a Short-Url-Service written in GoLang!

With bly.li, you can easily create and share short, memorable links to any webpage. Simply enter the URL you want to shorten, and bly.li will generate a shortened link that you can share on social media, in emails, or anywhere else you want to share a link.

<image style="background:darkgray; padding:1rem; margin:auto; border-radius:1rem" src="./etc/assets/blyli.arch.svg">

| Service | Description                                                                                 |
| ------- | ------------------------------------------------------------------------------------------- |
| Traefik | Reverse-Proxy                                                                               |
| Shortn  | Shortener-Service - Handels new shortn-requests and stores them into redis                  |
| BlowUp  | Redirects to the saved URL - So it blows up the short-url to full and redirects via 302     |
| Redis   | The data-platform for this app. Stores shortens and is used for messaging via redis Pub/Sub |

### Project-Structur

```
└──src
    └──services
    │   ├──shortn
    │   └──blowup
    └──shared
```

### Deployment

Deployment works via docker.
Images are hosted on dockerhub. Feel free to use the docker-compose.demo.yml as a blueprint for your deployment.

| Service                  | URL                                                                                  | Tag    | Description                     |
| ------------------------ | ------------------------------------------------------------------------------------ | ------ | ------------------------------- |
| mensch0mat/bly.li.blowup | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.blowup/) | latest | latest stable version           |
| "                        | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.blowup/) | main   | on push builds from main-branch |
| mensch0mat/bly.li.shortn | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.blowup/) | latest | latest stable version           |
| "                        | [hub.docker.com](https://hub.docker.com/repository/docker/mensch0mat/bly.li.blowup/) | main   | on push builds from main-branch |
