package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/consul/v2"
	"strconv"
)

var reg registry.Registry

func init(){
	reg = consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
}

func main() {
	service := web.NewService(
		web.Name("math_service"),
		web.Address(":50000"),
		web.Handler(Initweb()),
		web.Registry(reg),
		)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Initweb() *gin.Engine{
	r := gin.Default()
	r.GET("/math/add",func(c *gin.Context){
		x,_ := strconv.Atoi(c.Query("x"))
		y,_ := strconv.Atoi(c.Query("y"))
		z := x+y
		c.String(200,fmt.Sprintf("z=%d",z))
	})
	return r
}