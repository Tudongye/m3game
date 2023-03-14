package asyncser

import (
	"m3game/example/proto/pb"
)

var (
	_channmsgs []*pb.ChannelMsg
)

func AppendMsg(msg *pb.ChannelMsg) {
	_channmsgs = append(_channmsgs, msg)
}

func GetMsg() []*pb.ChannelMsg {
	return _channmsgs
}
