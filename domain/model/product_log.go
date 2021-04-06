package model
type ProductLog struct{
	Id int64 `gorm:"primary_key;not_null;auto_increment;" json:"ID"`
	Pid int64 `json:"pid"`//物品的id
	Wid int64 `json:"wid"`//借出员工ID
	BorrowTime int64 `json:"borrow_time"`//出借时间
	ScheduleTime int64 `json:"schedule_time"`//预定归还时间
	ReturnTime int64 `json:"ReturnTime"`//实际归还时间,给0代表是借出并为归还。
	Description string `json:"Description"`
	BelongArea int64 `json:"belong_area"`//所属库房
	PreHash [32]byte `json:"pre_hash"`//记录上一节点哈希
	Hash [32]byte `json:"hash"`
}