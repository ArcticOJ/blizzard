package contests

import (
	"blizzard/blizzard/judge"
	"blizzard/blizzard/models"
	"blizzard/blizzard/models/extra"
	"blizzard/blizzard/pb"
	"fmt"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

func CreateFile(id string, file *multipart.FileHeader) *pb.File {
	f, e := file.Open()
	if e != nil {
		return nil
	}
	buf, e := io.ReadAll(f)
	if e != nil {
		return nil
	}
	return &pb.File{
		Id:     id,
		Buffer: buf,
	}
}

func Submit(ctx *extra.Context) models.Response {
	// TODO: finalize response piping from judge to client
	shouldStream := ctx.FormValue("streamed") == "true"
	code, e := ctx.FormFile("code")
	problem := ctx.Param("id")
	if e != nil {
		return ctx.Bad("Invalid submission.")
	}
	id := fmt.Sprintf("%s_%s_%d_%s", strings.ReplaceAll(ctx.Get("user").(uuid.UUID).String(), "-", ""), problem, time.Now().Unix(), strconv.FormatUint(rand.Uint64(), 16))
	file := CreateFile(id, code)
	if ok, client := judge.TrySelectClient(ctx.Request().Context()); ok {
		conn, e := client.Judge(ctx.Request().Context(), file)
		if e != nil {
			return ctx.InternalServerError("Could not establish connection to judging server.")
		}
		var stream *models.ResponseStream
		if shouldStream {
			stream = ctx.StreamResponse()
		}
		for true {
			res, err := conn.Recv()
			if err != nil {
				_ = conn.Close()
				break
			}
			if shouldStream {
				_ = stream.Write(res)
			}
		}
	} else {
		return ctx.InternalServerError("Could not find a suitable judge server.")
	}
	return nil
}
