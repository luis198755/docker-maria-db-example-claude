package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type CenturionProject struct {
	ID           int             `json:"id"`
	CenturionID  string          `json:"centurion_id"`
	Matricula    string          `json:"matricula"`
	Ubicacion    string          `json:"ubicacion"`
	Proyecto     string          `json:"proyecto"`
	FechaEntrega string          `json:"fecha_entrega"`
	Programa     json.RawMessage `json:"programa"`
	Password     string          `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(db:3306)/centurion_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/projects", getProjects)
	r.GET("/projects/:id", getProject)
	r.POST("/projects", createProject)
	r.PUT("/projects/:id", updateProject)
	r.DELETE("/projects/:id", deleteProject)
	r.PATCH("/projects/:id/programa", updatePrograma) // New route for updating programa

	r.Run(":8080")
}

func getProjects(c *gin.Context) {
	rows, err := db.Query("SELECT id, centurion_id, matricula, ubicacion, proyecto, fecha_entrega, programa, password FROM centurion_projects")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var projects []CenturionProject
	for rows.Next() {
		var p CenturionProject
		var programJSON sql.NullString
		var fechaEntrega []uint8
		if err := rows.Scan(&p.ID, &p.CenturionID, &p.Matricula, &p.Ubicacion, &p.Proyecto, &fechaEntrega, &programJSON, &p.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		p.FechaEntrega = string(fechaEntrega)
		if programJSON.Valid {
			p.Programa = json.RawMessage(programJSON.String)
		}
		projects = append(projects, p)
	}

	c.JSON(http.StatusOK, projects)
}

func getProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var p CenturionProject
	var programJSON sql.NullString
	var fechaEntrega []uint8
	err = db.QueryRow("SELECT id, centurion_id, matricula, ubicacion, proyecto, fecha_entrega, programa, password FROM centurion_projects WHERE id = ?", id).Scan(&p.ID, &p.CenturionID, &p.Matricula, &p.Ubicacion, &p.Proyecto, &fechaEntrega, &programJSON, &p.Password)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p.FechaEntrega = string(fechaEntrega)
	if programJSON.Valid {
		p.Programa = json.RawMessage(programJSON.String)
	}

	c.JSON(http.StatusOK, p)
}

func createProject(c *gin.Context) {
	var p CenturionProject
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	programJSON, err := json.Marshal(p.Programa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Programa JSON: " + err.Error()})
		return
	}

	fechaEntrega, err := time.Parse("2006-01-02", p.FechaEntrega)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	result, err := db.Exec("INSERT INTO centurion_projects (centurion_id, matricula, ubicacion, proyecto, fecha_entrega, programa, password) VALUES (?, ?, ?, ?, ?, ?, ?)",
		p.CenturionID, p.Matricula, p.Ubicacion, p.Proyecto, fechaEntrega, string(programJSON), p.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	c.JSON(http.StatusCreated, p)
}

func updateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var p CenturionProject
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	var programJSON sql.NullString
	if p.Programa != nil {
		jsonBytes, err := json.Marshal(p.Programa)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Programa JSON: " + err.Error()})
			return
		}
		programJSON = sql.NullString{String: string(jsonBytes), Valid: true}
	}

	fechaEntrega, err := time.Parse("2006-01-02", p.FechaEntrega)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	_, err = db.Exec("UPDATE centurion_projects SET centurion_id = ?, matricula = ?, ubicacion = ?, proyecto = ?, fecha_entrega = ?, programa = ?, password = ? WHERE id = ?",
		p.CenturionID, p.Matricula, p.Ubicacion, p.Proyecto, fechaEntrega, programJSON, p.Password, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func deleteProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := db.Exec("DELETE FROM centurion_projects WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

type ProgramaUpdate struct {
	Fases      map[string][]int               `json:"fases"`
	Escenarios map[string][]int               `json:"escenarios"`
	Ciclos     map[string][]int               `json:"ciclos"`
	Eventos    map[string][]int               `json:"eventos"`
}

func updatePrograma(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var programaUpdate ProgramaUpdate

	if err := c.ShouldBindJSON(&programaUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Convert the programaUpdate to JSON
	programJSON, err := json.Marshal(programaUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Programa structure: " + err.Error()})
		return
	}

	// Update the database
	_, err = db.Exec("UPDATE centurion_projects SET programa = ? WHERE id = ?", string(programJSON), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Fetch the updated project to return
	var updatedProject CenturionProject
	var fetchedProgramJSON sql.NullString
	var fechaEntrega []uint8
	err = db.QueryRow("SELECT id, centurion_id, matricula, ubicacion, proyecto, fecha_entrega, programa, password FROM centurion_projects WHERE id = ?", id).Scan(
		&updatedProject.ID, &updatedProject.CenturionID, &updatedProject.Matricula, &updatedProject.Ubicacion, 
		&updatedProject.Proyecto, &fechaEntrega, &fetchedProgramJSON, &updatedProject.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated project: " + err.Error()})
		return
	}

	updatedProject.FechaEntrega = string(fechaEntrega)
	if fetchedProgramJSON.Valid {
		updatedProject.Programa = json.RawMessage(fetchedProgramJSON.String)
	}

	c.JSON(http.StatusOK, updatedProject)
}
