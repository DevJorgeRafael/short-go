package context

import "context"

func GetUserID(ctx context.Context) string{
	userID, _ := ctx.Value(UserIdKey).(string)
	return userID
}