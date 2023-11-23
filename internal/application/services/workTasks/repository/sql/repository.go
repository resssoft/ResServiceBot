package sql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"fun-coice/internal/application/services/workTasks/repository"
	"fun-coice/internal/application/services/workTasks/track"
	"github.com/doug-martin/goqu/v9"
)

const tableName = "timeTrack_tracks"

var _ repository.Repository = (*RepositorySQL)(nil)

type RepositorySQL struct {
	storage *sql.DB
	builder goqu.DialectWrapper
}

func NewSQLRepo(DB *sql.DB) (repository.Repository, error) {
	repo := &RepositorySQL{
		storage: DB,
		builder: goqu.Dialect("sqlite3"),
	}
	err := repo.Migrate()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RepositorySQL) NewItemLink() *track.Track {
	newItem := new(track.Track)
	return newItem
}

func (r *RepositorySQL) NewItemsLink() track.Tracks {
	return make(map[int64]track.Track)
}

func (r *RepositorySQL) Migrate() error {
	query := `CREATE TABLE IF NOT EXISTS ` + tableName + `
(
    id            INTEGER PRIMARY KEY,
    bot_name      TEXT     null,
    msg_id        INTEGER     null,
    track_id      TEXT     null,
    track_json    TEXT     null,
    status        INTEGER     null,
    user_id       INTEGER     null
);`
	statement, err := r.storage.Prepare(query) // Prepare SQL Statement
	if err != nil {
		return fmt.Errorf("migrate msg repo err prepare: %w ", err)
	}
	_, err = statement.Exec()
	if err != nil {
		return fmt.Errorf("migrate msg repo err: %w ", err)
	}
	return nil
}

func (r *RepositorySQL) Create(ctx context.Context, track track.Track) error {
	jsonData, err := json.Marshal(track)
	if err != nil {
		return err
	}
	jsonDataStr := string(jsonData)
	queryObject := r.builder.
		From(goqu.I(tableName)).
		Insert().
		Cols(goqu.C("bot_name"),
			goqu.C("msg_id"),
			goqu.C("track_id"),
			goqu.C("track_json"),
			goqu.C("user_id"),
			goqu.C("status")).
		Vals(goqu.Vals{
			&track.BotName,
			&track.MsgId,
			&track.Code,
			&jsonDataStr,
			&track.UserId,
			&track.Status,
		})

	query, args, err := queryObject.ToSQL()
	if err != nil {
		return err
	}
	//logger.RepositoryLogger(ctx).Info("Commissions", zap.String("query", query))
	_, err = r.storage.Exec(query, args...)
	return err
}

func (r *RepositorySQL) Update(ctx context.Context, item *track.Track) error {
	var err error
	return err
}

func (r *RepositorySQL) Get(ctx context.Context, ID int64) (*track.Track, error) {
	var err error
	return nil, err
}

func (r *RepositorySQL) List(ctx context.Context, filter track.TrackFilter) (track.Tracks, error) {
	queryObj := r.builder.From(tableName).Select("track_json")
	queryObj = r.Filter(ctx, queryObj, filter)
	query, _, err := queryObj.ToSQL()
	rows, err := r.storage.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := r.NewItemsLink()

	for rows.Next() {
		item := r.NewItemLink()
		jsonData := ""
		if err = rows.Scan(
			&jsonData,
		); err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(jsonData), item)
		if err != nil {
			return nil, err
		}
		items[item.UserId] = *item
	}
	return items, err
}

func (r *RepositorySQL) Delete(ctx context.Context, ID int64) error {
	var err error
	return err
}

func (r *RepositorySQL) Filter(ctx context.Context, obj *goqu.SelectDataset, filter track.TrackFilter) *goqu.SelectDataset {
	if filter.Code != nil {
		obj.Where(goqu.Ex{"code": filter.Code})
	}
	if filter.BotName != nil {
		obj.Where(goqu.Ex{"bot_name": filter.BotName})
	}
	if filter.UserId != nil {
		obj.Where(goqu.Ex{"user_id": filter.UserId})
	}
	if filter.MsgId != nil {
		obj.Where(goqu.Ex{"msg_id": filter.MsgId})
	}
	if filter.Status != 0 {
		obj.Where(goqu.Ex{"status": filter.Status})
	}
	return obj
}

//Condition{Field: "status", Value: *filter.Status}
