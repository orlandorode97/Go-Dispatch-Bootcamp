package presenter

import (
	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/orlandorode97/go-disptach/usercase/presenter"
)

var (
	UrbanLayout = "2006-01-02T15:04:05.999Z"
	UserLayout  = "Mon, 02-January-2006 15:04"
)

type wordsPresenter struct{}

func New() presenter.WordsPresenter {
	return &wordsPresenter{}
}

// ResponseWords return the list of definitions fulfilling the WordsPresenter interface
func (w *wordsPresenter) ResponseWords(result *model.Result) (*model.Result, error) {
	return result, nil
}
