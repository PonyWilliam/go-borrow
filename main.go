package main

import (
	"github.com/PonyWilliam/go-borrow/domain/repository"
	"github.com/PonyWilliam/go-borrow/domain/server"
	"github.com/PonyWilliam/go-borrow/handler"
	ProductLog "github.com/PonyWilliam/go-borrow/proto"
	common "github.com/PonyWilliam/go-common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	consul "github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/breaker/hystrix/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	"github.com/opentracing/opentracing-go"
	"strconv"
	"time"
)

func main() {
	consulConfig,err := common.GetConsualConfig("1.116.62.214",8500,"/micro/config")
	//配置中心
	if err != nil{
		log.Fatal(err)
	}
	consulRegistry := consul.NewRegistry(
		func(options *registry.Options){
			options.Addrs = []string{"1.116.62.214"}
			options.Timeout = time.Second * 10
		},
	)

	t,io,err := common.NewTracer("go.micro.service.borrow",":6834")
	if err != nil{
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	srv := micro.NewService(
		micro.Name("go.micro.service.borrow"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8086"),
		micro.Registry(consulRegistry),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(common.QPS)),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(hystrix.NewClientWrapper()),
		)
	srv.Init()
	go common.PrometheusBoot("5002")
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")
	db,err := gorm.Open("mysql",
		mysqlInfo.User+":"+mysqlInfo.Pwd+"@tcp("+mysqlInfo.Host + ":"+ strconv.FormatInt(mysqlInfo.Port,10) +")/"+mysqlInfo.DataBase+"?charset=utf8&parseTime=True&loc=Local",
	)
	if err != nil{
		log.Error(err)

	}
	defer db.Close()
	db.SingularTable(true)
	rp := repository.NewProductLogRepository(db)
	err = rp.InitTable()
	if err!=nil{
		log.Error(err)
	}
	service := server.NewWorkerService(rp)
	err = ProductLog.RegisterBorrowHandler(srv.Server(),&handler.ProductLog{service})
	if err:=srv.Run();err!=nil{
		log.Error(err)
	}

}