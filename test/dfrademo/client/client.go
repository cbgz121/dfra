package main

import (
	"bytes"
	"dfra/diface"
	"dfra/pack"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type Message struct {
	id      uint32
	dataLen uint32
	data    []byte
}

func Pack(msg diface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	//将dataLen 写进打databuf中,小端模式
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将MsgId 写进打databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	//将data数据写入databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func heartBeat2(conn net.Conn) {
	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()

OverHeartBeat:
	for {
		time.Sleep(5 * time.Second)
		binMsg, err := Pack(pack.NewMsgPackage(1, []byte("dfra client Test Message")))
		if err != nil {
			fmt.Println("data pack err:", err)
			continue
		}
		conn.Write(binMsg)

		select {
		// 20s后退出心跳包
		case <-timer.C:
			break OverHeartBeat
		default:
		}
	}
}

func main() {
	fmt.Println("client start...")

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	// 客户端循环发心跳包
	go heartBeat2(conn)

	for {

		//发送封包的message消息  MsgID：0
		dp := pack.NewDataPack()
		binaryMsg, err := dp.Pack(pack.NewMsgPackage(1, []byte("dfra client Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write error", err)
			return
		}

		//先读取流中的head部分  得到ID 和 datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		//将二进制的head拆包到msg 结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}

		//再根据Datalen进行第二次读取  将data读出来
		var data []byte
		if msgHead.GetMsgLen() > 0 {
			data = make([]byte, msgHead.GetMsgLen())
			if io.ReadFull(conn, data); err != nil {
				fmt.Println("read msg data error", err)
				break
			}
		}

		fmt.Println("-->Recv Server Msg:ID = ", msgHead.GetMsgId(), "len = ", msgHead.GetMsgLen(), "data = ", string(data))
		time.Sleep(1 * time.Second)
	}
}
