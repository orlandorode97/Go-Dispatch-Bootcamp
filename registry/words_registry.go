package registry

import (
	"github.com/orlandorode97/go-disptach/interface/controller"
	p "github.com/orlandorode97/go-disptach/interface/presenter"
	re "github.com/orlandorode97/go-disptach/interface/repository"
	"github.com/orlandorode97/go-disptach/usercase/interactor"
	"github.com/orlandorode97/go-disptach/usercase/presenter"
	"github.com/orlandorode97/go-disptach/usercase/repository"
)

func (r *registry) NewWordsController() controller.WordsController {
	return controller.NewWordsController(r.NewWordsInteractor())
}

func (r *registry) NewWordsInteractor() interactor.WordsInteractor {
	return interactor.New(r.NewWordsRepository(), r.NewWordsPresenter())
}

func (r *registry) NewWordsRepository() repository.ImaggaRepository {
	return re.New(r.ImaggaClient)
}

func (r *registry) NewWordsPresenter() presenter.WordsPresenter {
	return p.New()
}
