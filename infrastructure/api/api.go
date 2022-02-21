package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/orlandorode97/go-disptach/domain/model"
	"github.com/orlandorode97/go-disptach/workerpool"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

const (
	Odd  string = "odd"
	Even string = "even"
)

var (
	Client       HTTPClient
	WorkerPool   *workerpool.WorkerPool
	mutex        sync.Mutex
	itemsCounter int32 = 0
)

const (
	imaggaURL = "https://api.imagga.com/v2/text"
	csvPath   = "data/images_words.csv"
)

type Imagga struct {
	ApiKey    string
	SecretKey string
	ApiURL    string
	CSVPath   string
}

func init() {
	Client = &http.Client{}
}

// New returns a new instance of the Imagga client
func New(apiKey, apiSecret string) *Imagga {
	return &Imagga{
		ApiKey:    apiKey,
		SecretKey: apiSecret,
		ApiURL:    imaggaURL,
		CSVPath:   csvPath,
	}
}

func (i *Imagga) UpdateCSVPath(path string) {
	i.CSVPath = path
}

// GetTextWords reach out Imagga API by image url
func (i *Imagga) GetTextWords(imageUrl string) (*model.Result, error) {
	var result *model.Result
	url := fmt.Sprintf("%s?image_url=%s", i.ApiURL, imageUrl)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	request.SetBasicAuth(i.ApiKey, i.SecretKey)

	response, err := Client.Do(request)
	if err != nil {
		return nil, err
	}
	err = errorStatus(response.StatusCode)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if len(result.Text) == 0 {
		return nil, model.ErrNotFound{ImageUrl: imageUrl}
	}
	return result, nil
}

// GetTextWordById reads a local csv to find the word by id paramater
func (i *Imagga) GetTextWordById(id string) (*model.Result, error) {
	result := new(model.Result)

	words, err := i.Read(id)
	if err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, model.ErrNotFoundInCSV{Id: id}
	}

	result.Text = words

	return result, nil
}

// GetConcurrentWords reads concurrently the local csv file
func (i *Imagga) GetConcurrentWords(idType string, items, itemsWorker int) (*model.Result, error) {
	result := new(model.Result)
	workers := items / itemsWorker
	results := make(chan model.Text, items)
	if WorkerPool == nil {
		WorkerPool = workerpool.New()
	}
	WorkerPool.InitChan()
	WorkerPool.AddWorkers(workers)
	itemsCounter = 0
	wg := sync.WaitGroup{}
	file, err := i.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	csvReader := csv.NewReader(file)

	go func() {
		for word := range results {
			result.Text = append(result.Text, word)
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			WorkerPool.Add(worker(csvReader, idType, items, itemsWorker, results))
			wg.Done()
		}()
	}

	wg.Wait()

	WorkerPool.ShutDown()

	return result, nil
}
func worker(reader *csv.Reader, idType string, total, itemsWork int, results chan<- model.Text) func() {
	return func() {
		counter := 0
		for {
			if int(itemsCounter) == total {
				break
			}
			if counter == itemsWork {
				break
			}
			mutex.Lock()
			line, err := reader.Read()
			mutex.Unlock()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			idCsv, err := strconv.Atoi(line[0])
			if err != nil {
				break
			}
			if includeWord(idType, idCsv) {

				word, err := parseWord(line)
				if err != nil {
					break
				}
				results <- word
				counter++
				atomic.AddInt32(&itemsCounter, 1)
			}
		}
	}
}

// Open returns a pointer of the local csv file
func (i *Imagga) Open() (*os.File, error) {
	file, err := os.OpenFile(i.CSVPath, os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Read takes every definition record from the csv file into a []Text
func (i *Imagga) Read(id string) ([]model.Text, error) {
	words := make([]model.Text, 0)
	file, err := i.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if id == line[0] {
			word, err := parseWord(line)
			if err != nil {
				return nil, err
			}
			words = append(words, word)
			break
		}
	}

	return words, nil
}

// Write updates the local csv file with incoming words
func (i *Imagga) Write(result *model.Result) error {
	file, err := i.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	csvWriter := csv.NewWriter(file)
	for _, word := range result.Text {
		err = csvWriter.Write([]string{
			uuid.NewString(),
			word.Data,
			strconv.Itoa(int(word.Coordinates.Height)),
			strconv.Itoa(int(word.Coordinates.Width)),
		})
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()
	if csvWriter.Error() != nil {
		return csvWriter.Error()
	}
	return nil
}

// errorStatus returns the possible errors from the External Imagga API
func errorStatus(code int) error {
	switch code {
	case http.StatusForbidden:
		return model.ErrMissingApiKey{}
	case http.StatusBadRequest:
		return model.ErrInvalidData{Field: "imageUrl"}
	default:
		return nil
	}
}

// includeWord returns a bool if the word id is even or odd based on idType
func includeWord(idType string, wordId int) bool {
	if idType == Even {
		return wordId%2 == 0
	}
	return wordId%2 != 0
}

func parseWord(str []string) (model.Text, error) {

	height, err := strconv.Atoi(str[2])
	if err != nil {
		return model.Text{}, err
	}

	width, err := strconv.Atoi(str[2])
	if err != nil {
		return model.Text{}, err
	}

	return model.Text{
		ID:   str[0],
		Data: str[1],
		Coordinates: model.Coordinates{
			Height: int64(height),
			Width:  int64(width),
		},
	}, nil
}
