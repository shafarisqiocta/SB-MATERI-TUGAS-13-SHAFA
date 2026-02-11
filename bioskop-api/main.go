package main

import (
	"bioskop-api/config"
	"net/http"

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

	config.ConnectDB()

	r.Run(":8080")
}
