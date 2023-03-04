package user

import (
	"backend/blizzard/models"
)

func ApiKey(ctx *models.Context) models.Response {
	switch ctx.Method() {
	case models.Get:

	}
	return nil
}
