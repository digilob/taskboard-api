package repository

import (
	"database/sql"

	"github.com/digilob/taskboard-api/pkg/api"
	_ "github.com/lib/pq"
)

type TaskRepository interface {
	CreateTask(task *api.Task) error
	GetTasks() ([]api.Task, error)
	UpdateTask(task *api.Task) error
	DeleteTask(id int) error
}

type PostgresTaskRepository struct {
	DB *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) TaskRepository {
	return &PostgresTaskRepository{DB: db}
}

// Implement the TaskRepository interface for PostgresTaskRepository

// CreateTask implementation
func (r *PostgresTaskRepository) CreateTask(task *api.Task) error {
	query := "INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id"
	err := r.DB.QueryRow(query, task.Title, task.Description, task.Status).Scan(&task.ID)
	return err
}

// GetTasks implementation
func (r *PostgresTaskRepository) GetTasks() ([]api.Task, error) {
	query := "SELECT id, title, description, status FROM tasks"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []api.Task
	for rows.Next() {
		var task api.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask implementation
func (r *PostgresTaskRepository) UpdateTask(task *api.Task) error {
	query := "UPDATE tasks SET title = $1, description = $2, status = $3 WHERE id = $4"
	result, err := r.DB.Exec(query, task.Title, task.Description, task.Status, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteTask implementation
func (r *PostgresTaskRepository) DeleteTask(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
