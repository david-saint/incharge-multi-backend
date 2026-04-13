package database

import (
	"log/slog"

	"github.com/incharge/server/internal/model"
	"gorm.io/gorm"
)

// AutoMigrate runs GORM auto-migrations to create/update all tables.
func AutoMigrate(db *gorm.DB) error {
	slog.Info("running database auto-migration")
	return db.AutoMigrate(
		&model.Country{},
		&model.State{},
		&model.Location{},
		&model.EducationLevel{},
		&model.ContraceptionReason{},
		&model.User{},
		&model.Profile{},
		&model.Clinic{},
		&model.Locatable{},
		&model.FaqGroup{},
		&model.Faq{},
		&model.Algorithm{},
		&model.Admin{},
		&model.PasswordReset{},
	)
}

// Seed populates reference data if tables are empty.
func Seed(db *gorm.DB) error {
	slog.Info("checking seed data")

	if err := seedContraceptionReasons(db); err != nil {
		return err
	}
	if err := seedEducationLevels(db); err != nil {
		return err
	}
	if err := seedFaqGroups(db); err != nil {
		return err
	}
	if err := seedCountries(db); err != nil {
		return err
	}
	if err := seedStates(db); err != nil {
		return err
	}

	slog.Info("seed data check complete")
	return nil
}

func seedContraceptionReasons(db *gorm.DB) error {
	var count int64
	db.Model(&model.ContraceptionReason{}).Count(&count)
	if count > 0 {
		return nil
	}

	slog.Info("seeding contraception reasons")
	reasons := []model.ContraceptionReason{
		{Value: "Completed family size"},
		{Value: "Child Spacing"},
		{Value: "Sexually Active with no plan for children at the moment"},
	}
	return db.Create(&reasons).Error
}

func seedEducationLevels(db *gorm.DB) error {
	var count int64
	db.Model(&model.EducationLevel{}).Count(&count)
	if count > 0 {
		return nil
	}

	slog.Info("seeding education levels")
	names := []string{
		"BArch", "BA", "B.Sc", "B.ENG", "LLB", "HNC", "HND",
		"ND", "M.Sc", "M.Eng", "Phd.", "Prof", "B.Tech", "Other",
	}
	levels := make([]model.EducationLevel, len(names))
	for i, n := range names {
		levels[i] = model.EducationLevel{Name: n}
	}
	return db.Create(&levels).Error
}

func seedFaqGroups(db *gorm.DB) error {
	var count int64
	db.Model(&model.FaqGroup{}).Count(&count)
	if count > 0 {
		return nil
	}

	slog.Info("seeding FAQ groups")
	names := []string{
		"Barrier Method", "Combined Oral Contraceptives",
		"Diaphragms and Spermicides", "Emergency Contraceptive Pills",
		"Female Sterilization", "Fertility Awareness", "Implants",
		"Injectables", "IUCD", "Lactational Amenorrhea",
		"Progestin Only Pills", "STIs", "Vasectomy",
	}
	groups := make([]model.FaqGroup, len(names))
	for i, n := range names {
		groups[i] = model.FaqGroup{Name: n}
	}
	return db.Create(&groups).Error
}

func seedCountries(db *gorm.DB) error {
	var count int64
	db.Model(&model.Country{}).Count(&count)
	if count > 0 {
		return nil
	}

	slog.Info("seeding countries")
	countries := []model.Country{
		{Name: "Nigeria", Code: "NG"},
		{Name: "United States", Code: "US"},
		{Name: "United Kingdom", Code: "GB"},
		{Name: "Ghana", Code: "GH"},
		{Name: "South Africa", Code: "ZA"},
		{Name: "Kenya", Code: "KE"},
		{Name: "Canada", Code: "CA"},
	}
	return db.Create(&countries).Error
}

func seedStates(db *gorm.DB) error {
	var count int64
	db.Model(&model.State{}).Count(&count)
	if count > 0 {
		return nil
	}

	slog.Info("seeding Nigerian states")
	states := []model.State{
		{Name: "Abia", Slug: "abia"},
		{Name: "Adamawa", Slug: "adamawa"},
		{Name: "Akwa Ibom", Slug: "akwa-ibom"},
		{Name: "Anambra", Slug: "anambra"},
		{Name: "Bauchi", Slug: "bauchi"},
		{Name: "Bayelsa", Slug: "bayelsa"},
		{Name: "Benue", Slug: "benue"},
		{Name: "Borno", Slug: "borno"},
		{Name: "Cross River", Slug: "cross-river"},
		{Name: "Delta", Slug: "delta"},
		{Name: "Ebonyi", Slug: "ebonyi"},
		{Name: "Edo", Slug: "edo"},
		{Name: "Ekiti", Slug: "ekiti"},
		{Name: "Enugu", Slug: "enugu"},
		{Name: "FCT", Slug: "fct"},
		{Name: "Gombe", Slug: "gombe"},
		{Name: "Imo", Slug: "imo"},
		{Name: "Jigawa", Slug: "jigawa"},
		{Name: "Kaduna", Slug: "kaduna"},
		{Name: "Kano", Slug: "kano"},
		{Name: "Katsina", Slug: "katsina"},
		{Name: "Kebbi", Slug: "kebbi"},
		{Name: "Kogi", Slug: "kogi"},
		{Name: "Kwara", Slug: "kwara"},
		{Name: "Lagos", Slug: "lagos"},
		{Name: "Nasarawa", Slug: "nasarawa"},
		{Name: "Niger", Slug: "niger"},
		{Name: "Ogun", Slug: "ogun"},
		{Name: "Ondo", Slug: "ondo"},
		{Name: "Osun", Slug: "osun"},
		{Name: "Oyo", Slug: "oyo"},
		{Name: "Plateau", Slug: "plateau"},
		{Name: "Rivers", Slug: "rivers"},
		{Name: "Sokoto", Slug: "sokoto"},
		{Name: "Taraba", Slug: "taraba"},
		{Name: "Yobe", Slug: "yobe"},
		{Name: "Zamfara", Slug: "zamfara"},
	}
	return db.Create(&states).Error
}
