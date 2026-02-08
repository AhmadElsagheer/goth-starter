package repo

import (
	"github.com/ahmad/gother-example/pkg/phone"
	"github.com/ahmad/gother-example/server/modules/auth"

	"github.com/google/uuid"
)

type DbUser struct {
	ID    uuid.UUID   `db:"id"`
	Name  string      `db:"name"`
	Roles []auth.Role `db:"roles"`

	Email       string             `db:"email"`
	PhoneNumber *phone.PhoneNumber `db:"phoneNumber"`
}

func (u DbUser) toUser() auth.User {

	return auth.User{
		ID:          u.ID,
		Name:        u.Name,
		Roles:       u.Roles,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}
}

func fromUser(u auth.User) DbUser {
	out := DbUser{
		ID:          u.ID,
		Name:        u.Name,
		Roles:       u.Roles,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}

	return out
}
