package dnet

import "dfra/diface"

type Requset struct {
	conn diface.IConnection //已经和客户端建立好的连接
	data diface.IMessage    //客户端请求的数据
}

func NewRequest(c diface.IConnection, msg diface.IMessage) diface.IRequset {
	return &Requset{
		conn: c,
		data: msg,
	}
}

// 得到当前连接
func (r *Requset) GetConnection() diface.IConnection {
	return r.conn
}

// 得到请求的消息数据
func (r *Requset) GetData() []byte {
	return r.data.GetData()
}

func (r *Requset) GetMsgId() uint32 {
	return r.data.GetMsgId()
}
