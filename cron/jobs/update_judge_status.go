package jobs

import (
	"blizzard/config"
	"blizzard/judge"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 2,
}

type (
	queues struct {
		Items []struct {
			Arguments struct {
				Name        string
				BootedSince uint64
				Memory      uint32
				Parallelism uint8
				OS          string
				Version     string
			} `json:"arguments"`
			Name string `json:"name"`
		} `json:"items"`
	}
	binding struct {
		RoutingKey string `json:"routing_key"`
		Arguments  struct {
			Name      string
			Arguments string
			Compiler  string
			Version   string
		} `json:"arguments"`
	}
)

func getQueues(ctx context.Context) *queues {
	conf := config.Config.RabbitMQ
	var res queues
	req, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d/api/queues/%s?page_size=500&page=1&name=judge-worker-.*-\\d{13}&use_regex=true&pagination=true", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	if e != nil {
		return nil
	}
	req.SetBasicAuth(conf.Username, conf.Password)
	r, e := client.Do(req)
	if e != nil {
		return nil
	}
	if json.NewDecoder(r.Body).Decode(&res) != nil {
		return nil
	}
	return &res
}

func getBindings(ctx context.Context) []binding {
	conf := config.Config.RabbitMQ
	var res []binding
	req, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d/api/exchanges/%s/submissions/bindings/source", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	if e != nil {
		return nil
	}
	req.SetBasicAuth(conf.Username, conf.Password)
	r, e := client.Do(req)
	if e != nil {
		return nil
	}
	if json.NewDecoder(r.Body).Decode(&res) != nil {
		return nil
	}
	return res
}

func groupBindings(b []binding) (m map[string][]judge.Runtime) {
	if len(b) == 0 {
		return nil
	}
	m = make(map[string][]judge.Runtime)
	for _, b2 := range b {
		m[b2.Arguments.Name] = append(m[b2.Arguments.Name], judge.Runtime{
			ID:        b2.RoutingKey,
			Compiler:  b2.Arguments.Compiler,
			Arguments: b2.Arguments.Arguments,
			Version:   b2.Arguments.Version,
		})
	}
	return m
}

func UpdateJudgeStatus(ctx context.Context) {
	judge.LockStatus()
	defer judge.UnlockStatus()
	for name := range judge.Status {
		judge.Status[name] = new(judge.Judge)
	}
	b := getBindings(ctx)
	q := getQueues(ctx)
	if q == nil {
		return
	}
	for _, x := range q.Items {
		judge.Status[x.Arguments.Name] = &judge.Judge{
			Alive:   true,
			Version: x.Arguments.Version,
			Info: &judge.Info{
				Memory:      x.Arguments.Memory,
				OS:          x.Arguments.OS,
				Parallelism: x.Arguments.Parallelism,
				BootedSince: x.Arguments.BootedSince,
			},
			Runtimes: nil,
		}
	}
	g := groupBindings(b)
	for name, runtime := range g {
		if j, ok := judge.Status[name]; ok {
			j.Runtimes = runtime
		}
	}
}
