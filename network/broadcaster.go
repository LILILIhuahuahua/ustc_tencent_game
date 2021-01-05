package network

import "github.com/LILILIhuahuahua/ustc_tencent_game/network/kcp"

//网络组件----广播者
type(
	Broadcaster struct {
		addr string
		server *kcp.KcpServer
		conctmap map[interface{}]*Connector
	}
)
//数据持有；连接者指针列表
func NewBroadcaster(address string) (*Broadcaster,error){
	s,err := kcp.NewKcpServer(address)
	if err != nil {
		return nil,err
	}
	return &Broadcaster{
		addr: address,
		conctmap: make(map[interface{}]*Connector),
		server: s,
	},nil
}

func(b *Broadcaster) Serv() error{
	err := b.server.Run()
	if err != nil {
		return err
	}
	return nil
}

//注册连接者
func (b *Broadcaster)RegisterConnector(c *Connector)  error{
	c.broader = b
	b.conctmap[c.Id]=c
	return nil
}

//删除连接者
func (b *Broadcaster)DeleteConnector(c *Connector)  error{
return nil
}

//推送广播
func (b *Broadcaster)NotifyAll(buff []byte) error {
	for  _,connector := range b.conctmap{
		err:= connector.Update(buff)
		if nil != err {
			return err
		}
	}
 	return nil
}

