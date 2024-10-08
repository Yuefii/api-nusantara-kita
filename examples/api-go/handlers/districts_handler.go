package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuefii/NusantaraKita/examples/api-go/dtos"
	"github.com/yuefii/NusantaraKita/examples/api-go/models"
)

func (handler *Handler) GetDistricts(ctx *gin.Context) {
	var district []models.Districts
	var totalItems int64

	showAll := ctx.Query("show_all") == "true"
	pageStr := ctx.Query("page")
	perPageStr := ctx.Query("per_page")
	regencyCode := ctx.Query("regency_code")

	page := 1
	perPage := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 {
			perPage = pp
		}
	}

	if showAll {

		if regencyCode != "" {
			result := handler.db.Preload("Regency").Where("regency_code = ?", regencyCode).Find(&district)
			if result.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
		} else {
			result := handler.db.Preload("Regency").Find(&district)
			if result.Error != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
				return
			}
		}
		totalItems = int64(len(district))

		var districtDTOs []dtos.ApiDTO
		for _, response := range district {
			districtDTOs = append(districtDTOs, dtos.ApiDTO{
				Code: response.Code,
				Name: response.Name,
			})
		}

		response := gin.H{
			"data": districtDTOs,
		}

		ctx.JSON(http.StatusOK, response)
		return
	}

	query := handler.db.Model(&models.Districts{})
	if regencyCode != "" {
		query = query.Where("regency_code = ?", regencyCode)
	}

	result := query.Count(&totalItems)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	result = query.Preload("Regency").Offset((page - 1) * perPage).Limit(perPage).Find(&district)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var districtDTOs []dtos.ApiDTO
	for _, response := range district {
		districtDTOs = append(districtDTOs, dtos.ApiDTO{
			Code: response.Code,
			Name: response.Name,
		})
	}

	totalPages := int(totalItems) / perPage
	if int(totalItems)%perPage != 0 {
		totalPages++
	}

	response := dtos.ResponseDTO{
		Data: districtDTOs,
		Pagination: dtos.PaginationDTO{
			CurrentPage: page,
			PerPage:     perPage,
			TotalItems:  int(totalItems),
			TotalPages:  totalPages,
		},
	}
	ctx.JSON(http.StatusOK, response)
}
