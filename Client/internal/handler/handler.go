package handler

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strconv"
	"time"

	knocking "Client/internal/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type KnockServer struct {
}

func NewKnockServer() *KnockServer {
	return &KnockServer{}
}

func Mount(r *gin.Engine) {
	r.GET("/Health/:count", CallGRPC)
}

func CallGRPC(ctx *gin.Context) {
	count := ctx.Param("count")
	intCount, err := strconv.Atoi(count)
	if err != nil {
		log.Printf("Parse string to int error: %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	conn, err := grpc.Dial(":9000", grpc.WithInsecure(), grpc.WithNoProxy())
	if err != nil {
		log.Printf("Couldn't connect to gRPC: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}
	defer conn.Close()

	server := NewKnockServer()

	client := knocking.NewKnockingClient(conn)
	err = server.AsyncKnocking(ctx, int32(intCount), client)
	if err != nil {
		log.Printf("Ping DBs Error: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s", err),
		})
		return
	}

	log.Println("Everything OK")
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"message": "Mongo: OK!, Postgres: OK!" + fmt.Sprintf(" Pings - %v", intCount),
	})
}

func (k *KnockServer) Worker(ctx context.Context, jobs <-chan int, results chan<- *knocking.ResponseParam, client knocking.KnockingClient) error {
	for j := range jobs {
		select {
		case <-ctx.Done():
			return nil
		default:
			s, err := client.KnockDB(ctx, &knocking.RequestParam{Count: int32(1)})
			fmt.Println(fmt.Sprintf("Ping: %d = %s --- err: %v", j, s, err))
			if err != nil {
				return err
			}
			results <- s
		}
	}
	return nil
}

func (k *KnockServer) AsyncKnocking(ctx context.Context, count int32, client knocking.KnockingClient) error {
	t1 := time.Now()

	g, ctx := errgroup.WithContext(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan int, int(count))
	results := make(chan *knocking.ResponseParam, int(count))

	go FillJobs(ctx, jobs, int(count))

	for w := 1; w <= 100; w++ {
		g.Go(func() error {
			err := k.Worker(ctx, jobs, results, client)
			if err != nil {
				return err
			}
			return nil
		})
	}

	go func() {
		if err := g.Wait(); err != nil {
			fmt.Println(err)
		}
		close(results)
	}()

	for r := range results {
		if r.Msg != "Mongo: OK!, Postgres: OK!" {
			cancel()
			return fmt.Errorf(r.Msg)
		}
	}

	t2 := time.Since(t1)
	fmt.Println(t2)
	return nil
}

func FillJobs(ctx context.Context, jobs chan<- int, count int) {
	defer close(jobs)

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			jobs <- i
		}
	}
}
