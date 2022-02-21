package repository

import (
	"github.com/orlandorode97/go-disptach/domain/model"
)

type ImaggaRepository interface {
	GetWordsByImageUrl(term string) (*model.Result, error)
	GetWordsById(id string) (*model.Result, error)
	GetConcurrentWords(idType string, taskSize, perWorker int) (*model.Result, error)
}
