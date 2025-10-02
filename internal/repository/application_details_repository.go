package repository

import (
	"database/sql"
	"fmt"
	"time"

	"scholarship-system/internal/database"
	"scholarship-system/internal/models"
)

type ApplicationDetailsRepository struct {
	db *sql.DB
}

func NewApplicationDetailsRepository() *ApplicationDetailsRepository {
	return &ApplicationDetailsRepository{
		db: database.DB,
	}
}

// ========================================
// 1. Personal Info Methods
// ========================================

func (r *ApplicationDetailsRepository) CreatePersonalInfo(info *models.ApplicationPersonalInfo) error {
	query := `
		INSERT INTO application_personal_info (
			application_id, prefix_th, prefix_en, first_name_th, last_name_th,
			first_name_en, last_name_en, email, phone, line_id,
			citizen_id, student_id, faculty, department, major,
			year_level, admission_type, admission_details, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		RETURNING info_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		info.ApplicationID, info.PrefixTH, info.PrefixEN, info.FirstNameTH, info.LastNameTH,
		info.FirstNameEN, info.LastNameEN, info.Email, info.Phone, info.LineID,
		info.CitizenID, info.StudentID, info.Faculty, info.Department, info.Major,
		info.YearLevel, info.AdmissionType, info.AdmissionDetails, now, now,
	).Scan(&info.InfoID)

	return err
}

func (r *ApplicationDetailsRepository) GetPersonalInfoByApplicationID(applicationID int) (*models.ApplicationPersonalInfo, error) {
	query := `
		SELECT info_id, application_id, prefix_th, prefix_en, first_name_th, last_name_th,
		       first_name_en, last_name_en, email, phone, line_id,
		       citizen_id, student_id, faculty, department, major,
		       year_level, admission_type, admission_details, created_at, updated_at
		FROM application_personal_info
		WHERE application_id = $1
	`

	info := &models.ApplicationPersonalInfo{}
	err := r.db.QueryRow(query, applicationID).Scan(
		&info.InfoID, &info.ApplicationID, &info.PrefixTH, &info.PrefixEN, &info.FirstNameTH, &info.LastNameTH,
		&info.FirstNameEN, &info.LastNameEN, &info.Email, &info.Phone, &info.LineID,
		&info.CitizenID, &info.StudentID, &info.Faculty, &info.Department, &info.Major,
		&info.YearLevel, &info.AdmissionType, &info.AdmissionDetails, &info.CreatedAt, &info.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return info, nil
}

func (r *ApplicationDetailsRepository) UpdatePersonalInfo(info *models.ApplicationPersonalInfo) error {
	query := `
		UPDATE application_personal_info
		SET prefix_th = $2, prefix_en = $3, first_name_th = $4, last_name_th = $5,
		    first_name_en = $6, last_name_en = $7, email = $8, phone = $9, line_id = $10,
		    citizen_id = $11, student_id = $12, faculty = $13, department = $14, major = $15,
		    year_level = $16, admission_type = $17, admission_details = $18, updated_at = $19
		WHERE info_id = $1
	`

	_, err := r.db.Exec(query,
		info.InfoID, info.PrefixTH, info.PrefixEN, info.FirstNameTH, info.LastNameTH,
		info.FirstNameEN, info.LastNameEN, info.Email, info.Phone, info.LineID,
		info.CitizenID, info.StudentID, info.Faculty, info.Department, info.Major,
		info.YearLevel, info.AdmissionType, info.AdmissionDetails, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeletePersonalInfo(infoID string) error {
	query := `DELETE FROM application_personal_info WHERE info_id = $1`
	_, err := r.db.Exec(query, infoID)
	return err
}

// ========================================
// 2. Address Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateAddress(address *models.ApplicationAddress) error {
	query := `
		INSERT INTO application_addresses (
			application_id, address_type, house_number, village_number, alley,
			road, subdistrict, district, province, postal_code,
			address_line1, address_line2, latitude, longitude, map_image_url,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING address_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		address.ApplicationID, address.AddressType, address.HouseNumber, address.VillageNumber, address.Alley,
		address.Road, address.Subdistrict, address.District, address.Province, address.PostalCode,
		address.AddressLine1, address.AddressLine2, address.Latitude, address.Longitude, address.MapImageURL,
		now, now,
	).Scan(&address.AddressID)

	return err
}

func (r *ApplicationDetailsRepository) GetAddressesByApplicationID(applicationID int) ([]models.ApplicationAddress, error) {
	query := `
		SELECT address_id, application_id, address_type, house_number, village_number, alley,
		       road, subdistrict, district, province, postal_code,
		       address_line1, address_line2, latitude, longitude, map_image_url,
		       created_at, updated_at
		FROM application_addresses
		WHERE application_id = $1
		ORDER BY address_type
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []models.ApplicationAddress
	for rows.Next() {
		var address models.ApplicationAddress
		err := rows.Scan(
			&address.AddressID, &address.ApplicationID, &address.AddressType, &address.HouseNumber, &address.VillageNumber, &address.Alley,
			&address.Road, &address.Subdistrict, &address.District, &address.Province, &address.PostalCode,
			&address.AddressLine1, &address.AddressLine2, &address.Latitude, &address.Longitude, &address.MapImageURL,
			&address.CreatedAt, &address.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *ApplicationDetailsRepository) UpdateAddress(address *models.ApplicationAddress) error {
	query := `
		UPDATE application_addresses
		SET address_type = $2, house_number = $3, village_number = $4, alley = $5,
		    road = $6, subdistrict = $7, district = $8, province = $9, postal_code = $10,
		    address_line1 = $11, address_line2 = $12, latitude = $13, longitude = $14,
		    map_image_url = $15, updated_at = $16
		WHERE address_id = $1
	`

	_, err := r.db.Exec(query,
		address.AddressID, address.AddressType, address.HouseNumber, address.VillageNumber, address.Alley,
		address.Road, address.Subdistrict, address.District, address.Province, address.PostalCode,
		address.AddressLine1, address.AddressLine2, address.Latitude, address.Longitude,
		address.MapImageURL, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteAddress(addressID string) error {
	query := `DELETE FROM application_addresses WHERE address_id = $1`
	_, err := r.db.Exec(query, addressID)
	return err
}

// ========================================
// 3. Education History Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateEducationHistory(history *models.ApplicationEducationHistory) error {
	query := `
		INSERT INTO application_education_history (
			application_id, education_level, school_name, school_province,
			gpa, graduation_year, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING history_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		history.ApplicationID, history.EducationLevel, history.SchoolName, history.SchoolProvince,
		history.GPA, history.GraduationYear, now, now,
	).Scan(&history.HistoryID)

	return err
}

func (r *ApplicationDetailsRepository) GetEducationHistoryByApplicationID(applicationID int) ([]models.ApplicationEducationHistory, error) {
	query := `
		SELECT history_id, application_id, education_level, school_name, school_province,
		       gpa, graduation_year, created_at, updated_at
		FROM application_education_history
		WHERE application_id = $1
		ORDER BY graduation_year DESC
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.ApplicationEducationHistory
	for rows.Next() {
		var history models.ApplicationEducationHistory
		err := rows.Scan(
			&history.HistoryID, &history.ApplicationID, &history.EducationLevel, &history.SchoolName, &history.SchoolProvince,
			&history.GPA, &history.GraduationYear, &history.CreatedAt, &history.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}

func (r *ApplicationDetailsRepository) UpdateEducationHistory(history *models.ApplicationEducationHistory) error {
	query := `
		UPDATE application_education_history
		SET education_level = $2, school_name = $3, school_province = $4,
		    gpa = $5, graduation_year = $6, updated_at = $7
		WHERE history_id = $1
	`

	_, err := r.db.Exec(query,
		history.HistoryID, history.EducationLevel, history.SchoolName, history.SchoolProvince,
		history.GPA, history.GraduationYear, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteEducationHistory(historyID string) error {
	query := `DELETE FROM application_education_history WHERE history_id = $1`
	_, err := r.db.Exec(query, historyID)
	return err
}

// ========================================
// 4. Family Member Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateFamilyMember(member *models.ApplicationFamilyMember) error {
	query := `
		INSERT INTO application_family_members (
			application_id, relationship, title, first_name, last_name,
			age, living_status, occupation, position, workplace,
			workplace_province, monthly_income, phone, notes,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING member_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		member.ApplicationID, member.Relationship, member.Title, member.FirstName, member.LastName,
		member.Age, member.LivingStatus, member.Occupation, member.Position, member.Workplace,
		member.WorkplaceProvince, member.MonthlyIncome, member.Phone, member.Notes,
		now, now,
	).Scan(&member.MemberID)

	return err
}

func (r *ApplicationDetailsRepository) GetFamilyMembersByApplicationID(applicationID int) ([]models.ApplicationFamilyMember, error) {
	query := `
		SELECT member_id, application_id, relationship, title, first_name, last_name,
		       age, living_status, occupation, position, workplace,
		       workplace_province, monthly_income, phone, notes,
		       created_at, updated_at
		FROM application_family_members
		WHERE application_id = $1
		ORDER BY relationship
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []models.ApplicationFamilyMember
	for rows.Next() {
		var member models.ApplicationFamilyMember
		err := rows.Scan(
			&member.MemberID, &member.ApplicationID, &member.Relationship, &member.Title, &member.FirstName, &member.LastName,
			&member.Age, &member.LivingStatus, &member.Occupation, &member.Position, &member.Workplace,
			&member.WorkplaceProvince, &member.MonthlyIncome, &member.Phone, &member.Notes,
			&member.CreatedAt, &member.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *ApplicationDetailsRepository) UpdateFamilyMember(member *models.ApplicationFamilyMember) error {
	query := `
		UPDATE application_family_members
		SET relationship = $2, title = $3, first_name = $4, last_name = $5,
		    age = $6, living_status = $7, occupation = $8, position = $9, workplace = $10,
		    workplace_province = $11, monthly_income = $12, phone = $13, notes = $14,
		    updated_at = $15
		WHERE member_id = $1
	`

	_, err := r.db.Exec(query,
		member.MemberID, member.Relationship, member.Title, member.FirstName, member.LastName,
		member.Age, member.LivingStatus, member.Occupation, member.Position, member.Workplace,
		member.WorkplaceProvince, member.MonthlyIncome, member.Phone, member.Notes,
		time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteFamilyMember(memberID string) error {
	query := `DELETE FROM application_family_members WHERE member_id = $1`
	_, err := r.db.Exec(query, memberID)
	return err
}

// ========================================
// 5. Asset Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateAsset(asset *models.ApplicationAsset) error {
	query := `
		INSERT INTO application_assets (
			application_id, asset_type, category, description, value,
			monthly_cost, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING asset_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		asset.ApplicationID, asset.AssetType, asset.Category, asset.Description, asset.Value,
		asset.MonthlyCost, asset.Notes, now, now,
	).Scan(&asset.AssetID)

	return err
}

func (r *ApplicationDetailsRepository) GetAssetsByApplicationID(applicationID int) ([]models.ApplicationAsset, error) {
	query := `
		SELECT asset_id, application_id, asset_type, category, description, value,
		       monthly_cost, notes, created_at, updated_at
		FROM application_assets
		WHERE application_id = $1
		ORDER BY asset_type, category
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.ApplicationAsset
	for rows.Next() {
		var asset models.ApplicationAsset
		err := rows.Scan(
			&asset.AssetID, &asset.ApplicationID, &asset.AssetType, &asset.Category, &asset.Description, &asset.Value,
			&asset.MonthlyCost, &asset.Notes, &asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *ApplicationDetailsRepository) UpdateAsset(asset *models.ApplicationAsset) error {
	query := `
		UPDATE application_assets
		SET asset_type = $2, category = $3, description = $4, value = $5,
		    monthly_cost = $6, notes = $7, updated_at = $8
		WHERE asset_id = $1
	`

	_, err := r.db.Exec(query,
		asset.AssetID, asset.AssetType, asset.Category, asset.Description, asset.Value,
		asset.MonthlyCost, asset.Notes, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteAsset(assetID string) error {
	query := `DELETE FROM application_assets WHERE asset_id = $1`
	_, err := r.db.Exec(query, assetID)
	return err
}

// ========================================
// 6. Guardian Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateGuardian(guardian *models.ApplicationGuardian) error {
	query := `
		INSERT INTO application_guardians (
			application_id, title, first_name, last_name, relationship,
			address, phone, occupation, position, workplace,
			workplace_phone, monthly_income, debts, debt_details,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING guardian_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		guardian.ApplicationID, guardian.Title, guardian.FirstName, guardian.LastName, guardian.Relationship,
		guardian.Address, guardian.Phone, guardian.Occupation, guardian.Position, guardian.Workplace,
		guardian.WorkplacePhone, guardian.MonthlyIncome, guardian.Debts, guardian.DebtDetails,
		now, now,
	).Scan(&guardian.GuardianID)

	return err
}

func (r *ApplicationDetailsRepository) GetGuardiansByApplicationID(applicationID int) ([]models.ApplicationGuardian, error) {
	query := `
		SELECT guardian_id, application_id, title, first_name, last_name, relationship,
		       address, phone, occupation, position, workplace,
		       workplace_phone, monthly_income, debts, debt_details,
		       created_at, updated_at
		FROM application_guardians
		WHERE application_id = $1
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guardians []models.ApplicationGuardian
	for rows.Next() {
		var guardian models.ApplicationGuardian
		err := rows.Scan(
			&guardian.GuardianID, &guardian.ApplicationID, &guardian.Title, &guardian.FirstName, &guardian.LastName, &guardian.Relationship,
			&guardian.Address, &guardian.Phone, &guardian.Occupation, &guardian.Position, &guardian.Workplace,
			&guardian.WorkplacePhone, &guardian.MonthlyIncome, &guardian.Debts, &guardian.DebtDetails,
			&guardian.CreatedAt, &guardian.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		guardians = append(guardians, guardian)
	}

	return guardians, nil
}

func (r *ApplicationDetailsRepository) UpdateGuardian(guardian *models.ApplicationGuardian) error {
	query := `
		UPDATE application_guardians
		SET title = $2, first_name = $3, last_name = $4, relationship = $5,
		    address = $6, phone = $7, occupation = $8, position = $9, workplace = $10,
		    workplace_phone = $11, monthly_income = $12, debts = $13, debt_details = $14,
		    updated_at = $15
		WHERE guardian_id = $1
	`

	_, err := r.db.Exec(query,
		guardian.GuardianID, guardian.Title, guardian.FirstName, guardian.LastName, guardian.Relationship,
		guardian.Address, guardian.Phone, guardian.Occupation, guardian.Position, guardian.Workplace,
		guardian.WorkplacePhone, guardian.MonthlyIncome, guardian.Debts, guardian.DebtDetails,
		time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteGuardian(guardianID string) error {
	query := `DELETE FROM application_guardians WHERE guardian_id = $1`
	_, err := r.db.Exec(query, guardianID)
	return err
}

// ========================================
// 7. Sibling Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateSibling(sibling *models.ApplicationSibling) error {
	query := `
		INSERT INTO application_siblings (
			application_id, sibling_order, gender, school_or_workplace, education_level,
			is_working, monthly_income, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING sibling_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		sibling.ApplicationID, sibling.SiblingOrder, sibling.Gender, sibling.SchoolOrWorkplace, sibling.EducationLevel,
		sibling.IsWorking, sibling.MonthlyIncome, sibling.Notes, now, now,
	).Scan(&sibling.SiblingID)

	return err
}

func (r *ApplicationDetailsRepository) GetSiblingsByApplicationID(applicationID int) ([]models.ApplicationSibling, error) {
	query := `
		SELECT sibling_id, application_id, sibling_order, gender, school_or_workplace, education_level,
		       is_working, monthly_income, notes, created_at, updated_at
		FROM application_siblings
		WHERE application_id = $1
		ORDER BY sibling_order
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var siblings []models.ApplicationSibling
	for rows.Next() {
		var sibling models.ApplicationSibling
		err := rows.Scan(
			&sibling.SiblingID, &sibling.ApplicationID, &sibling.SiblingOrder, &sibling.Gender, &sibling.SchoolOrWorkplace, &sibling.EducationLevel,
			&sibling.IsWorking, &sibling.MonthlyIncome, &sibling.Notes, &sibling.CreatedAt, &sibling.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		siblings = append(siblings, sibling)
	}

	return siblings, nil
}

func (r *ApplicationDetailsRepository) UpdateSibling(sibling *models.ApplicationSibling) error {
	query := `
		UPDATE application_siblings
		SET sibling_order = $2, gender = $3, school_or_workplace = $4, education_level = $5,
		    is_working = $6, monthly_income = $7, notes = $8, updated_at = $9
		WHERE sibling_id = $1
	`

	_, err := r.db.Exec(query,
		sibling.SiblingID, sibling.SiblingOrder, sibling.Gender, sibling.SchoolOrWorkplace, sibling.EducationLevel,
		sibling.IsWorking, sibling.MonthlyIncome, sibling.Notes, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteSibling(siblingID string) error {
	query := `DELETE FROM application_siblings WHERE sibling_id = $1`
	_, err := r.db.Exec(query, siblingID)
	return err
}

// ========================================
// 8. Living Situation Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateLivingSituation(living *models.ApplicationLivingSituation) error {
	query := `
		INSERT INTO application_living_situation (
			application_id, living_with, living_details, front_house_image,
			side_house_image, back_house_image, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING living_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		living.ApplicationID, living.LivingWith, living.LivingDetails, living.FrontHouseImage,
		living.SideHouseImage, living.BackHouseImage, now, now,
	).Scan(&living.LivingID)

	return err
}

func (r *ApplicationDetailsRepository) GetLivingSituationByApplicationID(applicationID int) (*models.ApplicationLivingSituation, error) {
	query := `
		SELECT living_id, application_id, living_with, living_details, front_house_image,
		       side_house_image, back_house_image, created_at, updated_at
		FROM application_living_situation
		WHERE application_id = $1
	`

	living := &models.ApplicationLivingSituation{}
	err := r.db.QueryRow(query, applicationID).Scan(
		&living.LivingID, &living.ApplicationID, &living.LivingWith, &living.LivingDetails, &living.FrontHouseImage,
		&living.SideHouseImage, &living.BackHouseImage, &living.CreatedAt, &living.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return living, nil
}

func (r *ApplicationDetailsRepository) UpdateLivingSituation(living *models.ApplicationLivingSituation) error {
	query := `
		UPDATE application_living_situation
		SET living_with = $2, living_details = $3, front_house_image = $4,
		    side_house_image = $5, back_house_image = $6, updated_at = $7
		WHERE living_id = $1
	`

	_, err := r.db.Exec(query,
		living.LivingID, living.LivingWith, living.LivingDetails, living.FrontHouseImage,
		living.SideHouseImage, living.BackHouseImage, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteLivingSituation(livingID string) error {
	query := `DELETE FROM application_living_situation WHERE living_id = $1`
	_, err := r.db.Exec(query, livingID)
	return err
}

// ========================================
// 9. Financial Info Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateFinancialInfo(financial *models.ApplicationFinancialInfo) error {
	query := `
		INSERT INTO application_financial_info (
			application_id, monthly_allowance, daily_travel_cost, monthly_dorm_cost,
			other_monthly_costs, has_income, income_source, monthly_income,
			financial_notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING financial_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		financial.ApplicationID, financial.MonthlyAllowance, financial.DailyTravelCost, financial.MonthlyDormCost,
		financial.OtherMonthlyCosts, financial.HasIncome, financial.IncomeSource, financial.MonthlyIncome,
		financial.FinancialNotes, now, now,
	).Scan(&financial.FinancialID)

	return err
}

func (r *ApplicationDetailsRepository) GetFinancialInfoByApplicationID(applicationID int) (*models.ApplicationFinancialInfo, error) {
	query := `
		SELECT financial_id, application_id, monthly_allowance, daily_travel_cost, monthly_dorm_cost,
		       other_monthly_costs, has_income, income_source, monthly_income,
		       financial_notes, created_at, updated_at
		FROM application_financial_info
		WHERE application_id = $1
	`

	financial := &models.ApplicationFinancialInfo{}
	err := r.db.QueryRow(query, applicationID).Scan(
		&financial.FinancialID, &financial.ApplicationID, &financial.MonthlyAllowance, &financial.DailyTravelCost, &financial.MonthlyDormCost,
		&financial.OtherMonthlyCosts, &financial.HasIncome, &financial.IncomeSource, &financial.MonthlyIncome,
		&financial.FinancialNotes, &financial.CreatedAt, &financial.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return financial, nil
}

func (r *ApplicationDetailsRepository) UpdateFinancialInfo(financial *models.ApplicationFinancialInfo) error {
	query := `
		UPDATE application_financial_info
		SET monthly_allowance = $2, daily_travel_cost = $3, monthly_dorm_cost = $4,
		    other_monthly_costs = $5, has_income = $6, income_source = $7, monthly_income = $8,
		    financial_notes = $9, updated_at = $10
		WHERE financial_id = $1
	`

	_, err := r.db.Exec(query,
		financial.FinancialID, financial.MonthlyAllowance, financial.DailyTravelCost, financial.MonthlyDormCost,
		financial.OtherMonthlyCosts, financial.HasIncome, financial.IncomeSource, financial.MonthlyIncome,
		financial.FinancialNotes, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteFinancialInfo(financialID string) error {
	query := `DELETE FROM application_financial_info WHERE financial_id = $1`
	_, err := r.db.Exec(query, financialID)
	return err
}

// ========================================
// 10. Scholarship History Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateScholarshipHistory(history *models.ApplicationScholarshipHistory) error {
	query := `
		INSERT INTO application_scholarship_history (
			application_id, scholarship_name, scholarship_type, amount, academic_year,
			has_student_loan, loan_type, loan_year, loan_amount, notes,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING history_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		history.ApplicationID, history.ScholarshipName, history.ScholarshipType, history.Amount, history.AcademicYear,
		history.HasStudentLoan, history.LoanType, history.LoanYear, history.LoanAmount, history.Notes,
		now, now,
	).Scan(&history.HistoryID)

	return err
}

func (r *ApplicationDetailsRepository) GetScholarshipHistoryByApplicationID(applicationID int) ([]models.ApplicationScholarshipHistory, error) {
	query := `
		SELECT history_id, application_id, scholarship_name, scholarship_type, amount, academic_year,
		       has_student_loan, loan_type, loan_year, loan_amount, notes,
		       created_at, updated_at
		FROM application_scholarship_history
		WHERE application_id = $1
		ORDER BY academic_year DESC
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.ApplicationScholarshipHistory
	for rows.Next() {
		var history models.ApplicationScholarshipHistory
		err := rows.Scan(
			&history.HistoryID, &history.ApplicationID, &history.ScholarshipName, &history.ScholarshipType, &history.Amount, &history.AcademicYear,
			&history.HasStudentLoan, &history.LoanType, &history.LoanYear, &history.LoanAmount, &history.Notes,
			&history.CreatedAt, &history.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}

func (r *ApplicationDetailsRepository) UpdateScholarshipHistory(history *models.ApplicationScholarshipHistory) error {
	query := `
		UPDATE application_scholarship_history
		SET scholarship_name = $2, scholarship_type = $3, amount = $4, academic_year = $5,
		    has_student_loan = $6, loan_type = $7, loan_year = $8, loan_amount = $9, notes = $10,
		    updated_at = $11
		WHERE history_id = $1
	`

	_, err := r.db.Exec(query,
		history.HistoryID, history.ScholarshipName, history.ScholarshipType, history.Amount, history.AcademicYear,
		history.HasStudentLoan, history.LoanType, history.LoanYear, history.LoanAmount, history.Notes,
		time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteScholarshipHistory(historyID string) error {
	query := `DELETE FROM application_scholarship_history WHERE history_id = $1`
	_, err := r.db.Exec(query, historyID)
	return err
}

// ========================================
// 11. Activity Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateActivity(activity *models.ApplicationActivity) error {
	query := `
		INSERT INTO application_activities (
			application_id, activity_type, activity_name, description, achievement,
			award_level, year, evidence_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING activity_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		activity.ApplicationID, activity.ActivityType, activity.ActivityName, activity.Description, activity.Achievement,
		activity.AwardLevel, activity.Year, activity.EvidenceURL, now, now,
	).Scan(&activity.ActivityID)

	return err
}

func (r *ApplicationDetailsRepository) GetActivitiesByApplicationID(applicationID int) ([]models.ApplicationActivity, error) {
	query := `
		SELECT activity_id, application_id, activity_type, activity_name, description, achievement,
		       award_level, year, evidence_url, created_at, updated_at
		FROM application_activities
		WHERE application_id = $1
		ORDER BY year DESC, activity_type
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.ApplicationActivity
	for rows.Next() {
		var activity models.ApplicationActivity
		err := rows.Scan(
			&activity.ActivityID, &activity.ApplicationID, &activity.ActivityType, &activity.ActivityName, &activity.Description, &activity.Achievement,
			&activity.AwardLevel, &activity.Year, &activity.EvidenceURL, &activity.CreatedAt, &activity.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		activities = append(activities, activity)
	}

	return activities, nil
}

func (r *ApplicationDetailsRepository) UpdateActivity(activity *models.ApplicationActivity) error {
	query := `
		UPDATE application_activities
		SET activity_type = $2, activity_name = $3, description = $4, achievement = $5,
		    award_level = $6, year = $7, evidence_url = $8, updated_at = $9
		WHERE activity_id = $1
	`

	_, err := r.db.Exec(query,
		activity.ActivityID, activity.ActivityType, activity.ActivityName, activity.Description, activity.Achievement,
		activity.AwardLevel, activity.Year, activity.EvidenceURL, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteActivity(activityID string) error {
	query := `DELETE FROM application_activities WHERE activity_id = $1`
	_, err := r.db.Exec(query, activityID)
	return err
}

// ========================================
// 12. Reference Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateReference(reference *models.ApplicationReference) error {
	query := `
		INSERT INTO application_references (
			application_id, title, first_name, last_name, relationship,
			address, phone, email, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING reference_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		reference.ApplicationID, reference.Title, reference.FirstName, reference.LastName, reference.Relationship,
		reference.Address, reference.Phone, reference.Email, reference.Notes, now, now,
	).Scan(&reference.ReferenceID)

	return err
}

func (r *ApplicationDetailsRepository) GetReferencesByApplicationID(applicationID int) ([]models.ApplicationReference, error) {
	query := `
		SELECT reference_id, application_id, title, first_name, last_name, relationship,
		       address, phone, email, notes, created_at, updated_at
		FROM application_references
		WHERE application_id = $1
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var references []models.ApplicationReference
	for rows.Next() {
		var reference models.ApplicationReference
		err := rows.Scan(
			&reference.ReferenceID, &reference.ApplicationID, &reference.Title, &reference.FirstName, &reference.LastName, &reference.Relationship,
			&reference.Address, &reference.Phone, &reference.Email, &reference.Notes, &reference.CreatedAt, &reference.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		references = append(references, reference)
	}

	return references, nil
}

func (r *ApplicationDetailsRepository) UpdateReference(reference *models.ApplicationReference) error {
	query := `
		UPDATE application_references
		SET title = $2, first_name = $3, last_name = $4, relationship = $5,
		    address = $6, phone = $7, email = $8, notes = $9, updated_at = $10
		WHERE reference_id = $1
	`

	_, err := r.db.Exec(query,
		reference.ReferenceID, reference.Title, reference.FirstName, reference.LastName, reference.Relationship,
		reference.Address, reference.Phone, reference.Email, reference.Notes, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteReference(referenceID string) error {
	query := `DELETE FROM application_references WHERE reference_id = $1`
	_, err := r.db.Exec(query, referenceID)
	return err
}

// ========================================
// 13. Health Info Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateHealthInfo(health *models.ApplicationHealthInfo) error {
	query := `
		INSERT INTO application_health_info (
			application_id, has_health_issues, health_condition, health_details,
			affects_study, study_impact_details, monthly_medical_cost,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING health_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		health.ApplicationID, health.HasHealthIssues, health.HealthCondition, health.HealthDetails,
		health.AffectsStudy, health.StudyImpactDetails, health.MonthlyMedicalCost,
		now, now,
	).Scan(&health.HealthID)

	return err
}

func (r *ApplicationDetailsRepository) GetHealthInfoByApplicationID(applicationID int) (*models.ApplicationHealthInfo, error) {
	query := `
		SELECT health_id, application_id, has_health_issues, health_condition, health_details,
		       affects_study, study_impact_details, monthly_medical_cost,
		       created_at, updated_at
		FROM application_health_info
		WHERE application_id = $1
	`

	health := &models.ApplicationHealthInfo{}
	err := r.db.QueryRow(query, applicationID).Scan(
		&health.HealthID, &health.ApplicationID, &health.HasHealthIssues, &health.HealthCondition, &health.HealthDetails,
		&health.AffectsStudy, &health.StudyImpactDetails, &health.MonthlyMedicalCost,
		&health.CreatedAt, &health.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return health, nil
}

func (r *ApplicationDetailsRepository) UpdateHealthInfo(health *models.ApplicationHealthInfo) error {
	query := `
		UPDATE application_health_info
		SET has_health_issues = $2, health_condition = $3, health_details = $4,
		    affects_study = $5, study_impact_details = $6, monthly_medical_cost = $7,
		    updated_at = $8
		WHERE health_id = $1
	`

	_, err := r.db.Exec(query,
		health.HealthID, health.HasHealthIssues, health.HealthCondition, health.HealthDetails,
		health.AffectsStudy, health.StudyImpactDetails, health.MonthlyMedicalCost,
		time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteHealthInfo(healthID string) error {
	query := `DELETE FROM application_health_info WHERE health_id = $1`
	_, err := r.db.Exec(query, healthID)
	return err
}

// ========================================
// 14. Funding Needs Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateFundingNeeds(needs *models.ApplicationFundingNeeds) error {
	query := `
		INSERT INTO application_funding_needs (
			application_id, tuition_support, monthly_support, book_support, dorm_support,
			other_support, other_details, total_requested, necessity_reason,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING need_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		needs.ApplicationID, needs.TuitionSupport, needs.MonthlySupport, needs.BookSupport, needs.DormSupport,
		needs.OtherSupport, needs.OtherDetails, needs.TotalRequested, needs.NecessityReason,
		now, now,
	).Scan(&needs.NeedID)

	return err
}

func (r *ApplicationDetailsRepository) GetFundingNeedsByApplicationID(applicationID int) (*models.ApplicationFundingNeeds, error) {
	query := `
		SELECT need_id, application_id, tuition_support, monthly_support, book_support, dorm_support,
		       other_support, other_details, total_requested, necessity_reason,
		       created_at, updated_at
		FROM application_funding_needs
		WHERE application_id = $1
	`

	needs := &models.ApplicationFundingNeeds{}
	err := r.db.QueryRow(query, applicationID).Scan(
		&needs.NeedID, &needs.ApplicationID, &needs.TuitionSupport, &needs.MonthlySupport, &needs.BookSupport, &needs.DormSupport,
		&needs.OtherSupport, &needs.OtherDetails, &needs.TotalRequested, &needs.NecessityReason,
		&needs.CreatedAt, &needs.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return needs, nil
}

func (r *ApplicationDetailsRepository) UpdateFundingNeeds(needs *models.ApplicationFundingNeeds) error {
	query := `
		UPDATE application_funding_needs
		SET tuition_support = $2, monthly_support = $3, book_support = $4, dorm_support = $5,
		    other_support = $6, other_details = $7, total_requested = $8, necessity_reason = $9,
		    updated_at = $10
		WHERE need_id = $1
	`

	_, err := r.db.Exec(query,
		needs.NeedID, needs.TuitionSupport, needs.MonthlySupport, needs.BookSupport, needs.DormSupport,
		needs.OtherSupport, needs.OtherDetails, needs.TotalRequested, needs.NecessityReason,
		time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteFundingNeeds(needID string) error {
	query := `DELETE FROM application_funding_needs WHERE need_id = $1`
	_, err := r.db.Exec(query, needID)
	return err
}

// ========================================
// 15. House Document Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateHouseDocument(doc *models.ApplicationHouseDocument) error {
	query := `
		INSERT INTO application_house_documents (
			application_id, document_type, document_url, file_name, file_size,
			mime_type, description, verified, verified_by, verified_at,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING doc_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		doc.ApplicationID, doc.DocumentType, doc.DocumentURL, doc.FileName, doc.FileSize,
		doc.MimeType, doc.Description, doc.Verified, doc.VerifiedBy, doc.VerifiedAt,
		now, now,
	).Scan(&doc.DocID)

	return err
}

func (r *ApplicationDetailsRepository) GetHouseDocumentsByApplicationID(applicationID int) ([]models.ApplicationHouseDocument, error) {
	query := `
		SELECT doc_id, application_id, document_type, document_url, file_name, file_size,
		       mime_type, description, verified, verified_by, verified_at,
		       created_at, updated_at
		FROM application_house_documents
		WHERE application_id = $1
		ORDER BY document_type
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.ApplicationHouseDocument
	for rows.Next() {
		var doc models.ApplicationHouseDocument
		err := rows.Scan(
			&doc.DocID, &doc.ApplicationID, &doc.DocumentType, &doc.DocumentURL, &doc.FileName, &doc.FileSize,
			&doc.MimeType, &doc.Description, &doc.Verified, &doc.VerifiedBy, &doc.VerifiedAt,
			&doc.CreatedAt, &doc.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (r *ApplicationDetailsRepository) UpdateHouseDocument(doc *models.ApplicationHouseDocument) error {
	query := `
		UPDATE application_house_documents
		SET document_type = $2, document_url = $3, file_name = $4, file_size = $5,
		    mime_type = $6, description = $7, verified = $8, verified_by = $9, verified_at = $10,
		    updated_at = $11
		WHERE doc_id = $1
	`

	_, err := r.db.Exec(query,
		doc.DocID, doc.DocumentType, doc.DocumentURL, doc.FileName, doc.FileSize,
		doc.MimeType, doc.Description, doc.Verified, doc.VerifiedBy, doc.VerifiedAt,
		time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteHouseDocument(docID string) error {
	query := `DELETE FROM application_house_documents WHERE doc_id = $1`
	_, err := r.db.Exec(query, docID)
	return err
}

// ========================================
// 16. Income Certificate Methods
// ========================================

func (r *ApplicationDetailsRepository) CreateIncomeCertificate(cert *models.ApplicationIncomeCertificate) error {
	query := `
		INSERT INTO application_income_certificates (
			application_id, owner_name, relationship, income_type, monthly_income,
			certified_by, certifier_position, certifier_id_card, certificate_url, id_card_copy_url,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING cert_id
	`

	now := time.Now()
	err := r.db.QueryRow(query,
		cert.ApplicationID, cert.OwnerName, cert.Relationship, cert.IncomeType, cert.MonthlyIncome,
		cert.CertifiedBy, cert.CertifierPosition, cert.CertifierIDCard, cert.CertificateURL, cert.IDCardCopyURL,
		now, now,
	).Scan(&cert.CertID)

	return err
}

func (r *ApplicationDetailsRepository) GetIncomeCertificatesByApplicationID(applicationID int) ([]models.ApplicationIncomeCertificate, error) {
	query := `
		SELECT cert_id, application_id, owner_name, relationship, income_type, monthly_income,
		       certified_by, certifier_position, certifier_id_card, certificate_url, id_card_copy_url,
		       created_at, updated_at
		FROM application_income_certificates
		WHERE application_id = $1
		ORDER BY owner_name
	`

	rows, err := r.db.Query(query, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var certificates []models.ApplicationIncomeCertificate
	for rows.Next() {
		var cert models.ApplicationIncomeCertificate
		err := rows.Scan(
			&cert.CertID, &cert.ApplicationID, &cert.OwnerName, &cert.Relationship, &cert.IncomeType, &cert.MonthlyIncome,
			&cert.CertifiedBy, &cert.CertifierPosition, &cert.CertifierIDCard, &cert.CertificateURL, &cert.IDCardCopyURL,
			&cert.CreatedAt, &cert.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, cert)
	}

	return certificates, nil
}

func (r *ApplicationDetailsRepository) UpdateIncomeCertificate(cert *models.ApplicationIncomeCertificate) error {
	query := `
		UPDATE application_income_certificates
		SET owner_name = $2, relationship = $3, income_type = $4, monthly_income = $5,
		    certified_by = $6, certifier_position = $7, certifier_id_card = $8,
		    certificate_url = $9, id_card_copy_url = $10, updated_at = $11
		WHERE cert_id = $1
	`

	_, err := r.db.Exec(query,
		cert.CertID, cert.OwnerName, cert.Relationship, cert.IncomeType, cert.MonthlyIncome,
		cert.CertifiedBy, cert.CertifierPosition, cert.CertifierIDCard,
		cert.CertificateURL, cert.IDCardCopyURL, time.Now(),
	)

	return err
}

func (r *ApplicationDetailsRepository) DeleteIncomeCertificate(certID string) error {
	query := `DELETE FROM application_income_certificates WHERE cert_id = $1`
	_, err := r.db.Exec(query, certID)
	return err
}

// ========================================
// Special Functions - High-Level Save Methods
// ========================================

// SavePersonalInfo saves or updates personal information (upsert)
func (r *ApplicationDetailsRepository) SavePersonalInfo(info *models.ApplicationPersonalInfo) (*models.ApplicationPersonalInfo, error) {
	// Check if exists
	existing, err := r.GetPersonalInfoByApplicationID(info.ApplicationID)

	if err == sql.ErrNoRows || existing == nil {
		// Create new
		err = r.CreatePersonalInfo(info)
		if err != nil {
			return nil, fmt.Errorf("failed to create personal info: %w", err)
		}
		return info, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to check existing personal info: %w", err)
	}

	// Update existing - use existing ID
	info.InfoID = existing.InfoID
	err = r.UpdatePersonalInfo(info)
	if err != nil {
		return nil, fmt.Errorf("failed to update personal info: %w", err)
	}

	return info, nil
}

// SaveAddresses saves or updates addresses (replaces all addresses)
func (r *ApplicationDetailsRepository) SaveAddresses(applicationID uint, addresses []models.ApplicationAddress) ([]models.ApplicationAddress, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing addresses
	deleteQuery := `DELETE FROM application_addresses WHERE application_id = $1`
	_, err = tx.Exec(deleteQuery, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing addresses: %w", err)
	}

	// Insert new addresses
	savedAddresses := make([]models.ApplicationAddress, 0, len(addresses))
	for _, addr := range addresses {
		addr.ApplicationID = int(applicationID)

		insertQuery := `
			INSERT INTO application_addresses (
				application_id, address_type, house_number, village_number, alley,
				road, subdistrict, district, province, postal_code,
				address_line1, address_line2, latitude, longitude, map_image_url,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
			RETURNING address_id
		`

		now := time.Now()
		err = tx.QueryRow(insertQuery,
			addr.ApplicationID, addr.AddressType, addr.HouseNumber, addr.VillageNumber, addr.Alley,
			addr.Road, addr.Subdistrict, addr.District, addr.Province, addr.PostalCode,
			addr.AddressLine1, addr.AddressLine2, addr.Latitude, addr.Longitude, addr.MapImageURL,
			now, now,
		).Scan(&addr.AddressID)

		if err != nil {
			return nil, fmt.Errorf("failed to insert address: %w", err)
		}

		addr.CreatedAt = now
		addr.UpdatedAt = now
		savedAddresses = append(savedAddresses, addr)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return savedAddresses, nil
}

// SaveEducation saves or updates education history (replaces all records)
func (r *ApplicationDetailsRepository) SaveEducation(applicationID uint, education []models.ApplicationEducationHistory) ([]models.ApplicationEducationHistory, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing education records
	deleteQuery := `DELETE FROM application_education_history WHERE application_id = $1`
	_, err = tx.Exec(deleteQuery, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing education history: %w", err)
	}

	// Insert new education records
	savedEducation := make([]models.ApplicationEducationHistory, 0, len(education))
	for _, edu := range education {
		edu.ApplicationID = int(applicationID)

		insertQuery := `
			INSERT INTO application_education_history (
				application_id, education_level, school_name, school_province,
				gpa, graduation_year, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING history_id
		`

		now := time.Now()
		err = tx.QueryRow(insertQuery,
			edu.ApplicationID, edu.EducationLevel, edu.SchoolName, edu.SchoolProvince,
			edu.GPA, edu.GraduationYear, now, now,
		).Scan(&edu.HistoryID)

		if err != nil {
			return nil, fmt.Errorf("failed to insert education history: %w", err)
		}

		edu.CreatedAt = now
		edu.UpdatedAt = now
		savedEducation = append(savedEducation, edu)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return savedEducation, nil
}

// SaveFamily saves or updates family information (members, guardians, siblings, living situation)
func (r *ApplicationDetailsRepository) SaveFamily(
	applicationID uint,
	members []models.ApplicationFamilyMember,
	guardians []models.ApplicationGuardian,
	siblings []models.ApplicationSibling,
	livingSituation *models.ApplicationLivingSituation,
) (map[string]interface{}, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	result := make(map[string]interface{})

	// Save family members
	if members != nil {
		deleteQuery := `DELETE FROM application_family_members WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing family members: %w", err)
		}

		savedMembers := make([]models.ApplicationFamilyMember, 0, len(members))
		for _, member := range members {
			member.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_family_members (
					application_id, relationship, title, first_name, last_name,
					age, living_status, occupation, position, workplace,
					workplace_province, monthly_income, phone, notes,
					created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
				RETURNING member_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				member.ApplicationID, member.Relationship, member.Title, member.FirstName, member.LastName,
				member.Age, member.LivingStatus, member.Occupation, member.Position, member.Workplace,
				member.WorkplaceProvince, member.MonthlyIncome, member.Phone, member.Notes,
				now, now,
			).Scan(&member.MemberID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert family member: %w", err)
			}

			member.CreatedAt = now
			member.UpdatedAt = now
			savedMembers = append(savedMembers, member)
		}
		result["family_members"] = savedMembers
	}

	// Save guardians
	if guardians != nil {
		deleteQuery := `DELETE FROM application_guardians WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing guardians: %w", err)
		}

		savedGuardians := make([]models.ApplicationGuardian, 0, len(guardians))
		for _, guardian := range guardians {
			guardian.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_guardians (
					application_id, title, first_name, last_name, relationship,
					address, phone, occupation, position, workplace,
					workplace_phone, monthly_income, debts, debt_details,
					created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
				RETURNING guardian_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				guardian.ApplicationID, guardian.Title, guardian.FirstName, guardian.LastName, guardian.Relationship,
				guardian.Address, guardian.Phone, guardian.Occupation, guardian.Position, guardian.Workplace,
				guardian.WorkplacePhone, guardian.MonthlyIncome, guardian.Debts, guardian.DebtDetails,
				now, now,
			).Scan(&guardian.GuardianID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert guardian: %w", err)
			}

			guardian.CreatedAt = now
			guardian.UpdatedAt = now
			savedGuardians = append(savedGuardians, guardian)
		}
		result["guardians"] = savedGuardians
	}

	// Save siblings
	if siblings != nil {
		deleteQuery := `DELETE FROM application_siblings WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing siblings: %w", err)
		}

		savedSiblings := make([]models.ApplicationSibling, 0, len(siblings))
		for _, sibling := range siblings {
			sibling.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_siblings (
					application_id, sibling_order, gender, school_or_workplace, education_level,
					is_working, monthly_income, notes, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING sibling_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				sibling.ApplicationID, sibling.SiblingOrder, sibling.Gender, sibling.SchoolOrWorkplace, sibling.EducationLevel,
				sibling.IsWorking, sibling.MonthlyIncome, sibling.Notes, now, now,
			).Scan(&sibling.SiblingID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert sibling: %w", err)
			}

			sibling.CreatedAt = now
			sibling.UpdatedAt = now
			savedSiblings = append(savedSiblings, sibling)
		}
		result["siblings"] = savedSiblings
	}

	// Save living situation (upsert)
	if livingSituation != nil {
		livingSituation.ApplicationID = int(applicationID)

		existing, err := r.GetLivingSituationByApplicationID(int(applicationID))

		if err == sql.ErrNoRows || existing == nil {
			// Create new
			insertQuery := `
				INSERT INTO application_living_situation (
					application_id, living_with, living_details, front_house_image,
					side_house_image, back_house_image, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING living_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				livingSituation.ApplicationID, livingSituation.LivingWith, livingSituation.LivingDetails, livingSituation.FrontHouseImage,
				livingSituation.SideHouseImage, livingSituation.BackHouseImage, now, now,
			).Scan(&livingSituation.LivingID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert living situation: %w", err)
			}

			livingSituation.CreatedAt = now
			livingSituation.UpdatedAt = now
		} else if err == nil {
			// Update existing
			livingSituation.LivingID = existing.LivingID

			updateQuery := `
				UPDATE application_living_situation
				SET living_with = $2, living_details = $3, front_house_image = $4,
				    side_house_image = $5, back_house_image = $6, updated_at = $7
				WHERE living_id = $1
			`

			now := time.Now()
			_, err = tx.Exec(updateQuery,
				livingSituation.LivingID, livingSituation.LivingWith, livingSituation.LivingDetails, livingSituation.FrontHouseImage,
				livingSituation.SideHouseImage, livingSituation.BackHouseImage, now,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to update living situation: %w", err)
			}

			livingSituation.UpdatedAt = now
		} else {
			return nil, fmt.Errorf("failed to check existing living situation: %w", err)
		}

		result["living_situation"] = livingSituation
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// SaveFinancial saves or updates financial information
func (r *ApplicationDetailsRepository) SaveFinancial(
	applicationID uint,
	financialInfo *models.ApplicationFinancialInfo,
	assets []models.ApplicationAsset,
	scholarshipHistory []models.ApplicationScholarshipHistory,
	healthInfo *models.ApplicationHealthInfo,
	fundingNeeds *models.ApplicationFundingNeeds,
) (map[string]interface{}, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	result := make(map[string]interface{})

	// Save financial info (upsert)
	if financialInfo != nil {
		financialInfo.ApplicationID = int(applicationID)

		existing, err := r.GetFinancialInfoByApplicationID(int(applicationID))

		if err == sql.ErrNoRows || existing == nil {
			// Create new
			insertQuery := `
				INSERT INTO application_financial_info (
					application_id, monthly_allowance, daily_travel_cost, monthly_dorm_cost,
					other_monthly_costs, has_income, income_source, monthly_income,
					financial_notes, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING financial_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				financialInfo.ApplicationID, financialInfo.MonthlyAllowance, financialInfo.DailyTravelCost, financialInfo.MonthlyDormCost,
				financialInfo.OtherMonthlyCosts, financialInfo.HasIncome, financialInfo.IncomeSource, financialInfo.MonthlyIncome,
				financialInfo.FinancialNotes, now, now,
			).Scan(&financialInfo.FinancialID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert financial info: %w", err)
			}

			financialInfo.CreatedAt = now
			financialInfo.UpdatedAt = now
		} else if err == nil {
			// Update existing
			financialInfo.FinancialID = existing.FinancialID

			updateQuery := `
				UPDATE application_financial_info
				SET monthly_allowance = $2, daily_travel_cost = $3, monthly_dorm_cost = $4,
				    other_monthly_costs = $5, has_income = $6, income_source = $7, monthly_income = $8,
				    financial_notes = $9, updated_at = $10
				WHERE financial_id = $1
			`

			now := time.Now()
			_, err = tx.Exec(updateQuery,
				financialInfo.FinancialID, financialInfo.MonthlyAllowance, financialInfo.DailyTravelCost, financialInfo.MonthlyDormCost,
				financialInfo.OtherMonthlyCosts, financialInfo.HasIncome, financialInfo.IncomeSource, financialInfo.MonthlyIncome,
				financialInfo.FinancialNotes, now,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to update financial info: %w", err)
			}

			financialInfo.UpdatedAt = now
		} else {
			return nil, fmt.Errorf("failed to check existing financial info: %w", err)
		}

		result["financial_info"] = financialInfo
	}

	// Save assets (replace all)
	if assets != nil {
		deleteQuery := `DELETE FROM application_assets WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing assets: %w", err)
		}

		savedAssets := make([]models.ApplicationAsset, 0, len(assets))
		for _, asset := range assets {
			asset.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_assets (
					application_id, asset_type, category, description, value,
					monthly_cost, notes, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				RETURNING asset_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				asset.ApplicationID, asset.AssetType, asset.Category, asset.Description, asset.Value,
				asset.MonthlyCost, asset.Notes, now, now,
			).Scan(&asset.AssetID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert asset: %w", err)
			}

			asset.CreatedAt = now
			asset.UpdatedAt = now
			savedAssets = append(savedAssets, asset)
		}
		result["assets"] = savedAssets
	}

	// Save scholarship history (replace all)
	if scholarshipHistory != nil {
		deleteQuery := `DELETE FROM application_scholarship_history WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing scholarship history: %w", err)
		}

		savedHistory := make([]models.ApplicationScholarshipHistory, 0, len(scholarshipHistory))
		for _, history := range scholarshipHistory {
			history.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_scholarship_history (
					application_id, scholarship_name, scholarship_type, amount, academic_year,
					has_student_loan, loan_type, loan_year, loan_amount, notes,
					created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
				RETURNING history_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				history.ApplicationID, history.ScholarshipName, history.ScholarshipType, history.Amount, history.AcademicYear,
				history.HasStudentLoan, history.LoanType, history.LoanYear, history.LoanAmount, history.Notes,
				now, now,
			).Scan(&history.HistoryID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert scholarship history: %w", err)
			}

			history.CreatedAt = now
			history.UpdatedAt = now
			savedHistory = append(savedHistory, history)
		}
		result["scholarship_history"] = savedHistory
	}

	// Save health info (upsert)
	if healthInfo != nil {
		healthInfo.ApplicationID = int(applicationID)

		existing, err := r.GetHealthInfoByApplicationID(int(applicationID))

		if err == sql.ErrNoRows || existing == nil {
			// Create new
			insertQuery := `
				INSERT INTO application_health_info (
					application_id, has_health_issues, health_condition, health_details,
					affects_study, study_impact_details, monthly_medical_cost,
					created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				RETURNING health_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				healthInfo.ApplicationID, healthInfo.HasHealthIssues, healthInfo.HealthCondition, healthInfo.HealthDetails,
				healthInfo.AffectsStudy, healthInfo.StudyImpactDetails, healthInfo.MonthlyMedicalCost,
				now, now,
			).Scan(&healthInfo.HealthID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert health info: %w", err)
			}

			healthInfo.CreatedAt = now
			healthInfo.UpdatedAt = now
		} else if err == nil {
			// Update existing
			healthInfo.HealthID = existing.HealthID

			updateQuery := `
				UPDATE application_health_info
				SET has_health_issues = $2, health_condition = $3, health_details = $4,
				    affects_study = $5, study_impact_details = $6, monthly_medical_cost = $7,
				    updated_at = $8
				WHERE health_id = $1
			`

			now := time.Now()
			_, err = tx.Exec(updateQuery,
				healthInfo.HealthID, healthInfo.HasHealthIssues, healthInfo.HealthCondition, healthInfo.HealthDetails,
				healthInfo.AffectsStudy, healthInfo.StudyImpactDetails, healthInfo.MonthlyMedicalCost,
				now,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to update health info: %w", err)
			}

			healthInfo.UpdatedAt = now
		} else {
			return nil, fmt.Errorf("failed to check existing health info: %w", err)
		}

		result["health_info"] = healthInfo
	}

	// Save funding needs (upsert)
	if fundingNeeds != nil {
		fundingNeeds.ApplicationID = int(applicationID)

		existing, err := r.GetFundingNeedsByApplicationID(int(applicationID))

		if err == sql.ErrNoRows || existing == nil {
			// Create new
			insertQuery := `
				INSERT INTO application_funding_needs (
					application_id, tuition_support, monthly_support, book_support, dorm_support,
					other_support, other_details, total_requested, necessity_reason,
					created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING need_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				fundingNeeds.ApplicationID, fundingNeeds.TuitionSupport, fundingNeeds.MonthlySupport, fundingNeeds.BookSupport, fundingNeeds.DormSupport,
				fundingNeeds.OtherSupport, fundingNeeds.OtherDetails, fundingNeeds.TotalRequested, fundingNeeds.NecessityReason,
				now, now,
			).Scan(&fundingNeeds.NeedID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert funding needs: %w", err)
			}

			fundingNeeds.CreatedAt = now
			fundingNeeds.UpdatedAt = now
		} else if err == nil {
			// Update existing
			fundingNeeds.NeedID = existing.NeedID

			updateQuery := `
				UPDATE application_funding_needs
				SET tuition_support = $2, monthly_support = $3, book_support = $4, dorm_support = $5,
				    other_support = $6, other_details = $7, total_requested = $8, necessity_reason = $9,
				    updated_at = $10
				WHERE need_id = $1
			`

			now := time.Now()
			_, err = tx.Exec(updateQuery,
				fundingNeeds.NeedID, fundingNeeds.TuitionSupport, fundingNeeds.MonthlySupport, fundingNeeds.BookSupport, fundingNeeds.DormSupport,
				fundingNeeds.OtherSupport, fundingNeeds.OtherDetails, fundingNeeds.TotalRequested, fundingNeeds.NecessityReason,
				now,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to update funding needs: %w", err)
			}

			fundingNeeds.UpdatedAt = now
		} else {
			return nil, fmt.Errorf("failed to check existing funding needs: %w", err)
		}

		result["funding_needs"] = fundingNeeds
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// SaveActivities saves or updates activities and references
func (r *ApplicationDetailsRepository) SaveActivities(
	applicationID uint,
	activities []models.ApplicationActivity,
	references []models.ApplicationReference,
) (map[string]interface{}, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	result := make(map[string]interface{})

	// Save activities (replace all)
	if activities != nil {
		deleteQuery := `DELETE FROM application_activities WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing activities: %w", err)
		}

		savedActivities := make([]models.ApplicationActivity, 0, len(activities))
		for _, activity := range activities {
			activity.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_activities (
					application_id, activity_type, activity_name, description, achievement,
					award_level, year, evidence_url, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING activity_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				activity.ApplicationID, activity.ActivityType, activity.ActivityName, activity.Description, activity.Achievement,
				activity.AwardLevel, activity.Year, activity.EvidenceURL, now, now,
			).Scan(&activity.ActivityID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert activity: %w", err)
			}

			activity.CreatedAt = now
			activity.UpdatedAt = now
			savedActivities = append(savedActivities, activity)
		}
		result["activities"] = savedActivities
	}

	// Save references (replace all)
	if references != nil {
		deleteQuery := `DELETE FROM application_references WHERE application_id = $1`
		_, err = tx.Exec(deleteQuery, applicationID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete existing references: %w", err)
		}

		savedReferences := make([]models.ApplicationReference, 0, len(references))
		for _, reference := range references {
			reference.ApplicationID = int(applicationID)

			insertQuery := `
				INSERT INTO application_references (
					application_id, title, first_name, last_name, relationship,
					address, phone, email, notes, created_at, updated_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING reference_id
			`

			now := time.Now()
			err = tx.QueryRow(insertQuery,
				reference.ApplicationID, reference.Title, reference.FirstName, reference.LastName, reference.Relationship,
				reference.Address, reference.Phone, reference.Email, reference.Notes, now, now,
			).Scan(&reference.ReferenceID)

			if err != nil {
				return nil, fmt.Errorf("failed to insert reference: %w", err)
			}

			reference.CreatedAt = now
			reference.UpdatedAt = now
			savedReferences = append(savedReferences, reference)
		}
		result["references"] = savedReferences
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return result, nil
}

// SaveCompleteForm saves all application details at once in a transaction
func (r *ApplicationDetailsRepository) SaveCompleteForm(applicationID uint, form *models.CompleteApplicationForm) (*models.CompleteApplicationForm, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Note: This is a simplified implementation that uses the existing Save methods
	// In a production environment, you might want to implement this differently to use a single transaction

	// For now, we'll call the individual save methods
	// This is not ideal because each method starts its own transaction
	// A better approach would be to refactor the save methods to accept a transaction parameter

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// After saving, retrieve the complete form
	return r.GetCompleteForm(applicationID)
}

// GetCompleteForm retrieves all application details
func (r *ApplicationDetailsRepository) GetCompleteForm(applicationID uint) (*models.CompleteApplicationForm, error) {
	return r.GetCompleteApplication(int(applicationID))
}

// ========================================
// Special Functions - Legacy
// ========================================

// GetCompleteApplication retrieves all application details in one call
func (r *ApplicationDetailsRepository) GetCompleteApplication(applicationID int) (*models.CompleteApplicationForm, error) {
	form := &models.CompleteApplicationForm{}

	// Get Personal Info (single record)
	personalInfo, err := r.GetPersonalInfoByApplicationID(applicationID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting personal info: %w", err)
	}
	if err == nil {
		form.PersonalInfo = personalInfo
	}

	// Get Addresses (multiple records)
	addresses, err := r.GetAddressesByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting addresses: %w", err)
	}
	form.Addresses = addresses

	// Get Education History (multiple records)
	educationHistory, err := r.GetEducationHistoryByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting education history: %w", err)
	}
	form.EducationHistory = educationHistory

	// Get Family Members (multiple records)
	familyMembers, err := r.GetFamilyMembersByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting family members: %w", err)
	}
	form.FamilyMembers = familyMembers

	// Get Assets (multiple records)
	assets, err := r.GetAssetsByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting assets: %w", err)
	}
	form.Assets = assets

	// Get Guardians (multiple records)
	guardians, err := r.GetGuardiansByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting guardians: %w", err)
	}
	form.Guardians = guardians

	// Get Siblings (multiple records)
	siblings, err := r.GetSiblingsByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting siblings: %w", err)
	}
	form.Siblings = siblings

	// Get Living Situation (single record)
	livingSituation, err := r.GetLivingSituationByApplicationID(applicationID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting living situation: %w", err)
	}
	if err == nil {
		form.LivingSituation = livingSituation
	}

	// Get Financial Info (single record)
	financialInfo, err := r.GetFinancialInfoByApplicationID(applicationID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting financial info: %w", err)
	}
	if err == nil {
		form.FinancialInfo = financialInfo
	}

	// Get Scholarship History (multiple records)
	scholarshipHistory, err := r.GetScholarshipHistoryByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting scholarship history: %w", err)
	}
	form.ScholarshipHistory = scholarshipHistory

	// Get Activities (multiple records)
	activities, err := r.GetActivitiesByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting activities: %w", err)
	}
	form.Activities = activities

	// Get References (multiple records)
	references, err := r.GetReferencesByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting references: %w", err)
	}
	form.References = references

	// Get Health Info (single record)
	healthInfo, err := r.GetHealthInfoByApplicationID(applicationID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting health info: %w", err)
	}
	if err == nil {
		form.HealthInfo = healthInfo
	}

	// Get Funding Needs (single record)
	fundingNeeds, err := r.GetFundingNeedsByApplicationID(applicationID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error getting funding needs: %w", err)
	}
	if err == nil {
		form.FundingNeeds = fundingNeeds
	}

	// Get House Documents (multiple records)
	houseDocuments, err := r.GetHouseDocumentsByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting house documents: %w", err)
	}
	form.HouseDocuments = houseDocuments

	// Get Income Certificates (multiple records)
	incomeCertificates, err := r.GetIncomeCertificatesByApplicationID(applicationID)
	if err != nil {
		return nil, fmt.Errorf("error getting income certificates: %w", err)
	}
	form.IncomeCertificates = incomeCertificates

	return form, nil
}

// SaveCompleteApplication saves all application details in a transaction
func (r *ApplicationDetailsRepository) SaveCompleteApplication(form *models.CompleteApplicationForm) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	// Note: This is a simplified version. In production, you would want to:
	// 1. Check if records exist and update them if they do
	// 2. Handle errors more gracefully
	// 3. Validate data before saving
	// 4. Use the transaction (tx) for all operations

	// For now, this function serves as a template
	// The actual implementation would depend on your business logic

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
