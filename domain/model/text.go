package model

// Struct to unmarshal the Urban dictionary response of definitions based on a term
// API For more infomation: https://rapidapi.com/community/api/urban-dictionary/
type ImaggaResponse struct {
	Result Result `json:"result"`
	Status Status `json:"status"`
}

type Result struct {
	Text []Text `json:"text"`
}

type Text struct {
	ID          string
	Coordinates Coordinates `json:"coordinates"`
	Data        string      `json:"data"`
}

type Coordinates struct {
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
	Xmax   int64 `json:"xmax"`
	Xmin   int64 `json:"xmin"`
	Ymax   int64 `json:"ymax"`
	Ymin   int64 `json:"ymin"`
}

type Status struct {
	Text string `json:"text"`
	Type string `json:"type"`
}
