package dlog_test

import (
	"dfra/dlog"
	"testing"
)

func TestStddlog(t *testing.T) {

	//测试 默认debug输出
	dlog.Debug("dfra debug content1")
	dlog.Debug("dfra debug content2")

	dlog.Debugf(" dfra debug a = %d\n", 10)

	//设置log标记位，加上长文件名称 和 微秒 标记
	dlog.ResetFlags(dlog.BitDate | dlog.BitLongFile | dlog.BitLevel)
	dlog.Info("dfra info content")

	//设置日志前缀，主要标记当前日志模块
	dlog.SetPrefix("MODULE")
	dlog.Error("dfra error content")

	//添加标记位
	dlog.AddFlag(dlog.BitShortFile | dlog.BitTime)
	dlog.Stack(" dfra Stack! ")

	//设置日志写入文件
	dlog.SetLogFile("./log", "testfile.log")
	dlog.Debug("===> dfra debug content ~~666")
	dlog.Debug("===> dfra debug content ~~888")
	dlog.Error("===> dfra Error!!!! ~~~555~~~")

	//关闭debug调试
	dlog.CloseDebug()
	dlog.Debug("===> 我不应该出现~！")
	dlog.Debug("===> 我不应该出现~！")
	dlog.Error("===> dfra Error  after debug close !!!!")
}
