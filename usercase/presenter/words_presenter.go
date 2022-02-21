package presenter

import "github.com/orlandorode97/go-disptach/domain/model"

type WordsPresenter interface {
	ResponseWords(result *model.Result) (*model.Result, error)
}
