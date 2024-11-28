package watcher

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"lib_isod_v2/file_service/internal/domain/models"

	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	log     *slog.Logger

	create   string
	recovery string
}

const (
	ods = "ods"
	csv = "csv"
)

// TODO: add hash files
func NewWatcher(log *slog.Logger, createPath string, recoveryPath string) (*Watcher, error) {
	const op = "watcher.NewWatcher"

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Watcher{
		watcher:  watcher,
		log:      log,
		create:   createPath,
		recovery: recoveryPath,
	}, nil
}

// TODO: Добавить workerpool, в который будут складываться файлы, которые пришли
// И из workerpool'a буду разбираться в одну таблицу для файлов и в сервис для разбора файла
// Или pipeline
func (w *Watcher) Run() (<-chan models.File, <-chan models.File) {
	const op = "watcher.Run"

	filesCreate := make(chan models.File)
	filesRecovery := make(chan models.File)

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				switch event.Op {
				case fsnotify.Remove, fsnotify.Rename, fsnotify.Chmod:
					w.watcher.Remove(event.Name)
					if err := w.watcher.Add(event.Name); err != nil {
						w.log.Info("error while resetting watch on file", op, slog.Any("error", err.Error()))
						continue
					}
				case fsnotify.Create:
					w.log.Info("saving file", op, slog.String("filename:", event.Name))
					time.Sleep(500 * time.Millisecond) // ждем, пока файл запишется
					// recovery
					if strings.Contains(event.Name, w.recovery) {
						if strings.Contains(event.Name, ods) {
							filesRecovery <- models.File{
								Uuid:       uuid.New(),
								FilePath:   event.Name,
								Type:       w.recovery,
								Extension:  ods,
								Created_at: time.Now(),
							}

							continue
						}

						filesRecovery <- models.File{
							Uuid:       uuid.New(),
							FilePath:   event.Name,
							Type:       w.recovery,
							Extension:  csv,
							Created_at: time.Now(),
						}
					} else { // create
						if strings.Contains(event.Name, csv) {
							filesCreate <- models.File{
								Uuid:       uuid.New(),
								FilePath:   event.Name,
								Type:       w.create,
								Extension:  csv,
								Created_at: time.Now(),
							}

							continue
						}

						filesCreate <- models.File{
							Uuid:       uuid.New(),
							FilePath:   event.Name,
							Type:       w.create,
							Extension:  ods,
							Created_at: time.Now(),
						}
					}

				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				w.log.Error("%s error:", op, slog.Any("error", err.Error()))
			}
		}
	}()

	err := w.watcher.Add("./tmp/recovery")
	if err != nil {
		w.log.Error("%s error recovery", op, slog.Any("error", err.Error()))
	}

	err = w.watcher.Add("./tmp/create")
	if err != nil {
		w.log.Error("%s error create", op, slog.Any("error", err.Error()))
	}

	return filesCreate, filesRecovery
}

func (w *Watcher) Close() error {
	const op = "watcher.Stop"

	w.log.Info("stopping", op, "watcher")

	return w.watcher.Close()
}



