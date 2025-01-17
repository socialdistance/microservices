package file

import (
	"lib_isod_v2/file_service/internal/domain/models"
	"log/slog"
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

type Watcher interface {
	Run() (<-chan models.File, <-chan models.File)
	Close() error
}

type Reader interface {
	ReadCSVFileCreate(path string) ([]models.CreatePerson, error)
	ReadODSFileCreate(path string) ([]models.CreatePerson, error)
	ReadCSVFileRecovery(path string) ([]models.RecoveryPerson, error)
	ReadODSFileRecovery(path string) ([]models.RecoveryPerson, error)
}

type File struct {
	log     *slog.Logger
	storage Storage
	watcher Watcher
	reader  Reader
	done    chan struct{}
}

func New(log *slog.Logger, storage Storage, watcher Watcher, reader Reader) *File {

	return &File{
		log:     log,
		storage: storage,
		watcher: watcher,
		reader:  reader,
		done:    make(chan struct{}),
	}
}

func (f *File) FileRun() {
	const op = "file.FileRun"

	f.log.Info(op, slog.String("Start", "watcher"))
	changesCreate, changesRecovery := f.watcher.Run()

	for {
		select {
		case <-f.done:
			f.log.Info("Stopping", op, "file run")
			return
		case fileCreate, ok := <-changesCreate:
			if !ok {
				continue
			}

			if fileCreate.Extension == "csv" {
				personsCreateCSV, err := f.reader.ReadCSVFileCreate(fileCreate.FilePath)
				if err != nil {
					f.log.Error(op, slog.Any("error read csv file create", err.Error()))
				}

				id, err := f.storage.SaveCreatePerson(personsCreateCSV, fileCreate)
				if err != nil {
					f.log.Error(op, slog.Any("error save person csv create", err.Error()))
				}

				f.log.Info("Save persons from .csv file by id in create", slog.Uint64("id:", id))
			}

			if fileCreate.Extension == "ods" {
				personsCreatOds, err := f.reader.ReadODSFileCreate(fileCreate.FilePath)
				if err != nil {
					f.log.Error(op, slog.Any("error read ods file create", err.Error()))
				}

				id, err := f.storage.SaveCreatePerson(personsCreatOds, fileCreate)
				if err != nil {
					f.log.Error(op, slog.Any("error save person ods create", err.Error()))
				}

				f.log.Info("Save persons from .ods file by id in create", slog.Uint64("id:", id))
			}

		case fileRecovery, ok := <-changesRecovery:
			if !ok {
				continue
			}

			if fileRecovery.Extension == "csv" {
				personsRecoveryCSV, err := f.reader.ReadCSVFileRecovery(fileRecovery.FilePath)
				if err != nil {
					f.log.Error(op, slog.Any("error read csv file recovery", err.Error()))
				}

				id, err := f.storage.SaveRecoveryPerson(personsRecoveryCSV, fileRecovery)
				if err != nil {
					f.log.Error(op, slog.Any("error save person csv recovery", err.Error()))
				}

				f.log.Info("Save persons from .csv file by id in recovery", slog.Uint64("id:", id))
			}

			if fileRecovery.Extension == "ods" {
				personsRecoveryOds, err := f.reader.ReadODSFileRecovery(fileRecovery.FilePath)
				if err != nil {
					f.log.Error(op, slog.Any("error read ods file recovery", err.Error()))
				}

				id, err := f.storage.SaveRecoveryPerson(personsRecoveryOds, fileRecovery)
				if err != nil {
					f.log.Error(op, slog.Any("error save person ods recovery", err.Error()))
				}

				f.log.Info("Save persons from .ods file by id in recovery", slog.Uint64("id:", id))
			}
		}
	}
}

func (f *File) Stop() {
	const op = "file.Stop"
	f.log.Info("stoping", op, "fileRun")

	close(f.done)
}
