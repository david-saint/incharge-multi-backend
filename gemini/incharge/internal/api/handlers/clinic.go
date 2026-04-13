package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"incharge/internal/database"
	"incharge/internal/models"
	"incharge/internal/utils"

	"github.com/gin-gonic/gin"
)

type ClinicHandler struct{}

func (h *ClinicHandler) ListClinics(c *gin.Context) {
	var clinics []models.Clinic
	query := database.DB.Model(&models.Clinic{})

	// Handle Preloads
	with := c.Query("with")
	if with != "" {
		if strings.Contains(with, "locations") {
			query = query.Preload("Locations")
		}
	}

	// Handle Trashed
	if c.Query("withTrashed") == "true" {
		query = query.Unscoped()
	} else if c.Query("onlyTrashed") == "true" {
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	}

	// Handle Search
	search := c.Query("search")
	if search != "" {
		// Complex search across fields and potentially relations
		query = query.Where("name LIKE ? OR address LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Handle Lat/Lng + Radius (Haversine)
	latStr := c.Query("latitude")
	lngStr := c.Query("longitude")

	if latStr != "" && lngStr != "" {
		lat, errLat := strconv.ParseFloat(latStr, 64)
		lng, errLng := strconv.ParseFloat(lngStr, 64)

		if errLat == nil && errLng == nil {
			radiusStr := c.Query("radius")
			radius := 10.0 // Default 10
			if radiusStr != "" {
				r, err := strconv.ParseFloat(radiusStr, 64)
				if err == nil {
					radius = r
				}
			}

			mode := c.Query("mode")
			if mode == "" {
				mode = "km"
			}

			multiplier := 6371.0 // km
			if mode == "mi" {
				multiplier = 3959.0
			}

			// Haversine formula in SQL
			haversine := fmt.Sprintf(
				"(%f * acos(cos(radians(%f)) * cos(radians(latitude)) * cos(radians(longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(latitude))))",
				multiplier, lat, lng, lat,
			)

			query = query.Select("*, " + haversine + " as actual_distance")
			query = query.Having("actual_distance <= ?", radius)

			// Handle Haversine specific sorting
			sort := c.Query("sort")
			if sort != "" {
				sortParts := strings.Split(sort, "|")
				if len(sortParts) == 2 && sortParts[0] == "distance" {
					dir := "asc"
					if strings.ToLower(sortParts[1]) == "desc" {
						dir = "desc"
					}
					query = query.Order("actual_distance " + dir)
				}
			}
		}
	} else {
		// Handle normal sort
		sort := c.Query("sort")
		if sort != "" {
			sortParts := strings.Split(sort, ",")
			for _, part := range sortParts {
				p := strings.Split(part, "|")
				if len(p) == 2 {
					query = query.Order(fmt.Sprintf("%s %s", p[0], p[1]))
				}
			}
		}
	}

	// Handle Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if err != nil {
		// If absent, returns all results, but for safety let's use a large number if parsing fails
		perPage = 1000
	}

	if perPage > 0 {
		offset := (page - 1) * perPage
		query = query.Offset(offset).Limit(perPage)
	}

	if err := query.Find(&clinics).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve clinics")
		return
	}

	// Post-process clinics to format distance string
	if latStr != "" && lngStr != "" {
		mode := c.Query("mode")
		if mode == "" {
			mode = "km"
		}
		for i := range clinics {
			if clinics[i].ActualDistance != nil {
				distStr := fmt.Sprintf("%.2f%s", *clinics[i].ActualDistance, mode)
				clinics[i].Distance = &distStr
			}
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Clinics retrieved", clinics)
}
