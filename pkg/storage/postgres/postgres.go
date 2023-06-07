package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(dbCon string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), dbCon)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

func (s *Storage) Tasks(taskID, authorID int) ([]Task, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT 
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			($1 = 0 OR id = $1) AND
			($2 = 0 OR author_id = $2)
		ORDER BY id;
	`,
		taskID,
		authorID,
	)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		tasks = append(tasks, t)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return tasks, rows.Err()
}

// NewTask создаёт новую задачу и возвращает её id.
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2) RETURNING id;
		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// обновление задачи
func (s *Storage) UpdateTask(t Task) error {
	if t.ID == 0 {
		return fmt.Errorf("ID не может быть 0")
	}
	//читаем, что есть уже в задаче с заданным ID
	var t_fromDB Task
	err := s.db.QueryRow(context.Background(), `
		SELECT
			id,
			opened,
			closed,
			author_id,
			assigned_id,
			title,
			content
		FROM tasks
		WHERE
			 id = $1;
	`,
		t.ID,
	).Scan(
		&t_fromDB.ID,
		&t_fromDB.Opened,
		&t_fromDB.Closed,
		&t_fromDB.AuthorID,
		&t_fromDB.AssignedID,
		&t_fromDB.Title,
		&t_fromDB.Content)
	if err != nil {
		return err
	}

	//проверяем, что изменяется
	if t.AssignedID == 0 {
		t.AssignedID = t_fromDB.AssignedID
	}
	if t.AuthorID == 0 {
		t.AuthorID = t_fromDB.AuthorID
	}
	if t.Closed == 0 {
		t.Closed = t_fromDB.Closed
	}
	if t.Content == "" {
		t.Content = t_fromDB.Content
	}
	if t.Title == "" {
		t.Title = t_fromDB.Title
	}
	if t.Opened == 0 {
		t.Opened = t_fromDB.Opened
	}
	//собствено меняем
	_, err = s.db.Exec(context.Background(), `
	UPDATE tasks
	SET  opened = $1, closed = $2, author_id = $3, assigned_id=$4, title=$5, content=$6
	WHERE id = $7;
	`,
		t.Opened,
		t.Closed,
		t.AuthorID,
		t.AssignedID,
		t.Title,
		t.Content,
		t.ID,
	)

	return err
}

// удаление задачи по ее ID
func (s *Storage) DeleteTask(taskID int) error {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM tasks
	WHERE id = $1;
	`,
		taskID,
	)
	return err
}
