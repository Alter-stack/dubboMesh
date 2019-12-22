package etcdRegister

import (
	"context"
	"dubboMesh/agent/register/endpoints"
	"dubboMesh/util"
	"fmt"
	etcd "github.com/coreos/etcd/clientv3"
	"log"
	"strconv"
	"strings"
	"time"
)


type EtcdManager struct {
	Host string
	Port int
}

func getInterfaceKey(serviceName string, port int) string {
	return serviceName + util.GetHostName() + ":" + strconv.Itoa(port)
}


func (e *EtcdManager) Register(serviceName string, port int) {
	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", e.Host, e.Port)},
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, _ := context.WithCancel(context.Background())

	// 创建一个5秒的租约
	leaseGrantResponse, err := cli.Grant(ctx, 5)
	if err != nil {
		log.Fatal(err)
	}

	interfaceKey := getInterfaceKey(serviceName, port)
	log.Println("register interface to etcd", interfaceKey)

	_, err = cli.Put(ctx, interfaceKey, "", etcd.WithLease(leaseGrantResponse.ID))
	if err != nil {
		log.Fatal(err)
	}

	// 续期保持节点存在
	ch, err := cli.KeepAlive(ctx, leaseGrantResponse.ID)
	if err != nil {
		log.Fatal(err)
	}
	ka := <- ch
	log.Printf("ttl: ", ka.TTL)

	// TODO LB 能力值的更新
	//for {
	//	time.Sleep(time.Second * 10)
	//	_, err = cli.Put(context.TODO(), interfaceKey, "", etcd.WithLease(leaseGrantResponse.ID))
	//	if err != nil {
	//		continue
	//	}
	//}
}

func (e *EtcdManager) Find(serviceName string) []*endpoints.EndPoint {
	var endpointList []*endpoints.EndPoint

	cli, err := etcd.New(etcd.Config{
		Endpoints:   []string{fmt.Sprintf("http://%s:%d", e.Host, e.Port)},
		DialTimeout: time.Second * 3,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	if resp, err := cli.Get(context.Background(), serviceName, etcd.WithPrefix()); err != nil {
		log.Fatal(err)
	} else {
		for _, kv := range resp.Kvs {
			keyInfo := string(kv.Key)
			index := strings.LastIndex(keyInfo, "/")
			tail := keyInfo[index:]
			keyValuePairs := strings.Split(tail, ":")
			addr := keyValuePairs[0]
			port, _ := strconv.Atoi(keyValuePairs[1])
			endpoint := endpoints.NewEndPoint(addr, port)
			endpointList = append(endpointList, endpoint)
		}
	}
	return endpointList
}
