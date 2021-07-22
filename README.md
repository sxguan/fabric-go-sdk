# fabric-go-sdk

`GOPATH`设置为`root/go`
进入`GOPATH/src`拉取项目

```
cd GOPATH/src && git clone https://github.com/sxguan/fabric-go-sdk.git
```

启动节点

```
cd /fabric-go-sdk/fixtures/ && docker-compose up -d
```

启动项目

```
cd .. && go build && ./fabric-go-sdk
```

