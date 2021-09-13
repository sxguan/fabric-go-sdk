# fabric-go-sdk
本项目基于hyperledger fabric 2.x网络
## 基本流程
### 拉取项目
`GOPATH`设置为`/root/go`
进入`GOPATH/src`

```
cd $GOPATH/src && git clone https://github.com/sxguan/fabric-go-sdk.git
```

### 启动节点

```
cd ./fabric-go-sdk/fixtures/ && docker-compose up -d
```

### 启动项目

```
cd .. && go build && ./fabric-go-sdk
```
```
>> 开始创建通道......
>>>> 使用每个org的管理员身份更新锚节点配置...
>>>> 使用每个org的管理员身份更新锚节点配置完成
>> 创建通道成功
>> 加入通道......
>> 加入通道成功
>> 开始打包链码......
>> 打包链码成功
>> 开始安装链码......
>> 安装链码成功
>> 组织认可智能合约定义......
>>> chaincode approved by Org1 peers:
	peer0.org1.example.com:7051
	peer1.org1.example.com:9051
>> 组织认可智能合约定义完成
>> 检查智能合约是否就绪......
LifecycleCheckCCCommitReadiness cc = simplecc, = {map[Org1MSP:true]}
LifecycleCheckCCCommitReadiness cc = simplecc, = {map[Org1MSP:true]}
>> 智能合约已经就绪
>> 提交智能合约定义......
>> 智能合约定义提交完成
>> 调用智能合约初始化方法......
>> 完成智能合约初始化
>> 通过链码外部服务设置链码状态......
>> 设置链码状态完成
<--- 添加信息　--->： 18c0c86ce029d7de04461484976c5151992864b52ca28905d0ccf911443fdfcb
<--- 查询信息　--->： 123
```

