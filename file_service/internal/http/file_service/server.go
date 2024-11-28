package http

import (
	"lib_isod_v2/file_service/internal/domain/models"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Storage interface {
	SaveCreatePerson(persons []models.CreatePerson, file models.File) (uint64, error)
	SaveRecoveryPerson(persons []models.RecoveryPerson, file models.File) (uint64, error)
	GetCreatePersonByAdmin(limit uint64) ([]models.DTOCreatePersonByAdmin, error)
	GetRecoveryPersonByAdmin(limit uint64) ([]models.DTORecoveryPersonByAdmin, error)
	SearchCreatePerson(value string, limit uint64) ([]models.DTOCreatePersonByAdmin, error)
	SearchRecoveryPerson(value string, limit uint64) ([]models.DTORecoveryPersonByAdmin, error)
	UpdateRecordComplete(id uint64, tableName string, complete bool) error
}

type Routers struct {
	log     *slog.Logger
	storage Storage
}

func NewRouter(log *slog.Logger, storage Storage) *Routers {
	return &Routers{
		log:     log,
		storage: storage,
	}
}

func (r *Routers) PersonsByCreateByAdmin(c echo.Context) error {
	const op = "http.server.PersonsByCreateByAdmin"

	dto := new(LimitDTO)

	if err := c.Bind(dto); err != nil {
		r.log.Info("error binding to dto", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request person create" + err.Error()})
	}

	data, err := r.storage.GetCreatePersonByAdmin(dto.Limit)
	if err != nil {
		r.log.Info("error get persons by create", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"error get persons by create" + err.Error()})
	}

	return c.JSON(http.StatusOK, data)
}

func (r *Routers) PersonsByRecoveryByAdmin(c echo.Context) error {
	const op = "http.server.PersonsByRecoveryByAdmin"

	dto := new(LimitDTO)

	if err := c.Bind(dto); err != nil {
		r.log.Info("error binding to dto", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request person recovery" + err.Error()})
	}

	data, err := r.storage.GetRecoveryPersonByAdmin(dto.Limit)
	if err != nil {
		r.log.Info("error get persons by create", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"error get persons by recovery" + err.Error()})
	}

	return c.JSON(http.StatusOK, data)
}

func (r *Routers) SearchByCreate(c echo.Context) error {
	const op = "http.server.SearchByCreate"

	dto := new(SearchDTO)

	if err := c.Bind(dto); err != nil {
		r.log.Info("error binding search by create", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request search create" + err.Error()})
	}

	data, err := r.storage.SearchCreatePerson(dto.Value, dto.Limit)
	if err != nil {
		r.log.Info("error search by create", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"error search by create" + err.Error()})
	}

	return c.JSON(http.StatusOK, data)
}

func (r *Routers) SearchByRecovery(c echo.Context) error {
	const op = "http.server.SearchByRecovery"

	dto := new(SearchDTO)

	if err := c.Bind(dto); err != nil {
		r.log.Info("error binding search by recovery", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request search recovery" + err.Error()})
	}

	data, err := r.storage.SearchRecoveryPerson(dto.Value, dto.Limit)
	if err != nil {
		r.log.Info("error search by recovery", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"error search by recovery" + err.Error()})
	}

	return c.JSON(http.StatusOK, data)
}

func (r *Routers) UpdateField(c echo.Context) error {
	const op = "http.server.UpdateField"

	dto := new(UpdateDTO)
	if err := c.Bind(dto); err != nil {
		r.log.Info("error binding update", op, err)
		return c.JSON(http.StatusBadRequest, HTTPError{"bad request update" + err.Error()})
	}

	err := r.storage.UpdateRecordComplete(dto.Id, dto.TableName, dto.Complete)
	if err != nil {
		r.log.Info("error update field by id %s", op, slog.Any("id", dto.Id))
		return c.JSON(http.StatusBadRequest, HTTPError{"error update field" + err.Error()})
	}

	return c.JSON(http.StatusOK, HTTPSuccess{"ok"})
}
