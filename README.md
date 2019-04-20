# Image spider

A image spider implemented using Golang

## usage

put a `config.json` file in current folder

default config is as below

```json
{
  "image": {
    "path": "d:\\t66yImage",
    "workerCount": 20,
    "isUniqFolder": true
  },
  "startPages": [
    "http://t66y.com/thread0806.php?fid=8",
    "http://t66y.com/thread0806.php?fid=16",
    "http://t66y.com/thread0806.php?fid=21"
  ],
  "pageLimit": 30,

  "nameLenLimit": 60,
  "net": {
    "timeOut": 30,
    "proxyUrl": "socks5://localhost:1080"
  }
}
```

```go
go get github.com/wuxiangzhou2010/imagespider
cd $GOPATH/src/github.com/wuxiangzhou2010/imagespider
go run github.com/wuxiangzhou2010/imagespider
```

## arch

![arch](./mis/arch.png)
