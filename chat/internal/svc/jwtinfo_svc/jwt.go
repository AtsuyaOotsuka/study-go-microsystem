package jwtinfo_svc

import "context"

type JwtStruct struct {
	UserID int
	Email  string
}

func NewJwtInfo(ctx context.Context) *JwtStruct {
	userID := ctx.Value(UserIDKey).(int)
	email := ctx.Value(EmailKey).(string)

	jwtinfo := &JwtStruct{
		UserID: userID,
		Email:  email,
	}

	return jwtinfo
}
