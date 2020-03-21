package short

import (
	"log"
	"short/base"
	"short/db"
	"short/utils"
	"time"
)

type Shorter struct {
	db      *db.ShortDB
	sf      *utils.SnowFlake
	baseStr string
}

func (s *Shorter) Short(longURL string) (shortURL string, err error) {
	id, err := s.sf.NextUID()
	if err != nil {
		return
	}
	shortURL = base.Uint64ToString(id, s.baseStr)
	err = s.db.Store(longURL, shortURL)
	return
}

func (s *Shorter) Long(shortURL string) (longURL string, err error) {
	return s.db.Query(shortURL)
}

func (s *Shorter) Close() {
	s.db.Close()
}

var DefaultShorter *Shorter

func StartShorter() {
	sdb := db.ShortDB{}
	if err := sdb.Open(); err != nil {
		log.Panic(err)
	}
	sf, err := utils.NewSnowflake(utils.SnowFlakeOptions{
		StartTime: time.Time{},
		MachineID: func() uint16 {
			return 1
		},
	})
	if err != nil {
		log.Panic(err)
	}
	DefaultShorter = &Shorter{
		db:      &sdb,
		sf:      sf,
		baseStr: utils.Conf.BaseString,
	}
}
