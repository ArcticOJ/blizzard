package contests

import (
	"backend/blizzard/models"
	"backend/blizzard/pb"
	"io"
	"mime/multipart"
)

func CreateFile(file *multipart.FileHeader) *pb.File {
	f, e := file.Open()
	if e != nil {
		return nil
	}
	buf, e := io.ReadAll(f)
	if e != nil {
		return nil
	}
	return &pb.File{
		Buffer: buf,
	}
}

func Submit(ctx *models.Context) models.Response {
	// TODO: finalize response piping from judge to client
	//shouldStream := ctx.FormValue("streamed") == "true"
	code, e := ctx.FormFile("code")
	if e != nil {
		return ctx.Bad("Invalid submission.")
	}
	file := CreateFile(code)
	if ok, client := ctx.TrySelectClient(ctx.Request().Context()); ok {
		conn, e := client.Judge(ctx.Request().Context(), file)
		if e != nil {
			return ctx.InternalServerError("Could not establish connection to judging server.")
		}
		for _, err := conn.Recv(); true; {
			if err != nil {
				_ = conn.Close()
				break
			}
		} /*
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
			}*/
	} else {
		return ctx.InternalServerError("Could not find a suitable judge server.")
	}
	return nil
}
