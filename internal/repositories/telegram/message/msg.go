package tgmessage

import (
	"context"
	"database/sql"
	"fmt"
	tgModel "fun-coice/internal/domain/commands/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
)

var _ = (tgModel.MsgRepository)(&Repository{})

const (
	tableName = "tg_messages"
)

type Repository struct {
	storage *sql.DB
	builder goqu.DialectWrapper
}

func New(DB *sql.DB) (*Repository, error) {
	repo := &Repository{
		storage: DB,
		builder: goqu.Dialect("sqlite3"),
	}
	err := repo.Migrate()
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Repository) Migrate() error {
	query := `CREATE TABLE IF NOT EXISTS tg_messages
(
    id         INTEGER PRIMARY KEY,
    bot_name    TEXT     null,
    msg_id    INTEGER     null,
    msg_type    TEXT     null,
    msg_json    TEXT     null,
    msg_direction INTEGER   null,
    chat_id    INTEGER     null,
    chat_name    TEXT     null,
    chat_type    TEXT     null,
    from_id    INTEGER     null,
    from_username    TEXT     null,
    from_fname    TEXT     null,
    from_lname    TEXT     null,
    from_lang    TEXT     null,
    msg_text    TEXT     null,
    date    INTEGER     null
);
`
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

func (r *Repository) NewItemLink() *tgModel.Message {
	newItem := new(tgModel.Message)
	newItem.TgMsg = &tgbotapi.Message{
		Chat: &tgbotapi.Chat{},
		From: &tgbotapi.User{},
	}
	return newItem
}

func (r *Repository) NewItemsLink() tgModel.Messages {
	return make([]*tgModel.Message, 0)
}

func (r *Repository) Create(ctx context.Context, item *tgModel.Message) error {

	chatName := item.TgMsg.Chat.UserName
	if item.TgMsg.Chat.ID < 0 {
		chatName = item.TgMsg.Chat.Title
	}
	queryObject := r.builder.
		From(goqu.I(tableName)).
		Insert().
		Cols(goqu.C("bot_name"),
			goqu.C("msg_id"),
			goqu.C("msg_type"),
			goqu.C("msg_json"),
			goqu.C("msg_direction"),
			goqu.C("chat_id"),
			goqu.C("chat_name"),
			goqu.C("chat_type"),
			goqu.C("from_id"),
			goqu.C("from_username"),
			goqu.C("from_fname"),
			goqu.C("from_lname"),
			goqu.C("from_lang"),
			goqu.C("msg_text"),
			goqu.C("date")).
		Vals(goqu.Vals{
			&item.BotName,
			&item.TgMsg.MessageID,
			&item.MsgType,
			&item.MsgJson,
			&item.MsgDirection,
			&item.TgMsg.Chat.ID,
			&chatName,
			&item.TgMsg.Chat.Type,
			&item.TgMsg.From.ID,
			&item.TgMsg.From.UserName,
			&item.TgMsg.From.FirstName,
			&item.TgMsg.From.LastName,
			&item.TgMsg.From.LanguageCode,
			&item.TgMsg.Text,
			&item.TgMsg.Date,
		})

	query, args, err := queryObject.ToSQL()
	if err != nil {
		return err
	}
	//logger.RepositoryLogger(ctx).Info("Commissions", zap.String("query", query))
	_, err = r.storage.Exec(query, args...)
	return err
}

func (r *Repository) Update(ctx context.Context, item *tgModel.Message) error {
	var err error
	return err
}

func (r *Repository) Get(ctx context.Context, ID int64) (*tgModel.Message, error) {
	var err error
	return nil, err
}

func (r *Repository) List(ctx context.Context, filter *tgModel.MsgFilter) (tgModel.Messages, error) {
	items := r.NewItemsLink()

	query := "SELECT * FROM " + tableName
	rows, err := r.storage.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := r.NewItemLink()
		if err = rows.Scan(
			&item.ID,
			&item.BotName,
			&item.TgMsg.MessageID,
			&item.MsgType,
			&item.MsgJson,
			&item.MsgDirection,
			&item.TgMsg.Chat.ID,
			&item.TgMsg.Chat.UserName,
			&item.TgMsg.Chat.Type,
			&item.TgMsg.From.ID,
			&item.TgMsg.From.UserName,
			&item.TgMsg.From.FirstName,
			&item.TgMsg.From.LastName,
			&item.TgMsg.From.LanguageCode,
			&item.TgMsg.Text,
			&item.TgMsg.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, err
}

func (r *Repository) Delete(ctx context.Context, ID int64) error {
	var err error
	return err
}
