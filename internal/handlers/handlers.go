package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gestion-depenses/internal/models"
	"gestion-depenses/internal/repository"
)

type Handler struct { repo *repository.Repository }

func New(repo *repository.Repository) *Handler { return &Handler{repo: repo} }

type PageData struct {
	Title string
	Expenses []models.Expense
	Categories []models.Category
	Expense models.Expense
	SelectedCategory int64
	From string
	To string
	Error string
}

func (h *Handler) render(w http.ResponseWriter, page string, data PageData) {
	t, err := template.ParseFiles("web/templates/layout.html", "web/templates/"+page)
	if err != nil { http.Error(w, "Erreur de template", http.StatusInternalServerError); return }
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "Erreur de rendu", http.StatusInternalServerError)
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	h.render(w, "about.html", PageData{Title: "À propos"})
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	var categoryID int64
	if raw := r.URL.Query().Get("category"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 0 { http.Error(w, "Catégorie invalide", http.StatusBadRequest); return }
		categoryID = id
	}
	from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")
	if !validOptionalDate(from) || !validOptionalDate(to) || (from != "" && to != "" && from > to) {
		http.Error(w, "Dates invalides", http.StatusBadRequest); return
	}
	expenses, err := h.repo.GetAll(categoryID, from, to)
	if err != nil { http.Error(w, "Impossible de lire les dépenses", http.StatusInternalServerError); return }
	categories, err := h.repo.Categories()
	if err != nil { http.Error(w, "Impossible de lire les catégories", http.StatusInternalServerError); return }
	h.render(w, "expenses.html", PageData{Title:"Mes dépenses", Expenses:expenses,
		Categories:categories, SelectedCategory:categoryID, From:from, To:to})
}

func (h *Handler) NewExpense(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.Categories()
	if err != nil { http.Error(w, "Erreur base de données", 500); return }
	h.render(w, "form.html", PageData{Title:"Ajouter une dépense", Categories:categories,
		Expense:models.Expense{Date:time.Now().Format("2006-01-02")}})
}

func (h *Handler) readExpense(r *http.Request) (models.Expense, error) {
	amount, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64)
	if err != nil || amount <= 0 { return models.Expense{}, strconv.ErrSyntax }
	categoryID, err := strconv.ParseInt(r.FormValue("category_id"), 10, 64)
	if err != nil || categoryID <= 0 { return models.Expense{}, strconv.ErrSyntax }
	date := r.FormValue("date")
	if !validOptionalDate(date) || date == "" { return models.Expense{}, strconv.ErrSyntax }
	description := strings.TrimSpace(r.FormValue("description"))
	if description == "" || len(description) > 250 { return models.Expense{}, strconv.ErrSyntax }
	return models.Expense{Amount:amount, Description:description, Date:date, CategoryID:categoryID}, nil
}

func validOptionalDate(value string) bool {
	if value == "" { return true }
	_, err := time.Parse("2006-01-02", value)
	return err == nil
}

func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { http.Error(w, "Formulaire invalide", 400); return }
	e, err := h.readExpense(r)
	if err != nil { http.Error(w, "Vérifiez le montant, la description, la date et la catégorie", 400); return }
	if err := h.repo.Create(e); err != nil { http.Error(w, "Enregistrement impossible", 500); return }
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func (h *Handler) EditExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 { http.NotFound(w, r); return }
	e, err := h.repo.GetByID(id)
	if err != nil { if err == sql.ErrNoRows { http.NotFound(w,r) } else { http.Error(w,"Erreur base de données",500) }; return }
	categories, err := h.repo.Categories()
	if err != nil { http.Error(w,"Erreur base de données",500); return }
	h.render(w, "form.html", PageData{Title:"Modifier une dépense", Expense:e, Categories:categories})
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 { http.NotFound(w,r); return }
	if err := r.ParseForm(); err != nil { http.Error(w,"Formulaire invalide",400); return }
	e, err := h.readExpense(r)
	if err != nil { http.Error(w,"Vérifiez les champs du formulaire",400); return }
	e.ID = id
	if err := h.repo.Update(e); err != nil { if err == sql.ErrNoRows { http.NotFound(w,r) } else { http.Error(w,"Modification impossible",500) }; return }
	http.Redirect(w,r,"/expenses",http.StatusSeeOther)
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"),10,64)
	if err != nil || id <= 0 { http.NotFound(w,r); return }
	if err := h.repo.Delete(id); err != nil { if err == sql.ErrNoRows { http.NotFound(w,r) } else { http.Error(w,"Suppression impossible",500) }; return }
	http.Redirect(w,r,"/expenses",http.StatusSeeOther)
}
