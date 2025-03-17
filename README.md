### ルートコンテナの立ち上げ
```sh
$ make up
$ make down #down
```

### コンテナに入る
この作業を`broker-1`, `2`, `3`の各コンテナでも実行する。

```sh
$ docker exec -it broker-1 sh
# $ docker exec -it broker-2 sh
# $ docker exec -it broker-3 sh

# Zookeeperの立ち上げ
/ \# docker compose -f ./src/compose.zookeeper.yml up -d
# Kafkaの立ち上げ
/ \# docker compose -f ./src/compose.kafka.yml up -d
```

### UI 起動用のコンテナに入る
```sh
$ docker exec -it kafka-ui sh
/ \# docker compose -f ./src/compose.ui.yml up -d
```
`localhost:8888`アクセスすると UI が表示される。

