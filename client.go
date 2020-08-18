package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client/selector"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/consul/v2"
	"io/ioutil"
	"net/http"
)

var reg1 registry.Registry

func init(){
	reg1 = consul.NewRegistry(registry.Addrs("127.0.0.1:8500"))
}

func main() {
	service := web.NewService(
		web.Name("other_service"),
		web.Address(":50001"),
		web.Handler(Initweb1()),
		web.Registry(reg1),
	)
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func Initweb1() *gin.Engine{
	r := gin.Default()
	r.GET("/other/add",func(c *gin.Context){
		content := Call(c.Query("x"),c.Query("y"))
		c.String(200,fmt.Sprintf("client: %s",content))
	})
	return r
}

func Call(x,y string) string{
	address := GetServiceAddress("math_service")
	url := fmt.Sprintf("http://"+address+"/math/add?x=%s&y=%s",x,y)
	response,_ := http.Get(url)
	defer response.Body.Close()
	content,_ := ioutil.ReadAll(response.Body)
	return string(content)
}

func GetServiceAddress(name string) (address string) {
	list,_ := reg1.GetService(name)
	var services []*registry.Service
	for _,value := range list {
		services = append(services,value)
	}
	next := selector.RoundRobin(services)
	if node,err := next(); err == nil{
		address = node.Address
	}
	return
}