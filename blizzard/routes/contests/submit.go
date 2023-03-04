package contests

import (
	"backend/blizzard/db/models/problems/submissions"
	"backend/blizzard/models"
	"fmt"
	"io"
	"time"
)

func Submit(ctx *models.Context) models.Response {
	shouldStream := ctx.FormValue("streamed") == "true"
	shouldStream = true
	/*_, e := ctx.FormFile("code")
	if e != nil {
		return ctx.Forbid()
	}*/

	conn, e := ctx.Server.Polar["north-pole"].Judge(ctx.Request().Context(), nil)
	if e != nil {
		return ctx.InternalServerError("Could not establish connection to judging server.")
	}
	for _, err := conn.Recv(); true; {
		if err != nil {
			_ = conn.Close()
			if err == io.EOF {
				fmt.Println("End")
			} else {
				fmt.Println(err)
			}
		}

	}
	if shouldStream {
		stream := ctx.StreamResponse()
		print(ctx.Param("id"))
		for i := 0; i < 10; i++ {
			r := submissions.TestResult{
				Memory:   10,
				Duration: 25,
				Verdict:  submissions.Accepted,
				Point:    10,
			}
			if i%3 == 0 {
				r = submissions.TestResult{
					Memory:   10,
					Duration: 25,
					Verdict:  submissions.WrongAnswer,
					Point:    0,
				}
			} else if i%2 == 0 {
				r = submissions.TestResult{
					Memory:   10,
					Duration: 25,
					Verdict:  submissions.WrongAnswer,
					Point:    0,
				}
			}
			if stream.Write(r) != nil {
				return nil
			}
			time.Sleep(time.Second * 1)
		}
	}
	return nil
}
