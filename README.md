Blossom-Reposter
================
Blossom-Reposter - Monitors & reposts activities from popular social networks to telegram and /fag/
<h1 align="center">
  <img src="https://raw.githubusercontent.com/wMw9/makaba-reposter/master/assets/img/post_sample.png" alt="Blossom Reposter"></a>
  <br>
</h1>
  <br>
  <img src="https://raw.githubusercontent.com/wMw9/makaba-reposter/master/assets/img/reposter_scheme.svg" alt="Reposter Scheme"></a>
  <h3 align="center">
  Reposter Scheme
  </h3>
#### "Best friend for telegram channel owners"

Managing your telegram channels and stalking your waifus has never been as fun nor easier than this. It comes with features like:

* Save instagram posts and stories leaving no trace
* Save posts from VK pages, publics, status
* Announce when streamer goes online using Twitch WebHooks
* Repost these activities to /fag/ or telegram channel for archives purposes
* more to come

## Installation
This is a personal project, so wasn't meant for public use. If you have questions how to configurate or install you should find it yourself
There are 3 ways of starting this service up. Classic way, docker-compose or kubernetes.
For docker-compose you just 
```docker-compose up -d ```
For kubernetes you should check [/k8s](k8s/) folder for kubernetes manifests *.yaml files. Ingress proxy and SSL certification configurations are up to you. I personally use nginx-ingress and cert-manager.

## License
Blossom-Reposter is MIT License.
