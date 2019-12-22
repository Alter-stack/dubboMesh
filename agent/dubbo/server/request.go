package server


var atomicLong int64

type AgentRequest struct {
	RequestID int64
	Mdata RpcInvocation
}

func InitAgentRequest(requestID int64) *AgentRequest {
	a := &AgentRequest{
		RequestID: requestID,
	}
	return a
}

func (a *AgentRequest) setData(msg RpcInvocation) {
	a.Mdata = msg
}

type RpcInvocation struct {
	InterfaceName, MethodName, ParameterTypes string
	Attachments map[string]string
	Arguments []byte
}