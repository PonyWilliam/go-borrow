package server

import "github.com/PonyWilliam/go-borrow/domain/repository"

type IProService interface {
	Borrow(WID int64,PID int64,ScheduleTime int64) (id int64,err error)
	Return(ID int64)error
	UpdateToOther(ID int64,WID int64)error//转借给其它人的记录
}
func NewWorkerService(pro repository.IProductLog)IProService{
	return &ProServices{pro}
}
type ProServices struct{
	pro repository.IProductLog
}
func(p *ProServices)Borrow(WID int64,PID int64,ScheduleTime int64)(int64,error){
	return p.pro.Borrow(WID,PID,ScheduleTime)
}
func(p *ProServices)Return(ID int64)error{
	return p.pro.Return(ID)
}
func(p *ProServices)UpdateToOther(ID int64,WID int64)error{
	return p.pro.UpdateToOther(ID,WID)
}
