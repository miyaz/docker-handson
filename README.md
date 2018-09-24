# docker hands-on

## コンテナ動作確認

### アプリ動作確認
ハンズオンで使うアプリの動作を確認します（スキップOK）

事前にGolangがインストールされている必要あります  
macOSお使いの方は、 `brew install go`　で入ります  
```
cd dockergo
go run main.go
# 別のターミナルを開いてcurlコマンド実行
curl 127.0.0.1:8080
```

### コンテナイメージのビルド
Dockerイメージを作成します
```
# Dockerイメージを作成し、example/echo:latest という名前でタグ付け
docker image build -t example/echo:latest .
```

### イメージ確認
```
docker image ls
# example/echoのイメージが表示されます
```

### docker起動
イメージを使ってdockerコンテナを起動します（スキップOK）
```
# ホスト側9000ポート　から コンテナの8080ポートにポートフォワード
# コンテナにgoechoという名前をつけています（名前指定で起動、停止、削除などができます）
docker container run -d -p 9000:8080 --name goecho example/echo:latest
docker container ps
# 別のターミナルを開いてcurlコマンド実行
curl 127.0.0.1:9000
docker container stop goecho
docker container rm goecho
```

## DockerSwarm環境作成

### Docker in Dockerホスト起動
１台のホスト（物理マシン）上に複数ホストを擬似的に作るためにDocker in Dockerの環境を作成します  
manager(クラスタ管理)/registry(Dockerレジストリ)/worker01~03(ワーカー)
```
cd swarmgo
# dindを使って5台のDockerホストを起動します
docker-compose up -d
# 起動確認
docker container ls
```

### Swarmクラスタ作成
Swarmで複数Dockerホストを跨ったクラスタを作成し管理します
```
docker container exec -it manager docker swarm init
# 表示されたJOINトークンを使い次のコマンドを実行
docker container exec -it worker01 \
docker swarm join --token SWMTKN-〜〜〜 172.20.0.3:2377  # ← swarm init実行時に表示された文字列を指定
# 'This node joined a swarm as a worker.' が表示されればOK
# worker02, worker03 でも同様に実行する
# Swarmクラスタのノード状態確認
docker container exec -it manager docker node ls
```

### DockerレジストリへイメージPush
registryというDockerレジストリ用のコンテナにイメージをpushする
```
# ホスト側からregistryコンテナにpushできるようにタグをつける
docker image tag example/echo:latest localhost:5000/example/echo:latest
# タグ付けイメージ確認(localhost:5000/example/echo存在確認)
docker image ls
# registryコンテナにイメージをpush
docker image push localhost:5000/example/echo:latest
```

### Pull動作確認
registryにpushしたイメージをworker01からpullできることを確認
```
docker container exec -it worker01 docker image pull registry:5000/example/echo:latest
```

## DockerSwarmクラスタ上でのアプリ動作確認

### Service作成
アプリケーションイメージ単位であるServiceを作成してみます（スキップOK）

```
# managerコンテナからserviceを作成(動かすアプリは前述のdockergoをechoという名前で指定）
docker container exec -it manager \
docker service create --replicas 1 --publish 8000:8080 --name echo registry:5000/example/echo:latest
# service確認
docker container exec -it manager docker service ls
# コンテナ数を増やしてみる
docker container exec -it manager docker service scale echo=6
# コンテナの起動状態を確認
docker container exec -it manager docker service ps echo | grep Running
# service[echo]を削除
docker container exec -it manager docker service rm echo
```

### アプリ用Stack作成
複数のSerivceをグルーピングした単位であるStackを作成します  
今回クライアントからアプリケーションにアクセスするためにStack内で同じoverlayネットワークを  
使うため必要になります
```
# handsonというoverlayネットワークを作成します
docker container exec -it manager docker network create --driver=overlay --attachable handson
# アプリ用Stack[echo]を作成します
docker container exec -it manager docker stack deploy -c /stack/webapi.yml echo
# echoスタックのService一覧を表示します
docker container exec -it manager docker stack services echo
```

### アプリアクセス用Stack作成
echo_nginxのServiceへの橋渡し(Proxy)するためのStackをmanagerコンテナに配置します
```
docker container exec -it manager docker stack deploy -c /stack/ingress.yml ingress
docker container exec -it manager docker stack services ingress
# -> ブラウザで http://127.0.0.1:8000/ にアクセス
```

### コンテナ配置確認
コンテナ配置確認用のStackを作る（スキップOK）
```
docker container exec -it manager docker stack deploy -c /stack/visualizer.yml visualizer
docker container exec -it manager docker stack services visualizer
# -> ブラウザで http://127.0.0.1:9000/ にアクセス
# アプリのコンテナ数を変えてみる
docker container exec -it manager docker service scale echo_api=5
# 再度visualizerやブラウザで配置変更を確認
curl 127.0.0.1:8000
```

## ハンズオン後の削除
```
docker container stop $(docker container ls -aq)
docker container rm $(docker container ls -aq)
docker image rm -f $(docker image ls -aq)
```

