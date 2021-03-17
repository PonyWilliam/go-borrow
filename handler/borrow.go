package handler

import (
	"context"
	"github.com/PonyWilliam/go-borrow/domain/server"
	borrow "github.com/PonyWilliam/go-borrow/proto"
	"github.com/PonyWilliam/go-common"
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
	_ = common.SwapTo(res, rsp)
	return nil
}
func(p *ProductLog)FindBorrowByID(ctx context.Context,req *borrow.ID_Request,rsp *borrow.Borrowlog_Response)error{
	res,err := p.IProService.FindBorrowByID(req.Id)
	if err!=nil{
		rsp = nil
		return nil
	}
	_ = common.SwapTo(res, rsp)
	return nil
}
func(p *ProductLog)FindBorrowByWID(ctx context.Context,req *borrow.WID_Request,rsp *borrow.Borrowlogs_Response)error{
	res,err := p.IProService.FindBorrowByWID(req.WID)
	if err!=nil{
		rsp.Logs = nil
		return nil
	}
	_ = common.SwapTo(res, rsp)
	return nil
}
func(p *ProductLog)FindBorrowByPID(ctx context.Context,req *borrow.PID_Request,rsp *borrow.Borrowlogs_Response)error{
	res,err := p.IProService.FindBorrowByProductID(req.PID)
	if err!=nil{
		rsp.Logs = nil
		return nil
	}
	_ = common.SwapTo(res, rsp)
	return nil
}