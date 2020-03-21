package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"short/utils"
)

type ShortDB struct {
	db *sql.DB
}

func (s *ShortDB) Open() error {
	db, err := sql.Open("mysql", utils.Conf.DB.DSN)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	db.SetMaxIdleConns(utils.Conf.DB.MaxIdleConn)
	db.SetMaxOpenConns(utils.Conf.DB.MaxOpenConn)
	s.db = db
	return nil
}

func (s *ShortDB) Close() {
	if s.db != nil {
		s.db.Close()
		s.db = nil
	}
}

func (s *ShortDB) Store(longURL, shortURL string) error {
	// TODO: improve duplication check and storage
	// since i use long_url as unique key, it can NOT be too long (MYSQL: max key length is 3072 bytes)
	// maybe its better to store short_url/long_url mapping relationship into redis set (to remove duplicate)
	// and use short_url as unique key in MYSQL for storage?
	stmt, err := s.db.Prepare(fmt.Sprintf(`REPLACE INTO short(long_url, short_url) VALUES(?, ?)`))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(longURL, shortURL)
	return err
}

func (s *ShortDB) Query(shortURL string) (longURL string, err error) {
	rows, err := s.db.Query(fmt.Sprintf(`SELECT long_url FROM short WHERE short_url=?`), shortURL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&longURL)
		if err != nil {
			return
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

// create database short
// create user 'short'@'%' identified by 'short'
// grant select, insert, delete on short.short to 'short'@'%'
