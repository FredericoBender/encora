package postgres

import (
	"encora/service"
	"log"

	sq "github.com/elgris/sqrl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// JobsExporter is a service that export all jobs extracted data to the database
type JobsExporter struct {
	DB         *sqlx.DB
	Logger     *log.Logger
	DebugLevel int8
}

func (s *JobsExporter) checkDependencies() error {
	if s.DB == nil {
		return errors.New("missing DB")
	}

	return nil
}

func (s *JobsExporter) Run(data *service.EncoraJobs) error {
	err := s.checkDependencies()
	if err != nil {
		return errors.Wrap(err, "dependencies:")
	}

	tx, err := s.DB.Beginx()
	if err != nil {
		return errors.Wrap(err, "db begin")
	}

	err = s.insertIntoJobs(data, tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return errors.Wrap(errors.Wrap(err, "insert into jobs"), errors.Wrap(rollbackErr, "tx rollback").Error())
		}
		return errors.Wrap(err, "insert into jobs")
	}
	return nil
}

func (s *JobsExporter) insertIntoJobs(jobs *service.EncoraJobs, tx *sqlx.Tx) error {
	for index, _ := range jobs.JobsTitle {
		query := psql.Insert("jobs").
			Columns(
				"title",
				"area",
				"country",
				"url",
				"description",
			)

		query = query.Values(
			jobs.JobsTitle[index],
			jobs.JobAreas[index],
			jobs.JobsCountries[index],
			jobs.JobsDetailsURLs[index],
			jobs.Description[index],
		)

		qSQL, args, err := query.ToSql()
		if err != nil {
			return errors.Wrap(err, "insert into jobs to SQL")
		}

		_, err = tx.Exec(qSQL, args...)
		if err != nil {
			return errors.Wrap(err, "insert into jobs tx exec")
		}
	}
	err := tx.Commit()
	if err != nil {
		return errors.Wrap(err, "tx commit jobs")
	}
	return nil
}
