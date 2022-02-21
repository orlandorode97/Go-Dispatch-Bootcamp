package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/orlandorode97/go-disptach/infrastructure/api"
	"github.com/orlandorode97/go-disptach/usercase/interactor"
)

type wordsController struct {
	wordsInteractor interactor.WordsInteractor
}

type WordsController interface {
	GetWords(w http.ResponseWriter, r *http.Request)
	GetWordsFromCSV(w http.ResponseWriter, r *http.Request)
	GetConcurrentWords(w http.ResponseWriter, r *http.Request)
}

func NewWordsController(r interactor.WordsInteractor) WordsController {
	return &wordsController{r}
}

// GetWords handles the requests and responses of the /words/ endpoint
func (wo *wordsController) GetWords(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	imageUrl, ok := params["url"]
	if !ok || imageUrl == "" {
		model.EncodeError(w, model.ErrInvalidData{Field: "url"})
		return
	}
	words, err := wo.wordsInteractor.Get(imageUrl)
	if err != nil {
		model.EncodeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(&words)
}

// GetWordsFromCSV handles the requests and responses of the /words/{id} endpoint
func (wo *wordsController) GetWordsFromCSV(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	words, err := wo.wordsInteractor.GetFromCSV(params["id"])
	if err != nil {
		model.EncodeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(&words)
}

// GetConcurrentWords handles the requests and responses of the /words-csv/ endpoint
func (wo *wordsController) GetConcurrentWords(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idType := params["type"]
	items := params["items"]
	itemsPerWorker := params["items_per_workers"]
	if strings.ToLower(idType) != api.Odd && strings.ToLower(idType) != api.Even {
		model.EncodeError(w, model.ErrInvalidData{Field: "type"})
		return
	}

	itemsResponse, err := valideRange(items, "items")
	if err != nil {
		model.EncodeError(w, err)
		return
	}
	perWorker, err := valideRange(itemsPerWorker, "items_per_workers")
	if err != nil {
		model.EncodeError(w, err)
		return
	}
	words, err := wo.wordsInteractor.GetConcurrent(idType, itemsResponse, perWorker)
	if err != nil {
		model.EncodeError(w, err)
		return
	}
	json.NewEncoder(w).Encode(&words)
}

func valideRange(value, name string) (int, error) {
	val, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if val < 0 {
		return 0, model.ErrInvalidData{Field: name}
	}
	return val, nil
}
