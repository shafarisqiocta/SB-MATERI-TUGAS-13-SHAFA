package main

import (
	"bioskop-api/config"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama"`
	Lokasi string  `json:"lokasi"`
	Rating float64 `json:"rating"`
}

func main() {
	r := gin.Default()
	//req POST
	r.POST("/bioskop", func(c *gin.Context) {
		var input Bioskop

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Nama == "" || input.Lokasi == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Nama dan Lokasi tidak boleh kosong",
			})
			return
		}

		query := `INSERT INTO bioskop (nama, lokasi, rating)
			  VALUES ($1, $2, $3) RETURNING id`

		err := config.DB.QueryRow(
			query,
			input.Nama,
			input.Lokasi,
			input.Rating,
		).Scan(&input.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, input)
	})
	//req GET all data
	r.GET("/bioskop", func(c *gin.Context) {
		rows, err := config.DB.Query("SELECT id,nama,lokasi,rating FROM bioskop")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var bioskops []Bioskop

		for rows.Next() {
			var b Bioskop
			err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			bioskops = append(bioskops, b)
		}
		c.JSON(http.StatusOK, bioskops)
	})
	//req GET berdasarkan ID
	r.GET("/bioskop/:id", func(c *gin.Context) {
		id := c.Param("id")

		var b Bioskop
		query := "SELECT id, nama, lokasi, rating FROM bioskop WHERE id=$1"
		err := config.DB.QueryRow(query, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, b)
	})
	//req Update
	r.PUT("bioskop/:id", func(c *gin.Context) {
		id := c.Param("id")

		var input Bioskop
		if err := c.ShouldBindBodyWithJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.Nama == "" || input.Lokasi == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Nama dan Lokasi tidak boleh kosong",
			})
			return
		}

		query := `UPDATE bioskop SET nama=$1, lokasi=$2, rating=$3 WHERE id=$4`

		result, err := config.DB.Exec(query, input.Nama, input.Lokasi, input.Rating, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate"})

	})
	//Req DELETE
	r.DELETE("/bioskop/:id", func(c *gin.Context) {
		id := c.Param("id")

		result, err := config.DB.Exec("DELETE FROM bioskop WHERE id=$1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Message": "Data berhasil dihapus"})

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default untuk lokal
	}

	r.Run(":" + port)
}
