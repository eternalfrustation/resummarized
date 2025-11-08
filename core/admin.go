package core

import "context"

type Admin struct {
	AdminID    int32  `db:"admin_id"`
	AdminEmail string `db:"admin_email"`
}

func (app *App) IsAdmin(ctx context.Context, uid string) (bool, error) {
	user, err := authClient.GetUser(ctx, uid)
	if err != nil {
		return false, err
	}
	admin, err := FetchRow[Admin](app, ctx, "sql/find_admin_by_email.sql", user.Email)
	if err != nil {
		return false, err
	}
	if admin != nil {
		return true, nil
	}
	return false, nil
}
