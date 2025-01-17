package http

import "github.com/labstack/echo/v4"

type LimitDTO struct {
	Limit uint64 `json:"limit" from:"limit" query:"limit"`
}

type SearchDTO struct {
	LimitDTO
	Value string `json:"value" from:"value" query:"value"`
}

type UpdateDTO struct {
	Id        uint64 `json:"id" from:"id" query:"id"`
	TableName string `json:"table_name" from:"table_name" query:"table_name"`
	Complete  bool   `json:"complete" from:"complete" query:"compelte"`
}

type MediaDto struct {
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"message"`
	Data       *echo.Map `json:"data"`
}
