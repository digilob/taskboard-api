package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Task represents a single task.
type Task struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type Repository interface {
	GetAllTasks() ([]Task, error)
	CreateTask(task Task) (Task, error)
	UpdateTask(id string, task Task) (Task, error)
	DeleteTask(id string) error
}

type PostgresTaskRepository struct {
	DB *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) Repository {
	return &PostgresTaskRepository{DB: db}
}

// Implement the Repository interface for PostgresTaskRepository
// CreateTask implementation
func (r *PostgresTaskRepository) CreateTask(task Task) (Task, error) {
	query := "INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id"
	err := r.DB.QueryRow(query, task.Title, task.Description).Scan(&task.ID)
	return task, err
}

// GetAllTasks implementation
func (r *PostgresTaskRepository) GetAllTasks() ([]Task, error) {
	query := "SELECT id, title, description FROM tasks"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// UpdateTask implementation
func (r *PostgresTaskRepository) UpdateTask(id string, task Task) (Task, error) {
	query := "UPDATE tasks SET title = $1, description = $2 WHERE id = $3 RETURNING id"
	err := r.DB.QueryRow(query, task.Title, task.Description, id).Scan(&task.ID)
	return task, err
}

// DeleteTask implementation
func (r *PostgresTaskRepository) DeleteTask(id string) error {
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
