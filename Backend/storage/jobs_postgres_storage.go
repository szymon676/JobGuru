package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/szymon676/jobguru/types"
)

type PostgreJobStorage struct {
	db *sql.DB
}

func NewPostgreJobStorage(db *sql.DB) *PostgreJobStorage {
	return &PostgreJobStorage{
		db: db,
	}
}

func (js PostgreJobStorage) CreateJob(req types.JobReq) error {
	query := "INSERT INTO jobs (user_id, title, company, skills, salary, description, currency, date, location) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	convskills := pq.Array(req.Skills)

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}

	_, err = js.db.Exec(query, req.UserID, req.Title, req.Company, convskills, req.Salary, req.Description, req.Currency, date, req.Location)
	if err != nil {
		return fmt.Errorf("insert into jobs failed: %v", err)
	}

	return nil
}

func (js PostgreJobStorage) GetJobs() ([]types.Job, error) {
	query := "SELECT * FROM jobs;"

	rows, err := js.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []types.Job
	for rows.Next() {
		job, err := scanJob(rows)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func (js PostgreJobStorage) GetJobsByUserID(userid int) ([]types.Job, error) {
	query := "SELECT * FROM jobs WHERE user_id = $1;"

	rows, err := js.db.Query(query, userid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []types.Job
	for rows.Next() {
		job, err := scanJob(rows)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, *job)
	}

	return jobs, nil
}

func (js PostgreJobStorage) UpdateJobByID(ID int, req types.JobReq) error {
	var count int

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return fmt.Errorf("failed to parse date: %v", err)
	}

	if err := js.db.QueryRow("SELECT COUNT(*) FROM jobs WHERE id = $1", ID).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("job with id %v does not exist", ID)
	}

	convskills := pq.Array(req.Skills)

	query := "UPDATE jobs SET user_id = $1, title = $2, company = $3, skills = $4, salary = $5, description = $6, currency = $7, date = $8, location = $9 WHERE id = $10;"
	_, err = js.db.Exec(query, req.UserID, req.Title, req.Company, convskills, req.Salary, req.Description, req.Currency, date, req.Location, ID)
	if err != nil {
		return err
	}

	return nil
}

func (js PostgreJobStorage) DeleteJobByID(ID int) error {
	query := "DELETE FROM jobs WHERE id = $1"

	_, err := js.db.Exec(query, ID)

	if err != nil {
		return err
	}
	return nil
}

func scanJob(rows *sql.Rows) (*types.Job, error) {
	job := new(types.Job)
	err := rows.Scan(
		&job.ID,
		&job.UserID,
		&job.Title,
		&job.Company,
		pq.Array(&job.Skills),
		&job.Salary,
		&job.Description,
		&job.Currency,
		&job.Date,
		&job.Location,
	)

	return job, err
}
