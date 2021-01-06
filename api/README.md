2020-12-31更新日志

1. 添加了HeroQuitNotify，当玩家退出游戏的时候客户端向服务端发送该request，服务器response之后再广播GameGlobalInfoNotify
2. EntityChangeRequest中eventType改成了枚举
3. ItemMsg中foodType改成枚举，改名成ItemType
4. Response.result改成枚举类型
5. playerSize改为float
6. 将所有和player或者snake有关的字段全部替换成了hero
7. 修改pb2为pb3



2022-1-1更新日志

1.GMessage.msgType改为MSG_TYPE



2022-1-6更新日志

1.Event_TYPE增加Hero移动事件，Hero变大事件

2.添加EnterGame包

3.添加ConnectNet结构