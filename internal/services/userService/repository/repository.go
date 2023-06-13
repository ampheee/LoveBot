package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"strconv"
	"tacy/internal/services/userService"
	"tacy/pkg/botlogger"
)

type Repository struct {
	dbPool *pgxpool.Pool
	logger zerolog.Logger
}

func (r *Repository) InsertPhotoFromAdmin(ctx context.Context, photo []byte) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn")
		return err
	}
	defer conn.Release()
	q := `INSERT INTO images (image) VALUES ($1)`
	_, err = conn.Exec(ctx, q, photo)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to insert image into db")
		return err
	}
	return nil
}

func (r *Repository) InsertComplimentFromAdmin(ctx context.Context, compliment string) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire connection while insert compliment")
	}
	defer conn.Release()
	q := `INSERT INTO compliments(compliment) VALUES ($1)`
	_, err = conn.Exec(ctx, q, compliment)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to add new compliment into db")
	}
	return err
}

func (r *Repository) InsertThoughtsFromUser(ctx context.Context, thoughts string) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while insert thoughts")
	}
	defer conn.Release()
	q := `INSERT INTO thougths(thought) VALUES ($1)`
	_, err = conn.Exec(ctx, q, thoughts)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to insert thought")
	}
	return err
}

func (r *Repository) GetAllCompliments(ctx context.Context) ([]userService.Compliment, error) {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while getallcompliments")
		return nil, err
	}
	defer conn.Release()
	q := `SELECT * FROM compliments`
	rows, err := conn.Query(ctx, q)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to parse photo while send q")
		return nil, err
	}
	var compliments []userService.Compliment
	for rows.Next() {
		var compliment userService.Compliment
		err = rows.Scan(&compliment.Id, &compliment.Compliment)
		if err != nil {
			return nil, err
		}
		compliments = append(compliments, compliment)
	}
	return compliments, err
}

func (s *Repository) GetAllPhotos(ctx context.Context) ([]userService.Photo, error) {
	conn, err := s.dbPool.Acquire(ctx)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to acquire conn while getphoto")
		return nil, err
	}
	defer conn.Release()
	q := `SELECT (id, image) FROM images where seen = false`
	rows, err := conn.Query(ctx, q)
	if err != nil {
		s.logger.Warn().Err(err).Msg("unable to parse photo while send q")
		return nil, err
	}
	var photos []userService.Photo
	for rows.Next() {
		var photo userService.Photo
		err = rows.Scan(&photo.Id, &photo.Photo)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, err
}

func (r *Repository) GetPhotoByRandom(ctx context.Context) (userService.Photo, error) {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while getphoto")
		return userService.Photo{}, err
	}
	defer conn.Release()
	qFalseTotal := `select count(id) from images where seen = false`
	total, err := conn.Query(ctx, qFalseTotal)
	if err != nil {
		return userService.Photo{}, err
	}
	var max int
	err = total.Scan(&max)
	if err != nil {
	}
	if max == 0 {
		err = r.RefreshImages(ctx)
	}
	q := `SELECT id, image FROM images where seen = false ORDER BY RANDOM() LIMIT 1`
	row := conn.QueryRow(ctx, q)
	var photo userService.Photo
	err = row.Scan(&photo.Id, &photo.Photo)
	if err != nil {
	}
	err = r.UpdateImageById(ctx, photo.Id)
	return photo, nil
}

func (r *Repository) GetComplimentByRandom(ctx context.Context) (string, error) {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while getCompliment")
		return "", err
	}
	defer conn.Release()
	qFalseTotal := `select count(id) from compliments where seen = false`
	total, err := conn.Query(ctx, qFalseTotal)
	if err != nil {
		return "", err
	}
	var max int
	err = total.Scan(&max)
	if err != nil {
		return "", err
	}
	if max == 0 {
		err = r.RefreshCompliments(ctx)
	}
	q := `SELECT id, compliment FROM compliments where seen = false ORDER BY RANDOM() LIMIT 1`
	row := conn.QueryRow(ctx, q)
	var compliment userService.Compliment
	err = row.Scan(&compliment.Id, &compliment.Compliment)
	if err != nil {
		return "", err
	}
	err = r.UpdateComplimentById(ctx, compliment.Id)
	return compliment.Compliment, nil
}

func (r *Repository) UpdateComplimentById(ctx context.Context, id int) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while UpdateComplimentById")
		return err
	}
	defer conn.Release()
	q := `update compliments set seen = $1 where id = $2`
	_, err = conn.Exec(ctx, q, true, id)
	if err != nil {
		r.logger.Warn().Err(err)
		return err
	}
	return nil
}

func (r *Repository) UpdateImageById(ctx context.Context, id int) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while UpdateImageById")
		return err
	}
	defer conn.Release()
	q := `update images set seen = $1 where id = $2`
	_, err = conn.Exec(ctx, q, true, id)
	if err != nil {
		r.logger.Warn().Err(err)
		return err
	}
	return nil
}

func (r *Repository) RefreshImages(ctx context.Context) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while refreshImages")
		return err
	}
	defer conn.Release()
	q := `update images set seen = false where seen = true`
	_, err = conn.Exec(ctx, q)
	if err != nil {
		r.logger.Warn().Err(err)
		return err
	}
	return nil
}

func (r *Repository) RefreshCompliments(ctx context.Context) error {
	conn, err := r.dbPool.Acquire(ctx)
	if err != nil {
		r.logger.Warn().Err(err).Msg("unable to acquire conn while RefreshCompliments")
		return err
	}
	defer conn.Release()
	q := `update compliments set seen = false where seen = true`
	_, err = conn.Exec(ctx, q)
	if err != nil {
		r.logger.Warn().Err(err)
		return err
	}
	return nil
}

func (r *Repository) OutputAllCompliments(ctx context.Context) ([]string, error) {
	compliments, err := r.GetAllCompliments(ctx)
	if err != nil {
		r.logger.Warn().Err(err)
		return nil, err
	}
	var resCompliments []string
	for _, compliment := range compliments {
		resCompliments = append(resCompliments, strconv.Itoa(compliment.Id)+". "+compliment.Compliment)
	}
	return resCompliments, nil
}

func NewUSRepo(pool *pgxpool.Pool) userService.USRepository {
	return &Repository{
		dbPool: pool,
		logger: botlogger.GetLogger(),
	}
}
