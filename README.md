# bly.li

### Welcome to bly.li, a Short-Url-Service written in GoLang!

With bly.li, you can easily create and share short, memorable links to any webpage. Simply enter the URL you want to shorten, and bly.li will generate a shortened link that you can share on social media, in emails, or anywhere else you want to share a link.

| Service | Description                                                                                 |
| ------- | ------------------------------------------------------------------------------------------- |
| Traefik | Reverseproxy                                                                                |
| Shortn  | Shortener-Service - Handels new shortn-requests and stores them into redis                  |
| BlowUp  | Redirects to the saved URL - So it blows up the short-url to full and ridercts via 302      |
| Redis   | The data-platform for this app. Stores shortens and is used for messaging via redis Pub/Sub |
