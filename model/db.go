package model

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.io/covid-19-api/cfg"
)

type DB struct {
	config *cfg.DbDef
}

func NewDB(dbConfig *cfg.DbDef) (*DB, error) {
	if dbConfig == nil {
		log.Error("DB config is empty")
		return nil, errors.New("DB config is empty")
	}
	return &DB{dbConfig}, nil
}
