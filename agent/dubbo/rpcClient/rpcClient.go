package rpcClient

import (
	"dubboMesh/agent/dubbo"
	"dubboMesh/agent/dubbo/server"
	"net"
)

type RpcClient struct {

}


func (r *RpcClient) Invoke(interfaceName, method, parameterTypesString, parameter string) ([]byte, error) {
	invocation := server.RpcInvocation{
		MethodName: method,
		ParameterTypes: parameterTypesString,
		Arguments: []byte(parameter),
	}
	invocation.Attachments = map[string]string{"path": interfaceName}


	req := server.InitAgentRequest()
	req.Mdata = invocation
	totalBuf, err := dubbo.PackRequest(req)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", "127.0.0.1:20880")
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(totalBuf)
	if err != nil {
		return nil, err
	}
	res := make([]byte, 2048)
	conn.Read(res)
	return res, nil
}