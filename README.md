<h1 align="center">
  <img src="https://raw.githubusercontent.com/wMw9/makaba-reposter/master/assets/img/post_sample.png" alt="Blossom Reposter">
  <img src="https://raw.githubusercontent.com/wMw9/makaba-reposter/master/assets/img/post_sample2.png" alt="Blossom Reposter">
</h1>
<p align="center">Monitors & reposts activities from popular social networks to telegram and /fag/</p>
<h1 align="center">Blossom-Reposter</h1>
Best friend for telegram channel owners.
Managing your telegram channels and stalking your waifus has never been as fun nor easier than this. It comes with features like:

* Save instagram posts and stories leaving no trace
* Save posts from VK pages, publics, status
* Announce when streamer goes online using Twitch WebHooks
* Repost these activities to /fag/ or telegram channel for archive purposes
* more to come

## Scheme

<img src="https://raw.githubusercontent.com/wMw9/makaba-reposter/master/assets/img/reposter_scheme.svg" alt="Reposter Scheme">

## Installation

This is a personal project wasn't meant for public use. If you have questions how to configurate or install you should find it yourself
There are 3 ways of starting this service up. Classic way, docker-compose or Kubernetes.
For docker-compose you just 
```
docker-compose up -d
```
For Kubernetes you should check [/k8s](k8s/) folder for manifests *.yaml files. Ingress proxy and SSL cert configurations are up to you. I personally use nginx-ingress and cert-manager.

## Used

* [PostgreSQL](https://hub.docker.com/_/postgres) - object-relational database system provides reliability and data integrity
* [RabbitMQ](https://hub.docker.com/_/rabbitmq) - an open source multi-protocol messaging broker
* [AMQP](https://github.com/streadway/amqp) -  Go client for AMQP 0.9.1
* [GORM](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang, aims to be developer friendly
* [GoRequest](https://github.com/parnurzeal/gorequest) - Simplified HTTP client
* [gjson](https://github.com/tidwall/gjson) - get json values quickly

## Contact

Ivan Smyshlyaev [@wmw](https://instagram.com/wmw)

## License

Blossom-Reposter is MIT License.
