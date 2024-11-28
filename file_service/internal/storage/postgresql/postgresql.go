package postgresql

import (
	"context"
	baseErr "errors"
	"fmt"
	"lib_isod_v2/file_service/internal/domain/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate ${MOQPATH}moq -skip-ensure -pkg mocks -out ./mocks2/repo_mock.go . OrderRepo
type Storage struct {
	ctx context.Context
	db  *pgxpool.Pool
}

const (
	// tables
	createsTable  = "creates"
	recoveryTable = "recoverys"
	fileTable     = "files"
)

func New(ctx context.Context, storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := pgxpool.Connect(ctx, storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db:  db,
		ctx: ctx,
	}, nil
}

func (s *Storage) Stop() {
	s.db.Close()
}

func (s *Storage) SaveCreatePerson(persons []models.CreatePerson, file models.File) (uint64, error) {
	const op = "storage.postgresql.SaveCreatePerson"

	tx, err := s.db.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("can't create tx: %s", err.Error())
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(s.ctx)
			if rollbackErr != nil {
				err = baseErr.Join(err, rollbackErr)
			}
		}
	}()

	builder := sq.Insert(createsTable).Columns("last_name",
		"first_name",
		"sur_name",
		"place_of_work",
		"position",
		"login",
		"tmp_passwd",
		"date_of_birthday",
		"snils",
		"file_uuid",
		"complete")

	for _, person := range persons {
		builder = builder.Values(person.LastName,
			person.FirstName,
			person.SurName,
			person.PlaceOfWork,
			person.Position,
			person.Login,
			person.TmpPasswd,
			person.DateOfBirthday,
			person.Snils,
			file.Uuid,
			person.Complete)
	}

	query, args, err := builder.Suffix("RETURNING id").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: can't build sql:%w", op, err)
	}

	// query for createTable
	rows, err := tx.Query(s.ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: row err: %w", op, err)
	}
	defer rows.Close()

	var ID uint64
	for rows.Next() {
		if scanErr := rows.Scan(&ID); scanErr != nil {
			return 0, fmt.Errorf("%s can't scan id: %s", op, scanErr.Error())
		}
	}

	query, args, err = sq.Insert(fileTable).Columns("uuid",
		"file_path",
		"type",
		"extention",
		"created_at").Values(file.Uuid,
		file.FilePath,
		file.Type,
		file.Extension,
		time.Now().Format(time.RFC3339)).PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: can't build sql:%w", op, err)
	}

	if _, err := tx.Exec(s.ctx, query, args...); err != nil {
		return 0, fmt.Errorf("tx.Exec: %s", err.Error())
	}

	if err := tx.Commit(s.ctx); err != nil {
		return 0, fmt.Errorf("%s can't commit tx: %s", op, err.Error())
	}

	return ID, nil
}

func (s *Storage) SaveRecoveryPerson(persons []models.RecoveryPerson, file models.File) (uint64, error) {
	const op = "storage.postgresql.SaveRecoveryPerson"

	tx, err := s.db.BeginTx(s.ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("can't create tx: %s", err.Error())
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(s.ctx)
			if rollbackErr != nil {
				err = baseErr.Join(err, rollbackErr)
			}
		}
	}()

	builder := sq.Insert(recoveryTable).Columns("last_name",
		"first_name",
		"sur_name",
		"place_of_work",
		"position",
		"login",
		"tmp_passwd",
		"complete",
		"file_uuid",
	)

	for _, person := range persons {
		builder = builder.Values(person.LastName,
			person.FirstName,
			person.SurName,
			person.PlaceOfWork,
			person.Position,
			person.Login,
			person.TmpPasswd,
			person.Complete,
			file.Uuid)
	}

	query, args, err := builder.Suffix("RETURNING id").PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: can't build sql:%w", op, err)
	}

	// query for recoveryTable
	rows, err := tx.Query(s.ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s: row err: %w", op, err)
	}
	defer rows.Close()

	var ID uint64
	for rows.Next() {
		if scanErr := rows.Scan(&ID); scanErr != nil {
			return 0, fmt.Errorf("can't scan id: %s", scanErr.Error())
		}
	}

	query, args, err = sq.Insert(fileTable).Columns("uuid",
		"file_path",
		"type",
		"extention",
		"created_at").Values(file.Uuid,
		file.FilePath,
		file.Type,
		file.Extension,
		time.Now().Format(time.RFC3339)).PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return 0, fmt.Errorf("%s: can't build sql:%w", op, err)
	}

	if _, err := tx.Exec(s.ctx, query, args...); err != nil {
		return 0, fmt.Errorf("%s tx.Exec: %s", op, err.Error())
	}

	if err := tx.Commit(s.ctx); err != nil {
		return 0, fmt.Errorf("%s can't commit tx: %s", op, err.Error())
	}

	return ID, nil
}

func (s *Storage) GetCreatePersonByAdmin(limit uint64) ([]models.DTOCreatePersonByAdmin, error) {
	const op = "storage.postgresql.GetCreatePersonByAdmin"

	sql := `
		SELECT creates.last_name, creates.first_name, creates.sur_name, creates.place_of_work, creates.position, creates.login, creates.tmp_passwd, creates.date_of_birthday, creates.snils, creates.complete, files.created_at, files.extention, files.file_path, files.type, files.uuid from creates JOIN files ON creates.file_uuid=files.uuid ORDER BY created_at desc limit $1;
	`

	dtoCreatePerson := models.DTOCreatePersonByAdmin{}
	dtoCreatePersons := []models.DTOCreatePersonByAdmin{}

	rows, err := s.db.Query(s.ctx, sql, limit)
	if err != nil {
		return dtoCreatePersons, fmt.Errorf("%s can't query orders: %s", op, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(&dtoCreatePerson.LastName,
			&dtoCreatePerson.FirstName,
			&dtoCreatePerson.SurName,
			&dtoCreatePerson.PlaceOfWork,
			&dtoCreatePerson.Position,
			&dtoCreatePerson.Login,
			&dtoCreatePerson.TmpPasswd,
			&dtoCreatePerson.DateOfBirthday,
			&dtoCreatePerson.Snils,
			&dtoCreatePerson.Complete,
			&dtoCreatePerson.Created_at,
			&dtoCreatePerson.Extension,
			&dtoCreatePerson.FilePath,
			&dtoCreatePerson.Type,
			&dtoCreatePerson.File_Uuid,
		)

		if scanErr != nil {
			return dtoCreatePersons, fmt.Errorf("%s can't scan get create: %s", op, scanErr.Error())
		}

		dtoCreatePersons = append(dtoCreatePersons, dtoCreatePerson)
	}

	return dtoCreatePersons, nil
}

func (s *Storage) GetRecoveryPersonByAdmin(limit uint64) ([]models.DTORecoveryPersonByAdmin, error) {
	const op = "storage.postgresql.GetRecoveryPersonByAdmin"

	sql := `
		SELECT recoverys.last_name, recoverys.first_name, recoverys.sur_name, recoverys.place_of_work, recoverys.position, recoverys.login, recoverys.tmp_passwd, recoverys.complete, files.created_at, files.extention, files.file_path, files.type, files.uuid from recoverys JOIN files on recoverys.file_uuid=files.uuid order by id limit $1;
	`

	dtoRecoveryPerson := models.DTORecoveryPersonByAdmin{}
	dtoRecoveryPersons := []models.DTORecoveryPersonByAdmin{}

	rows, err := s.db.Query(s.ctx, sql, limit)
	if err != nil {
		return dtoRecoveryPersons, fmt.Errorf("%s can't query orders: %s", op, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(&dtoRecoveryPerson.LastName,
			&dtoRecoveryPerson.FirstName,
			&dtoRecoveryPerson.SurName,
			&dtoRecoveryPerson.PlaceOfWork,
			&dtoRecoveryPerson.Position,
			&dtoRecoveryPerson.Login,
			&dtoRecoveryPerson.TmpPasswd,
			&dtoRecoveryPerson.Complete,
			&dtoRecoveryPerson.Created_at,
			&dtoRecoveryPerson.Extension,
			&dtoRecoveryPerson.FilePath,
			&dtoRecoveryPerson.Type,
			&dtoRecoveryPerson.File_Uuid,
		)

		if scanErr != nil {
			return dtoRecoveryPersons, fmt.Errorf("%s can't scan get recovery: %s", op, scanErr.Error())
		}

		dtoRecoveryPersons = append(dtoRecoveryPersons, dtoRecoveryPerson)
	}

	return dtoRecoveryPersons, nil
}

func (s *Storage) SearchCreatePerson(value string, limit uint64) ([]models.DTOCreatePersonByAdmin, error) {
	const op = "storage.postgresql.SearchCreatePerson"

	sql := `
		SELECT creates.last_name, creates.first_name, creates.sur_name, creates.place_of_work, creates.position, creates.login, creates.tmp_passwd, creates.date_of_birthday, creates.snils, creates.complete, files.created_at, files.extention, files.file_path, files.type, files.uuid from creates JOIN files on creates.file_uuid=files.uuid where to_tsvector(creates.last_name || ' ' || creates.first_name || ' ' || creates.sur_name || ' ' || creates.place_of_work || ' ' || creates.position || ' ' || creates.login || ' ' || creates.tmp_passwd || ' ' || creates.date_of_birthday || ' ' || creates.snils || ' ' || creates.complete || ' ' || files.created_at || ' ' || files.extention || ' ' || files.file_path || ' ' || files.type || ' ' || files.uuid) @@ to_tsquery($1) limit $2;
	`

	dtoCreatePerson := models.DTOCreatePersonByAdmin{}
	dtoCreatePersons := []models.DTOCreatePersonByAdmin{}

	rows, err := s.db.Query(s.ctx, sql, value, limit)
	if err != nil {
		return dtoCreatePersons, fmt.Errorf("%s can't query orders: %s", op, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(&dtoCreatePerson.LastName,
			&dtoCreatePerson.FirstName,
			&dtoCreatePerson.SurName,
			&dtoCreatePerson.PlaceOfWork,
			&dtoCreatePerson.Position,
			&dtoCreatePerson.Login,
			&dtoCreatePerson.TmpPasswd,
			&dtoCreatePerson.DateOfBirthday,
			&dtoCreatePerson.Snils,
			&dtoCreatePerson.Complete,
			&dtoCreatePerson.Created_at,
			&dtoCreatePerson.Extension,
			&dtoCreatePerson.FilePath,
			&dtoCreatePerson.Type,
			&dtoCreatePerson.File_Uuid,
		)

		if scanErr != nil {
			return dtoCreatePersons, fmt.Errorf("%s can't scan search: %s", op, scanErr.Error())
		}

		dtoCreatePersons = append(dtoCreatePersons, dtoCreatePerson)
	}

	return dtoCreatePersons, nil
}

func (s *Storage) SearchRecoveryPerson(value string, limit uint64) ([]models.DTORecoveryPersonByAdmin, error) {
	const op = "storage.postgresql.SearchRecoveryPerson"

	sql := `
		SELECT recoverys.last_name, recoverys.first_name, recoverys.sur_name, recoverys.place_of_work, recoverys.position, recoverys.login, recoverys.tmp_passwd, recoverys.complete, files.created_at, files.extention, files.file_path, files.type, files.uuid from recoverys JOIN files on recoverys.file_uuid=files.uuid where to_tsvector(recoverys.last_name || ' ' || recoverys.first_name || ' ' || recoverys.sur_name || ' ' || recoverys.place_of_work || ' ' || recoverys.position || ' ' || recoverys.login || ' ' || recoverys.tmp_passwd || ' ' || recoverys.complete || ' ' || files.created_at || ' ' || files.extention || ' ' || files.file_path || ' ' || files.type || ' ' || files.uuid) @@ to_tsquery($1) limit $2;
	`

	dtoRecoveryPerson := models.DTORecoveryPersonByAdmin{}
	dtoRecoveryPersons := []models.DTORecoveryPersonByAdmin{}

	rows, err := s.db.Query(s.ctx, sql, value, limit)
	if err != nil {
		return dtoRecoveryPersons, fmt.Errorf("%s can't query search: %s", op, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		scanErr := rows.Scan(&dtoRecoveryPerson.LastName,
			&dtoRecoveryPerson.FirstName,
			&dtoRecoveryPerson.SurName,
			&dtoRecoveryPerson.PlaceOfWork,
			&dtoRecoveryPerson.Position,
			&dtoRecoveryPerson.Login,
			&dtoRecoveryPerson.TmpPasswd,
			&dtoRecoveryPerson.Complete,
			&dtoRecoveryPerson.Created_at,
			&dtoRecoveryPerson.Extension,
			&dtoRecoveryPerson.FilePath,
			&dtoRecoveryPerson.Type,
			&dtoRecoveryPerson.File_Uuid,
		)

		if scanErr != nil {
			return dtoRecoveryPersons, fmt.Errorf("%s can't scan order: %s", op, scanErr.Error())
		}

		dtoRecoveryPersons = append(dtoRecoveryPersons, dtoRecoveryPerson)
	}

	return dtoRecoveryPersons, nil
}

func (s *Storage) UpdateRecordComplete(id uint64, tableName string, complete bool) error {
	const op = "storage.postgresql.UpdateRecordComplete"

	query, args, err := sq.Update(tableName).Set("complete", complete).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%s can't build sql: %s", op, err.Error())
	}

	_, err = s.db.Exec(s.ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s cant exec update: %s", op, err.Error())
	}

	return nil
}
