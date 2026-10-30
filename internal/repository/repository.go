package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"gestion-depenses/internal/models"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository { return &Repository{db: db} }

func Migrate(db *sql.DB) error {
	const schema = `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);
	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL CHECK(amount > 0),
		description TEXT NOT NULL,
		date TEXT NOT NULL,
		category_id INTEGER NOT NULL,
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);`
	if _, err := db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		return err
	}
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	_, err := db.Exec(`INSERT OR IGNORE INTO categories(name) VALUES
		('Alimentation'), ('Transport'), ('Logement'), ('Loisirs'), ('Santé'), ('Autres')`)
	return err
}

func (r *Repository) Categories() ([]models.Category, error) {
	rows, err := r.db.Query(`SELECT id, name FROM categories ORDER BY name`)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil { return nil, err }
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *Repository) GetAll(categoryID int64, from, to string) ([]models.Expense, error) {
	query := `SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name
		FROM expenses e JOIN categories c ON c.id=e.category_id WHERE 1=1`
	var args []any
	if categoryID > 0 {
		query += ` AND e.category_id=?`
		args = append(args, categoryID)
	}
	if from != "" {
		query += ` AND e.date>=?`
		args = append(args, from)
	}
	if to != "" {
		query += ` AND e.date<=?`
		args = append(args, to)
	}
	query += ` ORDER BY e.date DESC, e.id DESC`
	rows, err := r.db.Query(query, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.Category); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *Repository) GetByID(id int64) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(`SELECT id, amount, description, date, category_id FROM expenses WHERE id=?`, id).
		Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID)
	return e, err
}

func (r *Repository) Create(e models.Expense) error {
	_, err := r.db.Exec(`INSERT INTO expenses(amount, description, date, category_id) VALUES(?,?,?,?)`,
		e.Amount, e.Description, e.Date, e.CategoryID)
	return err
}

func (r *Repository) Update(e models.Expense) error {
	res, err := r.db.Exec(`UPDATE expenses SET amount=?, description=?, date=?, category_id=? WHERE id=?`,
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID)
	if err != nil { return err }
	n, err := res.RowsAffected()
	if err == nil && n == 0 { return sql.ErrNoRows }
	return err
}

func (r *Repository) Delete(id int64) error {
	res, err := r.db.Exec(`DELETE FROM expenses WHERE id=?`, id)
	if err != nil { return err }
	n, err := res.RowsAffected()
	if err == nil && n == 0 { return sql.ErrNoRows }
	return err
}

// BuildFilterDescription is used only for a human-readable diagnostic message.
func BuildFilterDescription(categoryID int64, from, to string) string {
	parts := []string{}
	if categoryID > 0 { parts = append(parts, fmt.Sprintf("catégorie=%d", categoryID)) }
	if from != "" { parts = append(parts, "à partir du "+from) }
	if to != "" { parts = append(parts, "jusqu'au "+to) }
	return strings.Join(parts, ", ")
}
