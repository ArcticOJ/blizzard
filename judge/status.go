package judge

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ArcticOJ/blizzard/v0/config"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/yudppp/throttle"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// TODO: group

var (
	client = http.Client{
		Timeout: time.Second,
	}
	currentStatus     = make(map[string]status)
	supportedRuntimes = mapset.NewSet[string]()
	t                 = throttle.New(time.Second)
	m                 sync.RWMutex
)

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

	status struct {
		Name        string    `json:"name"`
		Version     string    `json:"version"`
		Memory      uint32    `json:"memory"`
		OS          string    `json:"os"`
		Parallelism uint8     `json:"parallelism"`
		BootedSince uint64    `json:"bootedSince"`
		Runtimes    []runtime `json:"runtimes"`
	}
	runtime struct {
		ID        string `json:"id"`
		Compiler  string `json:"compiler"`
		Arguments string `json:"arguments"`
		Version   string `json:"version"`
	}
)

func getQueues(ctx context.Context) *queues {
	conf := config.Config.RabbitMQ
	var q queues
	req, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d/api/queues/%s?page_size=500&page=1&name=judge-worker-.*-\\d{13}&use_regex=true&pagination=true", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	if e != nil {
		return nil
	}
	req.SetBasicAuth(conf.Username, conf.Password)
	r, e := client.Do(req)
	if e != nil {
		return nil
	}
	if json.NewDecoder(r.Body).Decode(&q) != nil {
		return nil
	}
	return &q
}

func getAllBindings(ctx context.Context) (m map[string][]runtime) {
	m = make(map[string][]runtime)
	conf := config.Config.RabbitMQ
	var bl []binding
	req, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d/api/exchanges/%s/submissions/bindings/source", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost)), nil)
	if e != nil {
		return
	}
	req.SetBasicAuth(conf.Username, conf.Password)
	r, e := client.Do(req)
	if e != nil {
		return
	}
	if json.NewDecoder(r.Body).Decode(&bl) != nil {
		return
	}
	for _, _b := range bl {
		m[_b.Destination] = append(m[_b.Destination], runtime{
			ID:        _b.RoutingKey,
			Compiler:  _b.Arguments.Compiler,
			Arguments: _b.Arguments.Arguments,
			Version:   _b.Arguments.Version,
		})
	}
	return
}

func getBinding(ctx context.Context, queue string) (rt []runtime) {
	conf := config.Config.RabbitMQ
	var _b []binding
	req, e := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("http://%s:%d/api/bindings/%s/e/submissions/q/%s", conf.Host, conf.ManagerPort, url.PathEscape(conf.VHost), queue), nil)
	if e != nil {
		return
	}
	req.SetBasicAuth(conf.Username, conf.Password)
	r, e := client.Do(req)
	if e != nil {
		return
	}
	if json.NewDecoder(r.Body).Decode(&_b) != nil {
		return
	}
	rt = make([]runtime, len(_b))
	for i := range _b {
		rt[i] = runtime{
			ID:        _b[i].RoutingKey,
			Compiler:  _b[i].Arguments.Compiler,
			Arguments: _b[i].Arguments.Arguments,
			Version:   _b[i].Arguments.Version,
		}
	}
	return rt
}

func updateStatus(ctx context.Context) {
	q := getQueues(ctx)
	if q == nil {
		clear(currentStatus)
		supportedRuntimes.Clear()
		return
	}
	if len(currentStatus) == 0 {
		b := getAllBindings(ctx)
		for _, item := range q.Items {
			a := item.Arguments
			currentStatus[item.Name] = status{
				Name:        a.Name,
				Version:     a.Version,
				Memory:      a.Memory,
				OS:          a.OS,
				Parallelism: a.Parallelism,
				BootedSince: a.BootedSince,
				Runtimes:    b[item.Name],
			}
		}
	} else {
		tm := make(map[string]status)
		for _, item := range q.Items {
			if v, ok := currentStatus[item.Name]; ok {
				tm[item.Name] = v
				continue
			}
			a := item.Arguments
			tm[item.Name] = status{
				Name:        a.Name,
				Version:     a.Version,
				Memory:      a.Memory,
				OS:          a.OS,
				Parallelism: a.Parallelism,
				BootedSince: a.BootedSince,
				Runtimes:    getBinding(ctx, item.Name),
			}
		}
		// strip unavailable judges and its runtimes
		for k, v := range currentStatus {
			if _, ok := tm[k]; !ok {
				delete(currentStatus, k)
				for _, r := range v.Runtimes {
					supportedRuntimes.Remove(r.ID)
				}
			}
		}
		// propagate supported runtimes list
		for k, v := range tm {
			currentStatus[k] = v
			for _, r := range v.Runtimes {
				supportedRuntimes.Add(r.ID)
			}
		}
	}
}

func (w *worker) updateStatus() func() {
	t.Do(func() {
		m.Lock()
		defer m.Unlock()
		updateStatus(w.ctx)
	})
	m.RLock()
	return m.RUnlock
}

func (w *worker) GetStatus() map[string]status {
	defer w.updateStatus()()
	return currentStatus
}

func (w *worker) IsRuntimeSupported(rt string) bool {
	defer w.updateStatus()()
	return supportedRuntimes.Contains(rt)
}
