package model

type Links struct {
	Self string `json:"self"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

type Page struct {
	Items      []Game `json:"items"`
	Page       int64  `json:"page"`
	PageSize   int64  `json:"pageSize"`
	TotalItems int64  `json:"totalItems"`
	TotalPages int64  `json:"totalPages"`
	Links      Links  `json:"_links"`
}
