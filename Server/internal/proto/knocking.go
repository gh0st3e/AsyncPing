package knocking

import (
	"Server/internal/service"
	"context"
	"fmt"
	"strings"
	"time"
)

type KnockServer struct {
	service *service.Service
	UnimplementedKnockingServer
}

func NewKnockServer(service *service.Service) *KnockServer {
	return &KnockServer{service: service}
}

func checkError(err error) string {
	if strings.Contains(fmt.Sprintf("%s", err), "pq") {
		return "Mongo: OK, Postgres: TOO MANY CALLS"
	}
	return "Mongo: TOO MANY CALLS, Postgres: OK "
}

func (k *KnockServer) KnockDB(ctx context.Context, in *RequestParam) (*ResponseParam, error) {
	err := k.GoPing()
	if err != nil {
		fmt.Println(time.Now())
		fmt.Println(err)
		strError := checkError(err)
		fmt.Println(strError)

		return &ResponseParam{Msg: strError}, nil
	}
	return &ResponseParam{Msg: "Mongo: OK!, Postgres: OK!"}, nil
}

func (k *KnockServer) GoPing() error {
	err := k.service.PingMongo()
	if err != nil {
		return err
	}
	err = k.service.PingPSQL()
	if err != nil {
		return err
	}
	return nil
}
