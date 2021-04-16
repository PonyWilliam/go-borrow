package repository

import (
	"bytes"
	"crypto/sha256"
	"errors"
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
	CheckToOther(ID int64,WID int64)error
	FindBorrowAll()([]model.ProductLog,error)
	FindBorrowByID(ID int64)(model.ProductLog,error)
	FindBorrowByWID(WID int64)([]model.ProductLog,error)
	FindBorrowByProductID(PID int64)([]model.ProductLog,error)
	TestLog()(id int64,err error)
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
	data1 := []byte(strconv.FormatInt(WID,10))
	data2 := []byte(strconv.FormatInt(PID,10))
	data3 := []byte(strconv.FormatInt(ScheduleTime,10))
	secret := []byte("RFIDeX")
	lastData := &model.ProductLog{}
	temp := p.mysql.Find(&lastData)
	var data []byte
	var prehash []byte
	if temp == nil {
		//没有prehash
		data = bytes.Join([][]byte{data1,data2,data3,secret},[]byte{})
	}else{
		prehash = lastData.Hash
		data = bytes.Join([][]byte{data1,data2,data3,prehash,secret},[]byte{})
	}
	hash := sha256.New().Sum(data)
	log := &model.ProductLog{BorrowTime: time.Now().Unix(),Pid: PID,Wid: WID,ScheduleTime: ScheduleTime,ReturnTime:0,
		Description: "首次借出",Hash: hash,PreHash: prehash}
	return log.Id,p.mysql.Model(log).Create(&log).Error
}
func(p *ProductLog)Return(ID int64)error{
	log := &model.ProductLog{}
	p.mysql.Where("id = ?",ID).Last(&log)
	if log.ReturnTime != 0{
		return errors.New("已归还过该物品")
	}

	return p.mysql.Model(&model.ProductLog{}).Where("id = ?",ID).Update("ReturnTime",time.Now().Unix()).Error
}
func(p *ProductLog)UpdateToOther(ID int64,WID int64)error{
	log := &model.ProductLog{}
	p.mysql.Where("id = ?",ID).First(&log)
	times := time.Now().Unix()
	err := p.mysql.Model(log).Where("id = ?",ID).Update("ReturnTime",times).Error
	if err != nil{
		return err
	}
	p.mysql.Where("id = ?",ID).Find(&log)
	log.Description = "转借于:" + strconv.FormatInt(log.Wid,10)
	log.BorrowTime = times
	log.ReturnTime = 0
	log.Wid = WID
	log.Id = 0
	return p.mysql.Model(log).Create(&log).Error
}
func(p *ProductLog)CheckToOther(ID int64,WID int64)error{
	log := &model.ProductLog{}
	p.mysql.Where("id = ?",ID).First(&log)
	if log.ReturnTime != 0{
		return errors.New("已归还物品不能转移")
	}
	if log.Wid == WID{
		return errors.New("不能转借给自己")
	}

	times := time.Now().Unix()
	if times - log.ReturnTime < 3600 * 3{
		return errors.New("距离预定归还时间少于三小时物品不能转借")
	}

	return nil//都没有则可以转。
}

/*
	FindBorrowAll()([]*model.ProductLog,error)
	FindBorrowByID(ID int64)(*model.ProductLog,error)
	FindBorrowByWID(WID int64)([]*model.ProductLog,error)
	FindBorrowByProductID(PID int64)([]*model.ProductLog,error)
*/
func(p *ProductLog)FindBorrowAll()(logs []model.ProductLog,err error){
	return logs,p.mysql.Model(&logs).Find(&logs).Error
}
func(p *ProductLog)FindBorrowByID(ID int64)(log model.ProductLog,err error){
	return log,p.mysql.Where("id = ?",ID).Find(&log).Error
}
func(p *ProductLog)FindBorrowByWID(WID int64)(logs []model.ProductLog,err error){
	return logs,p.mysql.Where("wid = ?",WID).Find(&logs).Error
}
func(p *ProductLog)FindBorrowByProductID(PID int64)(logs []model.ProductLog,err error){
	return logs,p.mysql.Where("pid = ?",PID).Find(&logs).Error
}
func(p *ProductLog)TestLog()(int64,error){
	var logs []model.ProductLog
	err := p.mysql.Find(&logs).Error
	secret := []byte("RFIDeX")
	var data []byte
	if err != nil{
		return 0,err
	}
	var prehash []byte
	for _,v := range logs{
		data1 := []byte(strconv.FormatInt(v.Wid,10))
		data2 := []byte(strconv.FormatInt(v.Pid,10))
		data3 := []byte(strconv.FormatInt(v.ScheduleTime,10))
		if sha256.Sum224(v.PreHash) == sha256.Sum224(prehash){
			//没有数据
			data = bytes.Join([][]byte{data1,data2,data3,secret},[]byte{})
		}else{
			//前面已有区块
			data = bytes.Join([][]byte{data1,data2,data3,[]byte(prehash[:]),secret},[]byte{})
		}
		if sha256.Sum256(v.Hash) != sha256.Sum256(sha256.New().Sum(data)){
			//校验出错
			return v.Id,nil
		}
	}
	//全部完成没出错
	return -1,nil
}