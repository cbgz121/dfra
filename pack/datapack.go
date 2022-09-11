package pack

import (
	"bytes"
	"dfra/diface"
	"dfra/untils"
	"encoding/binary"
	"errors"
)

// 封包，拆包
type DataPack struct{}

// 拆包封包实例的一个初始化方法
func NewDataPack() diface.IDataPack {
	return &DataPack{}
}

// 获取包的头的长度的方法
func (d *DataPack) GetHeadLen() uint32 {
	//暂时固定return 8
	return 8
}

// 封包方法
func (d *DataPack) Pack(msg diface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法(将包的Head信息都读出来) 之后再根据信息里的data长度，再进行一次读
func (d *DataPack) Unpack(binaryData []byte) (diface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.id); err != nil {
		return nil, err
	}

	if untils.GlobalObject.MaxPackageSize > 0 && msg.dataLen > untils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil
}
