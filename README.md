# Single dgs

This branch deals with game service.
## Layout explained

- model
    放置实体模型
  
- service
    dgs对外暴露的接口
  
- network
    网络组件，收发KCP以及广播。


## Q&A


1. Check meaning of project layout from [here](https://github.com/golang-standards/project-layout) ,try to follow the standard as possible as you can.

## 仓库地址汇总
- 登录服务：https://github.com/imilano/auth
- 数据库服务: https://github.com/TianqiS/database_for_happyball

## 启动方法
在启动服务的时候需要在参数中指定grpc的服务地址以及端口，否则服务器启动不成功
格式如下`program -DBProxyHost localhost -DBProxyPort 8890`