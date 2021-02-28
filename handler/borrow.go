package handler

import (
	"context"
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
		return err
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
		return err
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
		return err
	}
	rsp.Message = "success update"
	rsp.Status = 1
	return nil
}