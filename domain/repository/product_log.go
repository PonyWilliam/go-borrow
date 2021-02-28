package repository

import (
	"errors"
	"fmt"
	"github.com/PonyWilliam/go-borrow/domain/model"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)
type IProductLog interface {
	InitTable() error
	Borrow(WID int64,PID int64,ScheduleTime int64) (id int64,err error)
	Return(ID int64)error
	UpdateToOther(ID int64,WID int64)error//转借给其它人的记录
}
func NewProductLogRepository(db *gorm.DB)	IProductLog{
	return &ProductLog{mysql: db}
}
type ProductLog struct{
	mysql *gorm.DB
}
func(p *ProductLog)InitTable()error{
	if p.mysql.HasTable(&model.ProductLog{}){
		return nil
	}
	return p.mysql.CreateTable(&model.ProductLog{}).Error
}
func(p *ProductLog)Borrow(WID int64,PID int64,ScheduleTime int64)(id int64,err error){
	log := &model.ProductLog{BorrowTime: time.Now().Unix(),PID: PID,WID: WID,ScheduleTime: ScheduleTime,ReturnTime:0,Description: "首次借出"}
	return log.ID,p.mysql.Model(log).Create(&log).Error
}
func(p *ProductLog)Return(ID int64)error{
	return p.mysql.Model(&model.ProductLog{}).Where("id = ?",ID).Update("ReturnTime",time.Now().Unix()).Error
}
func(p *ProductLog)UpdateToOther(ID int64,WID int64)error{
	log := &model.ProductLog{}
	p.mysql.Where("id = ?",ID).First(&log)

	fmt.Println(log)
	if log.ReturnTime != 0{
		return errors.New("已归还物品不能转移")
	}
	if log.WID == WID{
		return errors.New("不能转借给自己")
	}

	times := time.Now().Unix()
	if times - log.ReturnTime < 3600 * 3{
		return errors.New("距离预定归还时间少于三小时物品不能再次借出")
	}
	err := p.mysql.Model(log).Where("id = ?",ID).Update("ReturnTime",times).Error
	if err != nil{
		return err
	}
	p.mysql.Where("id = ?",ID).Find(&log)
	log.Description = "转借于:" + strconv.FormatInt(log.WID,10)
	log.BorrowTime = times
	log.ReturnTime = 0
	log.WID = WID
	log.ID = 0
	return p.mysql.Model(log).Create(&log).Error
}