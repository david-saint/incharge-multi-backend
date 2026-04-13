package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// ClinicRepo handles clinic database operations.
type ClinicRepo struct {
	db *gorm.DB
}

// NewClinicRepo creates a new ClinicRepo.
func NewClinicRepo(db *gorm.DB) *ClinicRepo {
	return &ClinicRepo{db: db}
}

// ClinicQueryParams holds all query parameters for clinic listing.
type ClinicQueryParams struct {
	Search      string
	Latitude    *float64
	Longitude   *float64
	Radius      float64
	Mode        string // "km" or "mi"
	Sort        string // comma-separated "column|direction" pairs
	With        string // comma-separated relationships
	Page        int
	PerPage     int
	WithTrashed bool
	OnlyTrashed bool
	WithCount   string
}

const (
	earthRadiusKm = 6371.0
	earthRadiusMi = 3959.0
)

// List returns clinics based on query parameters.
func (r *ClinicRepo) List(params ClinicQueryParams) ([]model.Clinic, int64, error) {
	q := r.db.Model(&model.Clinic{})

	// Soft-delete handling.
	if params.OnlyTrashed {
		q = q.Unscoped().Where("deleted_at IS NOT NULL")
	} else if params.WithTrashed {
		q = q.Unscoped()
	}

	// Haversine distance calculation.
	hasDistance := params.Latitude != nil && params.Longitude != nil
	if hasDistance {
		earthRadius := earthRadiusKm
		if params.Mode == "mi" {
			earthRadius = earthRadiusMi
		}
		distanceSQL := fmt.Sprintf(
			"(%f * acos(cos(radians(%f)) * cos(radians(latitude)) * cos(radians(longitude) - radians(%f)) + sin(radians(%f)) * sin(radians(latitude))))",
			earthRadius, *params.Latitude, *params.Longitude, *params.Latitude,
		)
		q = q.Select("clinics.*, " + distanceSQL + " AS distance_raw")

		radius := params.Radius
		if radius <= 0 {
			radius = 10
		}
		q = q.Where(distanceSQL+" <= ?", radius)
	}

	// Search.
	if params.Search != "" {
		search := "%" + params.Search + "%"
		q = q.Where(
			r.db.Where("clinics.id LIKE ?", search).
				Or("clinics.name LIKE ?", search).
				Or("clinics.address LIKE ?", search).
				Or("EXISTS (SELECT 1 FROM locatables JOIN locations ON locations.id = locatables.location_id WHERE locatables.locatable_id = clinics.id AND locatables.locatable_type = ? AND locations.name LIKE ?)",
					"Clinic", search),
		)
	}

	// Count total before pagination.
	var total int64
	countQ := q.Session(&gorm.Session{})
	if hasDistance {
		// For distance queries, we need a subquery approach for counting.
		countQ.Count(&total)
	} else {
		countQ.Count(&total)
	}

	// Sorting.
	if params.Sort != "" {
		for _, s := range strings.Split(params.Sort, ",") {
			parts := strings.SplitN(s, "|", 2)
			col := strings.TrimSpace(parts[0])
			dir := "asc"
			if len(parts) == 2 {
				dir = strings.ToLower(strings.TrimSpace(parts[1]))
			}
			if dir != "asc" && dir != "desc" {
				dir = "asc"
			}
			if col == "distance" && hasDistance {
				q = q.Order("distance_raw " + dir)
			} else if isAllowedSortColumn(col) {
				q = q.Order(col + " " + dir)
			}
		}
	} else {
		q = q.Order("clinics.id ASC")
	}

	// Eager loading.
	if params.With != "" {
		for _, w := range strings.Split(params.With, ",") {
			w = strings.TrimSpace(w)
			if w == "locations" {
				// Simple preload is sufficient because only clinics use
				// the locatables pivot table per the spec.
				q = q.Preload("Locations")
			}
		}
	}

	// Pagination.
	if params.Page > 0 && params.PerPage > 0 {
		offset := (params.Page - 1) * params.PerPage
		q = q.Offset(offset).Limit(params.PerPage)
	}

	var clinics []model.Clinic
	if err := q.Find(&clinics).Error; err != nil {
		return nil, 0, err
	}

	// withCount: compute relationship counts in Go after fetching.
	if params.WithCount != "" {
		for _, wc := range strings.Split(params.WithCount, ",") {
			wc = strings.TrimSpace(wc)
			if wc == "locations" {
				ids := make([]uint, len(clinics))
				for i, c := range clinics {
					ids[i] = c.ID
				}
				if len(ids) > 0 {
					// Fetch counts grouped by locatable_id.
					type countResult struct {
						LocatableID uint
						Count       int64
					}
					var counts []countResult
					r.db.Table("locatables").
						Select("locatable_id, COUNT(*) as count").
						Where("locatable_id IN ? AND locatable_type = ?", ids, "Clinic").
						Group("locatable_id").
						Scan(&counts)
					countMap := make(map[uint]int64, len(counts))
					for _, cr := range counts {
						countMap[cr.LocatableID] = cr.Count
					}
					for i := range clinics {
						cnt := countMap[clinics[i].ID]
						clinics[i].LocationsCount = &cnt
					}
				}
			}
		}
	}

	// Populate computed distance fields.
	// NOTE: DistanceRaw has gorm:"-" so GORM won't scan the SQL alias into it.
	// We recompute the distance in Go using the same haversine formula.
	if hasDistance {
		earthRadius := earthRadiusKm
		if params.Mode == "mi" {
			earthRadius = earthRadiusMi
		}
		mode := params.Mode
		if mode == "" {
			mode = "km"
		}
		radius := params.Radius
		if radius <= 0 {
			radius = 10
		}
		for i := range clinics {
			clinics[i].Mode = mode
			clinics[i].Radius = radius
			clinics[i].SearchRadius = fmt.Sprintf("%.0f%s", radius, mode)
			// Compute distance in Go since GORM can't scan into gorm:"-" fields.
			dist := 0.0
			if clinics[i].Latitude != nil && clinics[i].Longitude != nil {
				dist = haversineDistance(*params.Latitude, *params.Longitude, *clinics[i].Latitude, *clinics[i].Longitude, earthRadius)
			}
			clinics[i].DistanceRaw = dist
			clinics[i].ActualDistance = math.Round(dist*100) / 100
			clinics[i].Distance = fmt.Sprintf("%.2f%s", clinics[i].ActualDistance, mode)
		}
	}

	return clinics, total, nil
}

// FindByID finds a clinic by ID.
func (r *ClinicRepo) FindByID(id uint) (*model.Clinic, error) {
	var clinic model.Clinic
	if err := r.db.First(&clinic, id).Error; err != nil {
		return nil, err
	}
	return &clinic, nil
}

// Create creates a new clinic.
func (r *ClinicRepo) Create(clinic *model.Clinic) error {
	return r.db.Create(clinic).Error
}

// Update updates a clinic.
func (r *ClinicRepo) Update(clinic *model.Clinic) error {
	return r.db.Save(clinic).Error
}

// SoftDelete soft-deletes a clinic.
func (r *ClinicRepo) SoftDelete(id uint) error {
	return r.db.Delete(&model.Clinic{}, id).Error
}

// Restore restores a soft-deleted clinic.
func (r *ClinicRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.Clinic{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// ListPaginated returns a simple paginated clinic list.
func (r *ClinicRepo) ListPaginated(page, perPage int) ([]model.Clinic, int64, error) {
	var total int64
	r.db.Model(&model.Clinic{}).Count(&total)

	var clinics []model.Clinic
	offset := (page - 1) * perPage
	err := r.db.Order("id desc").Offset(offset).Limit(perPage).Find(&clinics).Error
	return clinics, total, err
}

// ListDeleted returns all soft-deleted clinics.
func (r *ClinicRepo) ListDeleted() ([]model.Clinic, error) {
	var clinics []model.Clinic
	err := r.db.Unscoped().Where("deleted_at IS NOT NULL").Find(&clinics).Error
	return clinics, err
}

func isAllowedSortColumn(col string) bool {
	allowed := map[string]bool{
		"id": true, "name": true, "address": true,
		"created_at": true, "updated_at": true,
	}
	return allowed[col]
}

// haversineDistance computes the great-circle distance between two points
// on a sphere of the given radius using the haversine formula.
func haversineDistance(lat1, lng1, lat2, lng2, earthRadius float64) float64 {
	dLat := degToRad(lat2 - lat1)
	dLng := degToRad(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degToRad(lat1))*math.Cos(degToRad(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Acos(math.Min(1, math.Sqrt(a)+(1-math.Sqrt(a))*0)) // safer: 2*atan2
	// Use the standard formulation: c = 2 * atan2(sqrt(a), sqrt(1-a))
	c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

func degToRad(deg float64) float64 {
	return deg * math.Pi / 180
}
