package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"lib_isod_v2/file_service/internal/domain/models"
	"log/slog"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/multiprocessio/go-openoffice"
)

type Reader struct {
	log *slog.Logger
}

func NewReader(log *slog.Logger) *Reader {
	return &Reader{
		log: log,
	}
}

// TODO: Если headers нет. Только для recovery
// func (r *Reader) ReadCSVFileCreateWithHeaders(path string) ([]models.CreatePerson, error) {
// 	const op = "reader.ReadCSVFileCreateWithHeaders"

// 	rFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
// 	if err != nil {
// 		r.log.Error("error open file:", op, slog.Any("error", err.Error()))
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}
// 	defer rFile.Close()

// 	persons := []models.CreatePerson{}

// 	if err := gocsv.UnmarshalFile(rFile, &persons); err != nil {
// 		r.log.Error("error unmarshall file:", op, slog.Any("error", err.Error()))
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return persons, nil
// }

// @TODO: как-то обрабатывать, если не попадаем по количеству элементов, программа паникует
// person.LastName = row[0]
// person.FirstName = row[1]
// person.SurName = row[2]
// person.PlaceOfWork = row[3]
// person.Position = row[4]
// person.Login = row[5]
// person.TmpPasswd = row[6]
// person.DateOfBirthday = row[7]
// person.Complete = false
// persons = append(persons, person)

func (r *Reader) ReadCSVFileCreate(path string) ([]models.CreatePerson, error) {
	const op = "reader.ReadCSVFileCreate"

	personsWithOutHeaders := []models.CreatePerson{}

	rFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		r.log.Error("error open file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rFile.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(rFile)
		r.Comma = ','
		return r
	})

	if err := gocsv.UnmarshalWithoutHeaders(rFile, &personsWithOutHeaders); err != nil {
		r.log.Error("error unmarshall file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, если p.LastName == "last_name", значит у нас есть headers и мы их вырезаем
	for _, p := range personsWithOutHeaders {
		if p.LastName == "last_name" {
			return personsWithOutHeaders[1:], nil // [1:] потому что вырезаем из 0 элемента название колонок
		}
	}

	return personsWithOutHeaders, nil
}

// TODO: Проверять, если нет headers, то не нужно вырезать первый элемент
func (r *Reader) ReadODSFileCreate(path string) ([]models.CreatePerson, error) {
	const op = "reader.ReadODSFileCreate"

	f, err := openoffice.OpenODS(path)
	if err != nil {
		r.log.Error("error open file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	doc, err := f.ParseContent()
	if err != nil {
		r.log.Error("error parse file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	person := models.CreatePerson{}
	persons := []models.CreatePerson{}

	for _, t := range doc.Sheets {
		for _, row := range t.Strings() {
			person.LastName = row[0]
			person.FirstName = row[1]
			person.SurName = row[2]
			person.PlaceOfWork = row[3]
			person.Position = row[4]
			person.Login = row[5]
			person.TmpPasswd = row[6]
			person.DateOfBirthday = row[7]
			person.Complete = false
			persons = append(persons, person)
		}
	}

	// Если persons[0].LastName == "last_name" == true, значит headers есть и не нужно вырезать 0 элемент
	if persons[0].LastName == "last_name" {
		return persons[1:], nil // [1:] потому что вырезаем из 0 элемента название колонок
	}

	return persons, nil
}

func (r *Reader) ReadCSVFileRecovery(path string) ([]models.RecoveryPerson, error) {
	const op = "reader.ReadCSVFileRecovery"

	personsWithOutHeaders := []models.RecoveryPerson{}

	rFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		r.log.Error("error open file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rFile.Close()

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(rFile)
		r.Comma = ','
		return r
	})

	if err := gocsv.UnmarshalWithoutHeaders(rFile, &personsWithOutHeaders); err != nil {
		r.log.Error("error unmarshall file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, если p.LastName == "last_name", значит у нас есть headers и мы их вырезаем
	if personsWithOutHeaders[0].LastName == "last_name" {
		return personsWithOutHeaders[1:], nil // [1:] потому что вырезаем из 0 элемента название колонок
	}

	return personsWithOutHeaders, nil
}

func (r *Reader) ReadODSFileRecovery(path string) ([]models.RecoveryPerson, error) {
	const op = "reader.ReadODSFileCreate"

	f, err := openoffice.OpenODS(path)
	if err != nil {
		r.log.Error("error open file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	doc, err := f.ParseContent()
	if err != nil {
		r.log.Error("error parse file:", op, slog.Any("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	person := models.RecoveryPerson{}
	persons := []models.RecoveryPerson{}

	for _, t := range doc.Sheets {
		for _, row := range t.Strings() {
			person.LastName = row[0]
			person.FirstName = row[1]
			person.SurName = row[2]
			person.PlaceOfWork = row[3]
			person.Position = row[4]
			person.Login = row[5]
			person.TmpPasswd = row[6]
			person.Complete = false
			persons = append(persons, person)
		}
	}

	// Проверяем, если p.LastName == "last_name", значит у нас есть headers и мы их вырезаем
	if persons[0].LastName == "last_name" {
		return persons[1:], nil // [1:] потому что вырезаем из 0 элемента название колонок
	}

	return persons, nil
}
