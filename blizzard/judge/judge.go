package judge

import (
	"blizzard/blizzard/config"
	"blizzard/blizzard/db/models/contest"
	"blizzard/blizzard/logger"
	models "blizzard/blizzard/models"
	"blizzard/blizzard/pb"
	"blizzard/blizzard/server/utils"
	"context"
	"fmt"
	"go.arsenm.dev/drpc/muxconn"
	"net"
	"storj.io/drpc/drpcmetadata"
	"sync/atomic"
	"time"
)

var judges []*Client

var currentIndex uint32

type Info struct {
	Name        string  `json:"name"`
	IsAlive     bool    `json:"isAlive"`
	Latency     float64 `json:"latency"`
	BootedSince int64   `json:"bootedSince"`
	*pb.InstanceSpecification
}

func GetInfo(ctx context.Context) (info []*Info) {
	for _, client := range judges {
		alive, latency := client.Ping(keyedContext(ctx, client))
		inf := &Info{Name: client.Name, Latency: latency, IsAlive: alive, BootedSince: utils.BootTimestamp.Unix(), InstanceSpecification: client.Specs}
		inf.BootTimestamp = nil
		info = append(info, inf)
	}
	return
}

func initJudges(_judges map[string]*models.Judge) {
	for name, judge := range _judges {
		client := &Client{
			rpc:        nil,
			Name:       name,
			address:    judge.Address,
			privateKey: judge.Key,
		}
		judges = append(judges, client)
		dial, e := net.DialTimeout("tcp", judge.Address, time.Millisecond*300)
		if e != nil {
			logger.Logger.Err(e)
			continue
		}
		conn, _ := muxconn.New(dial)
		client.rpc = pb.NewDRPCIglooClient(conn)
		if specs, e := client.rpc.Specification(keyedContext(context.Background(), client), nil); e == nil {
			client.Specs = specs
		}
	}
}

func next() *Client {
	n := atomic.AddUint32(&currentIndex, 1)
	return judges[(int(n)-1)%len(judges)]
}

func processResult(r *pb.JudgeResult) interface{} {
	if r.GetCase() != nil {
		return r.GetCase()
	}
	if r.GetFinal() != nil {
		return finalResult{Verdict: processFinalVerdict(nil, nil)}
	}
	return nil
}

func judge(client *Client, submission *pb.Submission) {
	stream, e := client.rpc.Judge(keyedContext(context.Background(), client), submission)
	for {
		r, e := stream.Recv()
		if e != nil {
			// TODO: handle stream error
			break
		}
		ResultWatcher.emit(submission.Id, r)
		if r.GetFinal() != nil {
			ResultWatcher.Dispose(submission.Id)
			break
		}
		fmt.Println(r)
	}
	if e != nil {
		logger.Logger.Err(e)
		return
	}
}

func Enqueue(sub *pb.Submission, _sub *contest.Submission) func() {
	var c *Client
	i := 0
	for {
		if i == len(judges) {
			// stop as no judges have been found after one cycle
			c = nil
			return nil
		}
		c = next()
		if checkAlive(c) {
			break
		}
		i++
	}
	return func() {
		go judge(c, sub)
	}
}

func keyedContext(ctx context.Context, client *Client) context.Context {
	return drpcmetadata.Add(ctx, "key", client.privateKey)
}

func Init() {
	initJudges(config.Config.Judges)
}
