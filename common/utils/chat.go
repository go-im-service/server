package utils

import (
	"fmt"
	"server/common/model"
)

const (
	UserSecondCount      = 100   // 用户间相互聊天，默认一秒最多99条消息，超过会产生相同的聊天记录ID
	GroupSecondCount     = 10000 // 用户间相互聊天，默认一秒最多9999条消息
	BroadcastSecondCount = 1     // 每秒最多广播一次，(业务层需要根据情况控制，理论上可能一天只允许广播一次)
)

// GetChatMsgID 获取聊天消息的ID，会为每一个聊天记录产生一个当前聊天下的唯一ID(并非全局ID)
// 为了可以让每个聊天下的ID都唯一，方便聊天记录查询，删除，客户端ACK消息等
// ID根据时间生成，保证单调递增，同时记录一个counter值(可以借助redis实现)来保证唯一性。
// 如果聊天每秒内产生的聊天记录大于secondCount，就会出现重复的MsgID。
// 如果群太大，或者支持并发推送消息，可以考虑加大secondCount，或者使用纳秒。也可以做唯一性检验，对重复ID拒绝发送
func GetChatMsgID(time, counter, secondCount int64) int64 {
	return time*secondCount + counter%secondCount
}

// GetChatID 获取聊天ID
// 为每一个单人聊天，群组聊天，系统通知聊天等生成一个唯一的聊天ID，用来标识
// appID为必须字段，每个业务通过appID区别
// 如果是一对一聊天，则使用两端的id号拼接，type和clientID组合小的在前边
// 如果是群组聊天或者广播号，则使用群组的type+id来标识
func GetChatID(appID int, type1 int, clientID1 int64, type2 int, clientID2 int64) string {
	if type1 == model.ClientIdentityTypeGroup {
		return fmt.Sprintf("%v:group:%v:%v", appID, type1, clientID1)
	}
	if type2 == model.ClientIdentityTypeGroup {
		return fmt.Sprintf("%v:group:%v:%v", appID, type2, clientID2)
	}
	if type1 == model.ClientIdentityTypeSysBroadcast {
		return fmt.Sprintf("%v:broadcast:%v:%v", appID, type1, clientID1)
	}
	if type2 == model.ClientIdentityTypeSysBroadcast {
		return fmt.Sprintf("%v:broadcast:%v:%v", appID, type2, clientID2)
	}

	if type1 < type2 {
		return fmt.Sprintf("%v:single:%v:%v:to:%v:%v", appID, type1, clientID1, type2, clientID2)
	}
	if type1 > type2 {
		return fmt.Sprintf("%v:single:%v:%v:to:%v:%v", appID, type2, clientID2, type1, clientID1)
	}

	// type1 == type2
	if clientID1 < clientID2 {
		return fmt.Sprintf("%v:single:%v:%v:to:%v:%v", appID, type1, clientID1, type2, clientID2)
	}

	return fmt.Sprintf("%v:single:%v:%v:to:%v:%v", appID, type2, clientID2, type1, clientID1)
}

// GetChatDirectionID 获取聊天方向id
// 例如：A和B聊天，AB对于聊天记录的读取位置、配置等是不一样的，需要分别标记
// 采用：appID:type:id:type:id的形式
func GetChatDirectionID(appID int, type1 int, clientID1 int64, type2 int, clientID2 int64) string {
	return fmt.Sprintf("%v:%v:%v:%v:%v", appID, type1, clientID1, type2, clientID2)
}
