package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/lvestera/slot-machine/internal/models"
)

const maxRetries = 3
const defaultDelay time.Duration = 1
const writeDBDelay time.Duration = 5
const dbConfig = "host=localhost user=gouser password=gouser dbname=gouser_db sslmode=disable"

const (
	tableCheckSQL  = "SELECT 1 FROM results;"
	tableCreateSQL = "CREATE TABLE IF NOT EXISTS results (" +
		"player int," +
		"spin bigint," +
		"result varchar(10)," +
		"win int" +
		");"

	insertSQL = "INSERT INTO results (player, spin, result, win) VALUES ($1, $2, $3, $4) "
)

type DBRepository struct {
	DB *sql.DB
}

func NewDBRepository() (*DBRepository, error) {
	db, err := sql.Open("pgx", dbConfig)
	if err != nil {
		return nil, err
	}

	//create repository
	rep := &DBRepository{DB: db}

	//check connection
	if err = rep.Ping(); err != nil {
		log.Print("Error db connection " + err.Error())
		return nil, err
	}

	log.Print("Db connection OK")
	log.Print("DB string " + dbConfig)

	//check table
	rows, tableCheck := db.Query(tableCheckSQL)
	if tableCheck != nil {
		_, err := db.ExecContext(context.Background(), tableCreateSQL)
		if err != nil {
			return nil, err
		}
	} else {
		err = rows.Err()
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
	}

	return rep, nil
}

func (rep *DBRepository) Add(m models.Result) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), writeDBDelay*time.Second)
	defer cancel()

	_, err := rep.DB.ExecContext(ctx, insertSQL, m.Player, m.Spin, m.Result, m.Win)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (rep *DBRepository) AddBatch(models map[int]models.Result) (bool, error) {
	tx, err := rep.DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(context.Background(), insertSQL)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	for _, m := range models {
		ctx, cancel := context.WithTimeout(context.Background(), writeDBDelay*time.Second)
		defer cancel()

		_, err := stmt.ExecContext(ctx, m.Player, m.Spin, m.Result, m.Win)
		if err != nil {
			return false, err
		}

	}
	return true, tx.Commit()

}

func (rep *DBRepository) Ping() error {
	return rep.DB.PingContext(context.Background())
}
