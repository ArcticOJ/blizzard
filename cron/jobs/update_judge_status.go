package jobs

import (
	"blizzard/cache/stores"
	"blizzard/config"
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
		RoutingKey  string `json:"routing_key"`
		Destination string `json:"destination"`
		Arguments   struct {
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

func groupBindings(b []binding) (m map[string][]stores.JudgeRuntime, runtimes []interface{}) {
	if len(b) == 0 {
		return nil, nil
	}
	m = make(map[string][]stores.JudgeRuntime)
	for _, b2 := range b {
		runtimes = append(runtimes, b2.RoutingKey)
		m[b2.Destination] = append(m[b2.Destination], stores.JudgeRuntime{
			ID:        b2.RoutingKey,
			Compiler:  b2.Arguments.Compiler,
			Arguments: b2.Arguments.Arguments,
			Version:   b2.Arguments.Version,
		})
	}
	return
}

func UpdateJudgeStatus(ctx context.Context) {
	status := make(map[string]stores.JudgeStatus)
	b := getBindings(ctx)
	q := getQueues(ctx)
	var judgeList []interface{}
	if q == nil {
		return
	}
	g, rts := groupBindings(b)
	for _, x := range q.Items {
		judgeList = append(judgeList, x.Arguments.Name)
		status[x.Arguments.Name] = stores.JudgeStatus{
			Version:     x.Arguments.Version,
			Memory:      x.Arguments.Memory,
			OS:          x.Arguments.OS,
			Parallelism: x.Arguments.Parallelism,
			BootedSince: x.Arguments.BootedSince,
			Runtimes:    g[x.Name],
		}
	}

	if buf, e := json.Marshal(status); e == nil {
		stores.Judge.UpdateJudgeStatus(ctx, judgeList, string(buf), rts)
		return
	}
	stores.Judge.UpdateJudgeStatus(ctx, judgeList, "{}", rts)
}
