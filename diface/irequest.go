package diface

//IRequset接口， 实际上是把客户端请求的连接信息  和请求的数据 包装到了一个Request中

type IRequset interface {
	//得到当前连接
	GetConnection() IConnection

	//得到请求的消息数据
	GetData() []byte

	//得到请求的消息id
	GetMsgId() uint32
}
