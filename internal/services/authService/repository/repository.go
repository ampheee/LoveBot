package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
	"github.com/rs/zerolog"
	"tacy/config"
	"tacy/internal/services/authService"
	"tacy/pkg/botlogger"
)

type Repository struct {
	logger zerolog.Logger
	dbPool *pgxpool.Pool
	config config.Config
}

func (r *Repository) GetRoleById(ctx context.Context, id tg.UserID) (int, error) {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire connection")
	}
	defer conn.Release()
	qCheck := `SELECT EXISTS(SELECT userid FROM users WHERE userid = ($1))`
	var (
		found bool
		role  int
	)
	err = conn.QueryRow(ctx, qCheck, id).Scan(&found)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to found userid in db")
	}
	if found {
		qSelect := `SELECT role FROM users WHERE userid = ($1)`
		err = conn.QueryRow(ctx, qSelect, id).Scan(&role)
		if err != nil {
			r.logger.Warn().Err(err).Msg("unable to parse role from db")
		}
	} else {
		err = r.InsertNewUser(ctx, id, conn)
	}
	if err != nil {
		return 3, err
	}
	return role, nil
}

func (r *Repository) InsertNewUser(ctx context.Context, id tg.UserID, conn *pgxpool.Conn) error {
	var qInsert string
	if id.PeerID() == r.config.AcceptedUser {
		qInsert = `INSERT INTO users VALUES (($1), 2)`
	} else {
		qInsert = `INSERT INTO users VALUES (($1), 3)`
	}
	_, err := conn.Exec(ctx, qInsert, id)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to insert newuser into db")
	}
	return err
}

func NewAuthUsecase(pool *pgxpool.Pool) authService.Service {
	return &Repository{
		logger: botlogger.GetLogger(),
		dbPool: pool,
	}
}
