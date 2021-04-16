package handler

import (
	"context"
	"encoding/hex"
	"github.com/PonyWilliam/go-borrow/domain/model"
	"github.com/PonyWilliam/go-borrow/domain/server"
	borrow "github.com/PonyWilliam/go-borrow/proto"
	"strconv"
)

/*
rpc Borrow(Request)returns(Response);
rpc Return(Request)returns(Response);
rpc ToOther(to_other_request)returns(Response);
*/
type ProductLog struct{
	IProService server.IProService
}
func(p *ProductLog)Borrow(ctx context.Context,req *borrow.Borrow_Request,rsp *borrow.Response)error{
	id,err := p.IProService.Borrow(req.WorkerId,req.ProductId,req.ScheduleTime)
	if err!=nil{
		rsp.Message = err.Error()
		rsp.Status = 0
		return nil
	}
	rsp.Message = strconv.FormatInt(id,10)
	rsp.Status = 1
	return nil
}
func(p *ProductLog)Return(ctx context.Context,req *borrow.Returns_Request,rsp *borrow.Response)error{
	err := p.IProService.Return(req.Id)
	if err!=nil{
		rsp.Message = err.Error()
		rsp.Status = 0
		return nil
	}
	rsp.Message = "success return"
	rsp.Status = 1
	return nil
}
func(p *ProductLog)ToOther(ctx context.Context,req *borrow.ToOtherRequest,rsp *borrow.Response)error{
	err := p.IProService.UpdateToOther(req.Id,req.Wid)
	if err!=nil{
		rsp.Message = err.Error()
		rsp.Status = 0
		return nil
	}
	rsp.Message = "success update"
	rsp.Status = 1
	return nil
}
func(p *ProductLog)CheckToOther(ctx context.Context,req *borrow.ToOtherRequest,rsp *borrow.Response)error{
	err := p.IProService.CheckToOther(req.Id,req.Wid)
	if err!=nil{
		rsp.Message = err.Error()
		rsp.Status = 0
		return nil
	}
	rsp.Message = "success update"
	rsp.Status = 1
	return nil
}
func(p *ProductLog)FindBorrowAll(ctx context.Context,req *borrow.Null_Request,rsp *borrow.Borrowlogs_Response)error{
	res,err := p.IProService.FindBorrowAll()
	if err!=nil{
		rsp.Logs = nil
		return nil
	}
	for _,v := range res{
		rsp.Logs = append(rsp.Logs,Swap(v))
	}
	return nil
}
func(p *ProductLog)FindBorrowByID(ctx context.Context,req2 *borrow.ID_Request,rsp *borrow.Borrowlog_Response)error{
	req,err := p.IProService.FindBorrowByID(req2.Id)
	if err!=nil{
		rsp = nil
		return nil
	}
	rsp.ID = req.Id
	rsp.ReturnTime = req.ReturnTime
	rsp.BorrowTime = req.BorrowTime
	rsp.PID = req.Pid
	rsp.WID = req.Wid
	rsp.Description = req.Description
	rsp.PreHash = hex.EncodeToString(req.PreHash)
	rsp.Hash = hex.EncodeToString(req.Hash)
	rsp.BelongArea = req.BelongArea
	rsp.ScheduleTime = req.ScheduleTime
	return nil
}
func(p *ProductLog)FindBorrowByWID(ctx context.Context,req *borrow.WID_Request,rsp *borrow.Borrowlogs_Response)error{
	res,err := p.IProService.FindBorrowByWID(req.WID)
	if err!=nil{
		rsp.Logs = nil
		return err
	}
	for _,v := range res{
		rsp.Logs = append(rsp.Logs,Swap(v))
	}
	return nil
}
func(p *ProductLog)FindBorrowByPID(ctx context.Context,req *borrow.PID_Request,rsp *borrow.Borrowlogs_Response)error{
	res,err := p.IProService.FindBorrowByProductID(req.PID)
	if err!=nil{
		rsp.Logs = nil
		return err
	}
	for _,v := range res{
		rsp.Logs = append(rsp.Logs,Swap(v))
	}
	return nil
}
func(p *ProductLog)TestLog(ctx context.Context,req *borrow.Null_Request,rsp *borrow.Response_HashTest)error{
	var err error
	rsp.Id,err = p.IProService.TestLog()
	if err != nil{
		return err
	}
	return nil
}
func Swap(req model.ProductLog)*borrow.Borrowlog_Response{
	temp := borrow.Borrowlog_Response{}
	temp.ID = req.Id
	temp.ReturnTime = req.ReturnTime
	temp.BorrowTime = req.BorrowTime
	temp.PID = req.Pid
	temp.WID = req.Wid
	temp.Description = req.Description
	temp.PreHash = hex.EncodeToString(req.PreHash)
	temp.Hash = hex.EncodeToString(req.Hash)
	temp.BelongArea = req.BelongArea
	temp.ScheduleTime = req.ScheduleTime
	return &temp
}