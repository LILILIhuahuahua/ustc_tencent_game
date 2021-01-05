package network

import (
	"fmt"
	"github.com/LILILIhuahuahua/ustc_tencent_game/network/kcpnet"
	"github.com/xtaci/kcp-go"
	"log"
)

//网络组件----广播者
type(
	Broadcaster struct {
		addr string
		server *kcpnet.KcpServer
		conctmap map[interface{}]*Connector
	}
)
//数据持有；连接者指针列表
func NewBroadcaster(address string) (*Broadcaster,error){
	s,err := kcpnet.NewKcpServer(address)
	if err != nil {
		return nil,err
	}
	return &Broadcaster{
		addr: address,
		conctmap: make(map[interface{}]*Connector),
		server: s,
	},nil
}

//func(b *Broadcaster) Serv() error{
//	err := b.server.Run()
//	if err != nil {
//		return err
//	}
//	return nil
//}

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
			println(err)
			return err
		}
	}
 	return nil
}


func (b *Broadcaster) Serv() error {
	for {
		conn, err := b.server.Listen.AcceptKCP()
		if err != nil {
			return err
		}
		connector := NewConnector(b,*conn)
		err = b.RegisterConnector(connector)
		if err != nil {
			return err
		}
		go b.Handle(conn)
	}
}

func (b *Broadcaster) Handle(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		fmt.Println(string(buf))
		// 处理buf，改为事件驱动型...
		if err != nil {
			log.Println(err)
			return
		}
		err = b.NotifyAll(buf[:n])
		//n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

