package dubbo

import (
	"dubboMesh/agent/dubbo/server"
	"fmt"
	"testing"
)

func TestEncodeRequest(t *testing.T) {
	message := &server.AgentRequest{
		RequestID: 123,
	}
	res, err := PackRequest(message)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", res)
	fmt.Printf("%c", res)
}