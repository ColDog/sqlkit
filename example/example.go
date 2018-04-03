package example

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/coldog/sqlkit/db"
)

func uuid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type Base struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Base
	Email string
}

type UserRepo struct {
	db db.DB
}

func (rep *UserRepo) table() string { return "users" }

func (rep *UserRepo) Save(ctx context.Context, u *User) error {
	var stmt db.SQL
	if u.ID == "" {
		u.ID, u.CreatedAt, u.UpdatedAt = uuid(), time.Now(), time.Now()
		stmt = rep.db.Insert().Into(rep.table()).Record(u)
	} else {
		u.UpdatedAt = time.Now()
		stmt = rep.db.Update(rep.table()).Record(u)
	}
	return rep.db.Exec(ctx, stmt).Err()
}

func (rep *UserRepo) Get(ctx context.Context, id string) (*User, error) {
	u := &User{}
	err := rep.db.Query(
		ctx,
		rep.db.Select("*").From(rep.table()).Where("id = ?", id),
	).Decode(u)
	return u, err
}

func (rep *UserRepo) List(ctx context.Context) ([]*User, error) {
	u := []*User{}
	err := rep.db.Query(
		ctx,
		rep.db.Select("*").From(rep.table()),
	).Decode(&u)
	return u, err
}