package main

import (
	"context"
	message "dubboMesh/agent/dubbo/server/pb"
	"dubboMesh/agent/register/etcdRegister"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"sync/atomic"

	//"net/http"
	"net/url"
	"sync"
)

var lock sync.Mutex
var requestID int64

func main() {
	r := gin.Default()
	//r.GET("/create/", Hello)
	r.POST("/agent", HelloController)

	// Run http server
	if err := r.Run(":8888"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}


func HelloController(ctx *gin.Context) {
	InterfaceName := ctx.PostForm("interface")
	method := ctx.PostForm("method")
	parameterTypesString := ctx.PostForm("parameterTypesString")
	parameter := ctx.PostForm("parameter")
	//agentType := os.Getenv("type")
	consumer(InterfaceName, method, parameterTypesString, parameter)
}



func consumer(interfaceName string, method string, parameterTypesString string, parameter string) int64 {
	// 组装 agent 请求
	lock.Lock()
	register := new(etcdRegister.EtcdManager)
	endpointList := register.Find(interfaceName)
	lock.Unlock()
	randIndex := rand.Intn(len(endpointList))
	endpoint := endpointList[randIndex] // 随机抽取负载均衡

	seviceUrl := new(url.URL)
	seviceUrl.Host = endpoint.GetHost() + ":" + string(endpoint.GetPort())

	conn, err := grpc.Dial(seviceUrl.Host)
	if err != nil {
		return 0
	}
	cli := message.NewAgentServiceClient(conn)
	req := message.AgentRequest{
		RequestID: atomic.AddInt64(&requestID, 1),
		Interface:interfaceName,
		Method:method,
		ParameterTypesString:parameterTypesString,
		Parameter:parameter,
	}

	resp , err := cli.Server(context.Background(), &req)
	if err != nil {
		return 0
	}

	//request := new(http.Request)
	//request.URL = seviceUrl
	//request.Form.Add("interface", interfaceName)
	//request.Form.Add("method", method)
	//request.Form.Add("parameterTypesString", parameterTypesString)
	//request.Form.Add("parameter", parameter)
	//
	//httpClient := new(http.Client)
	//resp ,err := httpClient.Do(request)
	//if err != nil {
	//	//	log.Fatal(err)
	//	//}
	//var respBytes []byte
	//s, err := resp.Body.Read(respBytes)
	return resp.RespLen

}
