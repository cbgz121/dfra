package pack

import "dfra/diface"

type Message struct {
	id      uint32 //消息ID
	dataLen uint32 //消息长度
	data    []byte //消息的内容
}

func NewMsgPackage(id uint32, data []byte) diface.IMessage {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

// 获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.id
}

// 获取消息的长度
func (m *Message) GetMsgLen() uint32 {
	return m.dataLen
}

// 获取消息的内容
func (m *Message) GetData() []byte {
	return m.data
}

// 设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.id = id
}

// 设置消息的内容
func (m *Message) SetData(data []byte) {
	m.data = data
}

// 设置消息的长度
func (m *Message) SetMsgLen(datalen uint32) {
	m.dataLen = datalen
}
