package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mr-linch/go-tg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"tacy/internal/storage"
	"tacy/pkg/botlogger"
	"time"
)

type StorageS struct {
	dbPool *pgxpool.Pool
	logger zerolog.Logger
}

func (s *StorageS) GetCompliment(ctx context.Context) (string, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn while getCompliment")
		return "", err
	}
	defer conn.Release()
	q := `SELECT compliment FROM compliments`
	query, err := conn.Query(ctx, q)
	if err != nil {
		log.Warn().Err(err).Msg("something went wrong while get compliment")
	}
	var compliments []string
	for query.Next() {
		var compliment string
		err := query.Scan(&compliment)
		if err != nil {
			s.logger.Warn().Err(err).Msg("something wrong while scan compliments")
		}
		compliments = append(compliments, compliment)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return compliments[r.Intn(len(compliments))], nil
}

func (s *StorageS) GetPhoto(ctx context.Context) ([]byte, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn while getphoto")
		return nil, err
	}
	defer conn.Release()
	q := `SELECT image FROM images`
	rows, err := conn.Query(ctx, q)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to parse photo while send q")
		return nil, err
	}
	var photos [][]byte
	for rows.Next() {
		var photo []byte
		err := rows.Scan(&photo)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return photos[r.Intn(len(photos))], nil
}

func (s *StorageS) GetAllPhotos(ctx context.Context) ([][]byte, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn while getphoto")
		return nil, err
	}
	defer conn.Release()
	q := `SELECT image FROM images`
	rows, err := conn.Query(ctx, q)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to parse photo while send q")
		return nil, err
	}
	var photos [][]byte
	for rows.Next() {
		var photo []byte
		err := rows.Scan(&photo)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, err
}

func (s *StorageS) GetAllCompliments(ctx context.Context) ([]storage.Compliment, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn while getallcompliments")
		return nil, err
	}
	defer conn.Release()
	q := `SELECT * FROM compliments`
	rows, err := conn.Query(ctx, q)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to parse photo while send q")
		return nil, err
	}
	var compliments []storage.Compliment
	for rows.Next() {
		var compliment storage.Compliment
		err := rows.Scan(&compliment.Id, &compliment.Compliment)
		if err != nil {
			return nil, err
		}
		compliments = append(compliments, compliment)
	}
	return compliments, err
}

func (s *StorageS) InsertPhoto(ctx context.Context, buffer []byte) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn")
		return err
	}
	defer conn.Release()
	q := `INSERT INTO images (image) VALUES ($1)`
	_, err = conn.Exec(ctx, q, buffer)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to insert image into db")
		return err
	}
	return err
}

func (s *StorageS) InsertThoughts(ctx context.Context, thoughts string) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn while insert thoughts")
	}
	defer conn.Release()
	q := `INSERT INTO thougths(thought) VALUES ($1)`
	_, err = conn.Exec(ctx, q, thoughts)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to insert thought")
	}
	return err
}

func (s *StorageS) InsertCompliment(ctx context.Context, compliment string) error {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire connection while insert compliment")
	}
	defer conn.Release()
	q := `INSERT INTO compliments(compliment) VALUES ($1)`
	_, err = conn.Exec(ctx, q, compliment)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to add new compliment into db")
	}
	return err
}

func (s *StorageS) InsertNewUser(ctx context.Context, id tg.UserID, conn *pgxpool.Conn) error {
	qInsert := `INSERT INTO users VALUES (($1), 2)`
	_, err := conn.Exec(ctx, qInsert, id)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to insert newuser into db")
	}
	return err
}

func (s *StorageS) GetRoleById(ctx context.Context, id tg.UserID) (int, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire connection")
	}
	defer conn.Release()

	qCheck := `SELECT EXISTS(SELECT userid FROM users WHERE userid = ($1))`
	var found bool
	err = conn.QueryRow(ctx, qCheck, id).Scan(&found)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to found userid in db")
	}
	if found {
		var role int
		qSelect := `SELECT role FROM users WHERE userid = ($1)`
		err = conn.QueryRow(ctx, qSelect, id).Scan(&role)
		if err != nil {
			s.logger.Warn().Err(err).Msg("unable to parse role from db")
		}
		return role, err
	} else {
		err = s.InsertNewUser(ctx, id, conn)
	}
	return 2, err
}

func InitStorage(pool *pgxpool.Pool) storage.Storage {
	return &StorageS{
		dbPool: pool,
		logger: botlogger.GetLogger(),
	}
}
