package service

import (
	"github.com/sivaosorg/gocell/internal/repository"
	"github.com/sivaosorg/govm/dbx"
)

type CommonService interface {
	GetPsqlStatus() dbx.Dbx
}

type commonServiceImpl struct {
	commonRepository repository.CommonRepository
}

func NewCommonService(commonRepository repository.CommonRepository) CommonService {
	return &commonServiceImpl{
		commonRepository: commonRepository,
	}
}

func (s *commonServiceImpl) GetPsqlStatus() dbx.Dbx {
	return s.commonRepository.GetPsqlStatus()
}
