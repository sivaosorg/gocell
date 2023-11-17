package repository

import (
	"github.com/sivaosorg/govm/dbx"
	"github.com/sivaosorg/postgresconn/postgresconn"
)

type CommonRepository interface {
	GetPsqlStatus() dbx.Dbx
}

type commonRepositoryImpl struct {
	psql       *postgresconn.Postgres
	psqlStatus dbx.Dbx
}

func NewCommonRepository(psql *postgresconn.Postgres, psqlStatus dbx.Dbx) CommonRepository {
	return &commonRepositoryImpl{
		psql:       psql,
		psqlStatus: psqlStatus,
	}
}

func (repo *commonRepositoryImpl) GetPsqlStatus() dbx.Dbx {
	return repo.psqlStatus
}
