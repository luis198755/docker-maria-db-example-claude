package main

import (
	"encoding/json"  // Add this line
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
    ID      int             `json:"id"`
    Name    string          `json:"name"`
    Price   float64         `json:"price"`
    Program json.RawMessage `json:"program"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(db:3306)/ecommerce")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/products", getProducts)
	r.GET("/products/:id", getProduct)
	r.POST("/products", createProduct)
	r.PUT("/products/:id", updateProduct)
	r.DELETE("/products/:id", deleteProduct)

	log.Println("Server starting on port 8080...")
	r.Run(":8080")
}

func getProducts(c *gin.Context) {
    rows, err := db.Query("SELECT id, name, price, Program FROM products")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var p Product
        var programJSON sql.NullString
        if err := rows.Scan(&p.ID, &p.Name, &p.Price, &programJSON); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        if programJSON.Valid {
            p.Program = json.RawMessage(programJSON.String)
        }
        products = append(products, p)
    }

    c.JSON(http.StatusOK, products)
}

func getProduct(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var p Product
    var programJSON sql.NullString
    err = db.QueryRow("SELECT id, name, price, Program FROM products WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Price, &programJSON)
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if programJSON.Valid {
        p.Program = json.RawMessage(programJSON.String)
    }

    c.JSON(http.StatusOK, p)
}

func createProduct(c *gin.Context) {
    var p Product
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
        return
    }

    // Convert the Program field to a valid JSON string
    programJSON, err := json.Marshal(p.Program)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Program JSON: " + err.Error()})
        return
    }

    result, err := db.Exec("INSERT INTO products (name, price, Program) VALUES (?, ?, ?)", p.Name, p.Price, string(programJSON))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }

    id, _ := result.LastInsertId()
    p.ID = int(id)
    c.JSON(http.StatusCreated, p)
}

func updateProduct(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var p Product
    if err := c.ShouldBindJSON(&p); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
        return
    }

    // Convert the Program field to a valid JSON string
    var programJSON sql.NullString
    if p.Program != nil {
        jsonBytes, err := json.Marshal(p.Program)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Program JSON: " + err.Error()})
            return
        }
        programJSON = sql.NullString{String: string(jsonBytes), Valid: true}
    }

    _, err = db.Exec("UPDATE products SET name = ?, price = ?, Program = ? WHERE id = ?", p.Name, p.Price, programJSON, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
        return
    }

    // Fetch the updated product to return
    var updatedProduct Product
    var fetchedProgramJSON sql.NullString
    err = db.QueryRow("SELECT id, name, price, Program FROM products WHERE id = ?", id).Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Price, &fetchedProgramJSON)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching updated product: " + err.Error()})
        return
    }

    if fetchedProgramJSON.Valid {
        updatedProduct.Program = json.RawMessage(fetchedProgramJSON.String)
    }

    c.JSON(http.StatusOK, updatedProduct)
}

func deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}