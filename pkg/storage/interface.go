package storage

import "Mod30/pkg/storage/postgres"

type DBInterface interface {
	Tasks(int, int) ([]postgres.Task, error)
	NewTask(postgres.Task) (int, error)
	UpdateTask(postgres.Task) error
	DeleteTask(int) error
}
