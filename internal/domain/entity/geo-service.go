package entity

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
	House  string `json:"house"`
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
}

// SearchRequest структура для search-запросов
type SearchRequest struct {
	Query string `json:"query"`
}

// GeocodeRequest структура для geocode-запросов
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// ResponseBody Структура ответа
type ResponseBody struct {
	Addresses []*Address `json:"addresses"` // Слайс адресов
}
