[![golang-pipeline](https://github.com/Menschomat/bly.li/actions/workflows/push.yml/badge.svg)](https://github.com/Menschomat/bly.li/actions/workflows/push.yml) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Menschomat_bly.li&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=Menschomat_bly.li)

# bly.li

### Welcome to bly.li: A scalable Short URL Service in GoLang

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
- **Configurable Logging**: Centralized logging with Loki, configurable log levels and endpoints.
- **OpenID Connect Integration**: Built-in OIDC support for secure authentication.
- **Metrics & Monitoring**: Each service exposes a dedicated Prometheus metrics endpoint for comprehensive monitoring.
- **Graceful Shutdown**: All services support graceful shutdown and cleanup of resources.

<image style="background:darkgray; padding:1rem; margin:auto; border-radius:1rem" src="./etc/assets/blyli.arch.svg">

### Services

| Service | Description                                                                                                        | Main Port | Metrics Port |
| ------- | ------------------------------------------------------------------------------------------------------------------ | --------- | ------------ |
| Traefik | Reverse-Proxy and Load-balancer                                                                                    | :80/:443  | -            |
| Shortn  | Shortener-Service - Handles new shortn-requests, manages URL shortening with distributed counter system            | :8082     | :9082        |
| BlowUp  | URL Resolution Service - Resolves short URLs to their original form and handles redirects with configurable status | :8081     | :9081        |
| Dasher  | Dashboard Service - Provides analytics and management interface for shortened URLs                                 | :8083     | :9083        |
| Perso   | Personalization Service - Handles user preferences and custom URL management with periodic cleanup                 | :8084     | :9084        |
| Front   | Angular-based frontend with modern UI/UX, built with Tailwind CSS                                                  | :80       | -            |
| Redis   | High-performance data store for URL mappings and inter-service communication via Pub/Sub                           | :6379     | -            |
| MongoDB | Persistent storage for URL data and analytics                                                                      | :27017    | -            |

### Environment Variables

#### Shared Configuration
| Variable    | Default                           | Description                           |
| ----------- | --------------------------------- | ------------------------------------- |
| LOKI_URL    | http://loki:3100/loki/api/v1/push | Loki logging endpoint                 |
| LOKI_TENANT | single                            | Loki tenant ID                        |
| LOG_LEVEL   | info                              | Logging level (debug/info/warn/error) |
| INSTANCE_ID | [hostname]                        | Instance identifier for logging       |

#### MongoDB Configuration
| Variable         | Default                 | Description            |
| ---------------- | ----------------------- | ---------------------- |
| MONGO_DATABASE   | short_url_db            | MongoDB database name  |
| MONGN_SERVER_URL | mongodb://mongodb:27017 | MongoDB connection URL |

#### OIDC Configuration
| Variable       | Default          | Description                 |
| -------------- | ---------------- | --------------------------- |
| OIDC_CLIENT_ID | 12345            | OpenID Connect client ID    |
| OIDC_URL       | http://127.0.0.1 | OpenID Connect provider URL |

#### Shortn Service Configuration
| Variable               | Default            | Description                            |
| ---------------------- | ------------------ | -------------------------------------- |
| SERVER_PORT            | :8082              | HTTP server port                       |
| METRICS_PORT           | :9082              | Prometheus metrics endpoint port       |
| CORS_ALLOWED_ORIGINS   | https://*,http://* | Allowed CORS origins (comma-separated) |
| CORS_MAX_AGE           | 300                | CORS preflight max age                 |
| ZOOKEEPER_HOST         | zookeeper:2181     | Zookeeper connection host              |
| ZOOKEEPER_COUNTER_PATH | /shortn-ranges     | Zookeeper counter root path            |

#### BlowUp Service Configuration
| Variable      | Default | Description                      |
| ------------- | ------- | -------------------------------- |
| SERVER_PORT   | :8081   | HTTP server port                 |
| METRICS_PORT  | :9081   | Prometheus metrics endpoint port |
| REDIRECT_CODE | 302     | HTTP redirect status code        |

#### Dasher Service Configuration
| Variable     | Default | Description                      |
| ------------ | ------- | -------------------------------- |
| SERVER_PORT  | :8083   | HTTP server port                 |
| METRICS_PORT | :9083   | Prometheus metrics endpoint port |

#### Perso Service Configuration
| Variable         | Default | Description                      |
| ---------------- | ------- | -------------------------------- |
| SERVER_PORT      | :8084   | HTTP server port                 |
| METRICS_PORT     | :9084   | Prometheus metrics endpoint port |
| CLEANUP_INTERVAL | 24h     | Interval for cleanup operations  |

### Project Structure

```
└──src
    └──services
    |   ├─-front    # Angular frontend application
    │   ├──shortn   # URL shortening service with metrics
    │   ├──blowup   # URL resolution service with metrics
    │   ├──dasher   # Analytics dashboard service with metrics
    │   └──perso    # User personalization service with metrics
    └──shared       # Shared Go packages and utilities
        ├──config   # Configuration management
        ├──model    # Data models
        ├──mongo    # MongoDB client and operations
        ├──redis    # Redis client and operations
        └──utils    # Shared utilities
```

### Deployment

Deployment works via Docker. Images are hosted on GitHub Container Registry (ghcr.io). Use the docker-compose.yml as a blueprint for your deployment.

| Service                       | Description                  |
| ----------------------------- | ---------------------------- |
| ghcr.io/[owner]/bly.li/front  | Frontend Angular application |
| ghcr.io/[owner]/bly.li/blowup | URL resolution service       |
| ghcr.io/[owner]/bly.li/shortn | URL shortening service       |
| ghcr.io/[owner]/bly.li/dasher | Analytics dashboard service  |
| ghcr.io/[owner]/bly.li/perso  | User personalization service |

Each service is available with the following tags:
- `latest`: Latest stable release version (from tags)
- `main`: Latest build from the main branch

The images are built for linux/amd64 platform and include automated build caching via GitHub Actions.

### License

This project is licensed under the GNU Affero General Public License v3.0. You can view the full license text in the LICENSE file. This license ensures that the code remains free and open, promoting collaboration and sharing within the community.