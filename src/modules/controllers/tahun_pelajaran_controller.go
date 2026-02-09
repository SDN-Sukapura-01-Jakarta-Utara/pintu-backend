package controllers

import (
	"pintu-backend/src/modules/models"
	"pintu-backend/src/modules/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TahunPelajaranController handles HTTP requests for TahunPelajaran
type TahunPelajaranController struct {
	service services.TahunPelajaranService
}

// NewTahunPelajaranController creates a new TahunPelajaran controller
func NewTahunPelajaranController(service services.TahunPelajaranService) *TahunPelajaranController {
	return &TahunPelajaranController{service: service}
}

// Create creates a new TahunPelajaran
// @Summary Create new TahunPelajaran
// @Description Create a new TahunPelajaran
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Router /tahun_pelajaran [post]
func (c *TahunPelajaranController) Create(ctx *gin.Context) {
	var req models.TahunPelajaran
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Create(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": req})
}

// GetByID retrieves TahunPelajaran by ID
// @Summary Get TahunPelajaran by ID
// @Description Retrieve tahun_pelajaran details by ID
// @Tags TahunPelajaran
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404
// @Router /TahunPelajaran/{id} [get]
func (c *TahunPelajaranController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	data, err := c.service.GetByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// GetAll retrieves all tahun_pelajaran
// @Summary Get all TahunPelajaran
// @Description Retrieve all TahunPelajaran records
// @Tags TahunPelajaran
// @Produce json
// @Success 200
// @Failure 500
// @Router /tahun_pelajaran [get]
func (c *TahunPelajaranController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}

// Update updates TahunPelajaran
// @Summary Update TahunPelajaran
// @Description Update TahunPelajaran details
// @Tags tahun_pelajaran
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 400
// @Failure 404
// @Router /TahunPelajaran/{id} [put]
func (c *TahunPelajaranController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.TahunPelajaran
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = uint(id)
	if err := c.service.Update(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": req})
}

// Delete deletes TahunPelajaran by ID
// @Summary Delete tahun_pelajaran
// @Description Delete TahunPelajaran by ID
// @Tags TahunPelajaran
// @Produce json
// @Param id path int true "ID"
// @Success 200
// @Failure 404
// @Router /TahunPelajaran/{id} [delete]
func (c *TahunPelajaranController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
