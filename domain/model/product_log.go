package model
type ProductLog struct{
	ID int64 `gorm:"primary_key;not_null;auto_increment;" json:"ID"`
	PID int64 `json:"pid"`//物品的id
	WID int64 `json:"wid"`//借出员工ID
	BorrowTime int64 `json:"borrow_time"`//出借时间
	ScheduleTime int64 `json:"schedule_time"`//预定归还时间
	ReturnTime int64 `json:"ReturnTime"`//实际归还时间,给0代表是借出并为归还。
	Description string `json:"Description"`
	BelongArea int64 `json:"belong_area"`//所属库房
	HashCode string `json:"hash_code"`//为了让管理人员不作弊的校验码,但实际上如果删库就会没有意义，除非与员工出借记录分开管理。所以还是得用区块链保证无法作弊。
}