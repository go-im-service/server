package model

const (
	ClientIdentityTypeUser         = 1 // 一般的聊天用户
	ClientIdentityTypeGroup        = 2 // 群
	ClientIdentityTypeSysBroadcast = 5 // app的广播号，app发布自己的消息使用。支持向所有用户广播消息
	ClientIdentityTypeSysSingle    = 6 // app的单播号，app向单个用户发布消息使用，比如点赞、评论的通知
)

type ClientIdentity struct {
	AppID    int   `json:"app_id"`
	Type     int   `json:"type"`
	ClientID int64 `json:"client_id"`
}
