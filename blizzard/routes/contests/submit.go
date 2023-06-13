// TODO: complete submit route

package contests

import (
	"blizzard/blizzard/db/models/problems/submissions"
	"blizzard/blizzard/judge"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/pb/igloo"
	"blizzard/blizzard/utils"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
	"io"
	"mime/multipart"
)

type (
	Metadata struct {
		Cpu          uint32 `yaml:"cpu"`
		Memory       uint64 `yaml:"memory"`
		Duration     uint32 `yaml:"duration"`
		CaseCount    uint32 `yaml:"caseCount"`
		Input        string `yaml:"input"`
		Output       string `yaml:"output"`
		AllowPartial bool   `yaml:"allow_partial"`
	}
	Case struct {
		Time uint32 `yaml:"time"`
	}
	Result struct {
		Verdict  submissions.Verdict `json:"verdict"`
		Memory   uint64              `json:"memory"`
		Duration float64             `json:"duration"`
	}
)

func Parse() *Metadata {
	f := utils.ReadFile("/data/Dev/mock/uwu/metadata.yml")
	var metadata Metadata
	if e := yaml.Unmarshal(f, &metadata); e != nil {
		return nil
	}
	return &metadata
}

func prepare(id uint64, file *multipart.FileHeader) (*Metadata, *igloo.Submission) {
	f, e := file.Open()
	if e != nil {
		return nil, nil
	}
	buf, e := io.ReadAll(f)
	if e != nil {
		return nil, nil
	}
	m := Parse()
	return m, &igloo.Submission{
		Id:       id,
		Buffer:   buf,
		Language: "cpp11",
		Checker:  utils.ReadFile("/data/Dev/mock/test.py"),
		Metadata: &igloo.Metadata{
			Cpu:          m.Cpu,
			Memory:       m.Memory,
			Duration:     m.Duration,
			AllowPartial: m.AllowPartial,
			CaseCount:    m.CaseCount,
			Input:        m.Input,
			Output:       m.Output,
		},
	}
}

func Submit(ctx *extra.Context) models.Response {
	// TODO: finalize response piping from judge to client
	if ctx.RequireAuth() {
		return nil
	}
	code, e := ctx.FormFile("code")
	//problem := ctx.Param("id")
	if e != nil {
		return ctx.Bad("Invalid submission.")
	}
	_, file := prepare(1, code)
	if ok, client := judge.PickClient(context.Background()); ok {
		closed := false
		conn, e := client.Judge(judge.KeyContext(context.Background()), file)
		if e != nil {
			return ctx.InternalServerError("Could not establish connection to judging server.")
		}
		stream := ctx.StreamResponse()
		_res := make(chan interface{}, 3)
		go func() {
			<-ctx.Request().Context().Done()
			closed = true
			//close(_res)
		}()
		go func() {
			for r := range _res {
				if !closed {
					if _r, _ok := r.(*igloo.JudgeResult_Case); _ok {
						stream.Write(_r.Case)
					} else {
						stream.Write(r)
					}
				}
				// Handle response here.
				fmt.Println(r)
			}
		}()
		// TODO: remove this, create a long-running listener in the background and hook a listener and dispose on response close.
		go func() {
			for {
				res, err := conn.Recv()
				if err != nil {
					closed = true
					_ = conn.Close()
					break
				}
				if r, _ok := res.Result.(*igloo.JudgeResult_Case); _ok {
					_res <- r
				} else {
					r := res.GetFinal()
					_res <- echo.Map{
						"final":          true,
						"compilerOutput": r.CompilerOutput,
						"verdict":        r.Verdict,
					}
				}
			}
		}()
		for !closed {

		}
	} else {
		return ctx.InternalServerError("Could not find a suitable judge server.")
	}
	return nil
}
