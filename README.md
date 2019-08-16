# Image spider

A image spider implemented using Golang

## usage

put a `config.json` file in current folder

default config is as below

### Install Elasticsearch docker

```sh
docker pull elasticsearch:5.6.16
docker run -d -p 9200:9200  elasticsearch:5.6.16
```

```json
{
  "image": {
    "path": "d:\\t66yImage",
    "workerCount": 20,
    "isUniqFolder": true
  },
  "engine": {
    "workerCount": 2,
    "elasticUrl": "http://172.16.3.116:9200"
  },
  "init": {
    "topicPerPage": 100,
    "startPageNum": 20,
    "endPageNum": 100,
    "realms": [
      "http://t66y.com/thread0806.php?fid=8"
    ]
  },
  "nameLenLimit": 60,
  "net": {
    "timeOut": 30,
    "proxyUrl": "socks5://localhost:1080"
  }
}
```

### Downloaod image and save elastic search

```sh
go get github.com/Tomilla/imagespider
cd $GOPATH/src/github.com/Tomilla/imagespider
go run github.com/Tomilla/imagespider
```

### Web search page

```sh
cd $GOPATH/src/github.com/Tomilla/imagespider
go run github.com/Tomilla/imagespider/frontend
```

## arch

![arch](./mis/arch.png)
