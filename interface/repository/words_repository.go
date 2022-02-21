package repository

import (
	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/orlandorode97/go-disptach/infrastructure/api"
)

type imaggaRepository struct {
	imaggaClient *api.Imagga
}

func New(imagga *api.Imagga) *imaggaRepository {
	return &imaggaRepository{imagga}
}

func (i *imaggaRepository) GetWordsByImageUrl(url string) (*model.Result, error) {
	result, err := i.imaggaClient.GetTextWords(url)
	if err != nil {
		return nil, err
	}
	err = i.imaggaClient.Write(result)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (i *imaggaRepository) GetWordsById(id string) (*model.Result, error) {
	result, err := i.imaggaClient.GetTextWordById(id)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (i *imaggaRepository) GetConcurrentWords(idType string, taskSize, workers int) (*model.Result, error) {
	result, err := i.GetConcurrentWords(idType, taskSize, workers)
	if err != nil {
		return nil, err
	}
	return result, nil
}
