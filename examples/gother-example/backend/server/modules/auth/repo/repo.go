package repo

import (
	"context"

	"github.com/ahmad/gother-example/pkg/pg"
	"github.com/ahmad/gother-example/pkg/xerrors"

	"github.com/ahmad/gother-example/server/modules/auth"
	sqtypes "github.com/ahmad/gother-example/server/schema/postgres/types"

	"github.com/google/uuid"
)

//go:generate mockgen -source ./repo.go -destination ./mockrepo/mockrepo.go -package mockrepo
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (auth.User, error)
	GetByEmail(ctx context.Context, email string) (auth.User, error)
	GetByPhoneNumber(ctx context.Context, phoneNumber string) (auth.User, error)
}

type repo struct {
	userDao pg.Dao[sqtypes.USERS, DbUser]
}

func New(db pg.Database) Repository {
	return &repo{
		userDao: pg.NewDao[sqtypes.USERS, DbUser](db),
	}
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (auth.User, error) {
	table := r.userDao.Table()
	dbUser, err := r.userDao.Get(ctx, pg.Where(table.ID.Eq(id)))
	if err != nil {
		return auth.User{}, convertError(err)
	}
	return dbUser.toUser(), nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (auth.User, error) {
	table := r.userDao.Table()
	dbUser, err := r.userDao.Get(ctx, pg.Where(table.EMAIL.EqString(email)))
	if err != nil {
		return auth.User{}, convertError(err)
	}
	return dbUser.toUser(), nil
}

func (r *repo) GetByPhoneNumber(ctx context.Context, phoneNumber string) (auth.User, error) {
	table := r.userDao.Table()
	dbUser, err := r.userDao.Get(ctx, pg.Where(table.PHONENUMBER.EqString(phoneNumber)))
	if err != nil {
		return auth.User{}, convertError(err)
	}
	return dbUser.toUser(), nil
}

func convertError(err error) error {
	if err == pg.ErrNoRows {
		return xerrors.NewNotFound("no user found")
	}
	return err
}
