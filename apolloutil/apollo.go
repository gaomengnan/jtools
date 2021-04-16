package apolloutil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gaomengnan/jtools/apolloutil/proto"
	"google.golang.org/grpc"
	"time"
)

type Apollo struct {
	addr         string
	conn         *grpc.ClientConn
	apolloClient proto.ApolloClient
	subServer    proto.Apollo_SubscribeClient
	data         map[string]map[string]map[string]interface{}
}

func NewApollo(addr string) *Apollo {
	return &Apollo{
		addr: addr,
	}
}

func (apo *Apollo) Subscribe(key string) *Apollo {
	apo.connect()

	c, _ := apo.apolloClient.Subscribe(context.Background(), &proto.SubRequest{
		Key: key,
	})

	for {
		msg, err := c.Recv()
		if err != nil {
			break
		}
		if msg.Data != "" {
			_ = json.Unmarshal([]byte(msg.Data), &apo.data)
		}
		break
	}

	go func(c proto.Apollo_SubscribeClient, apollo *Apollo) {
		apo.receiveMsg()
	}(c, apo)
	apo.subServer = c
	return apo
}

func (apo *Apollo) receiveMsg() {

	for {
		time.Sleep(time.Microsecond * 300)
		msg, err := apo.subServer.Recv()
		if err != nil {
			break
		}

		if msg.Data != "" {
			_ = json.Unmarshal([]byte(msg.Data), &apo.data)
		}
	}

}

func (apo *Apollo) connect() {

	fmt.Println(apo.addr)
	conn, err := grpc.Dial(apo.addr, grpc.WithInsecure())
	if err != nil {
		return
	}
	//defer conn.Close()
	apo.conn = conn
	c := proto.NewApolloClient(conn)

	//reqstreamData := &proto.SubRequest{
	//	Key: "test.apollo-service-setting.application",
	//}
	//
	//res, _ := c.Subscribe(context.Background(), reqstreamData)
	//
	//for {
	//	aa, err := res.Recv()
	//	if err != nil {
	//		log.Println(err)
	//	}
	//	log.Println(aa)
	//}

	apo.apolloClient = c
}

func (apo *Apollo) Get() map[string]map[string]map[string]interface{} {
	return apo.data
}
