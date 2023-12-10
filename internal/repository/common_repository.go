package repository

import (
	"github.com/sivaosorg/govm/dbx"

	dbresolver "github.com/sivaosorg/db.resolver"
)

type CommonRepository interface {
	GetPsqlStatus() dbx.Dbx
}

type commonRepositoryImpl struct {
	resolver *dbresolver.MultiTenantDBResolver
}

func NewCommonRepository(resolver *dbresolver.MultiTenantDBResolver) CommonRepository {
	return &commonRepositoryImpl{
		resolver: resolver,
	}
}

func (repo *commonRepositoryImpl) GetPsqlStatus() dbx.Dbx {
	_, s := repo.resolver.GetDefaultConnector()
	return s
}
