package interactor

import (
	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/orlandorode97/go-disptach/usercase/presenter"
	"github.com/orlandorode97/go-disptach/usercase/repository"
)

type WordsInteractor interface {
	Get(term string) (*model.Result, error)
	GetFromCSV(id string) (*model.Result, error)
	GetConcurrent(idType string, taskSize, perWorker int) (*model.Result, error)
}

type wordsInteractor struct {
	imaggaRepository repository.ImaggaRepository
	wordsPresenter   presenter.WordsPresenter
}

func New(repository repository.ImaggaRepository, presenter presenter.WordsPresenter) WordsInteractor {
	return &wordsInteractor{repository, presenter}
}

func (w *wordsInteractor) Get(term string) (*model.Result, error) {
	result, err := w.imaggaRepository.GetWordsByImageUrl(term)
	if err != nil {
		return nil, err
	}
	return w.wordsPresenter.ResponseWords(result)
}

func (w *wordsInteractor) GetFromCSV(id string) (*model.Result, error) {
	result, err := w.imaggaRepository.GetWordsById(id)
	if err != nil {
		return nil, err
	}
	return w.wordsPresenter.ResponseWords(result)
}

func (w *wordsInteractor) GetConcurrent(idType string, taskSize, perWorker int) (*model.Result, error) {
	result, err := w.imaggaRepository.GetConcurrentWords(idType, taskSize, perWorker)
	if err != nil {
		return nil, err
	}
	return w.wordsPresenter.ResponseWords(result)
}
