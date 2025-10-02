package models

import (
	"database/sql"
	"time"
)

// ApplicationPersonalInfo represents personal information of applicant
type ApplicationPersonalInfo struct {
	InfoID         string         `json:"info_id" db:"info_id"`
	ApplicationID  int            `json:"application_id" db:"application_id"`
	PrefixTH       sql.NullString `json:"prefix_th,omitempty" db:"prefix_th"`
	PrefixEN       sql.NullString `json:"prefix_en,omitempty" db:"prefix_en"`
	FirstNameTH    string         `json:"first_name_th" db:"first_name_th"`
	LastNameTH     string         `json:"last_name_th" db:"last_name_th"`
	FirstNameEN    sql.NullString `json:"first_name_en,omitempty" db:"first_name_en"`
	LastNameEN     sql.NullString `json:"last_name_en,omitempty" db:"last_name_en"`
	Email          string         `json:"email" db:"email"`
	Phone          sql.NullString `json:"phone,omitempty" db:"phone"`
	LineID         sql.NullString `json:"line_id,omitempty" db:"line_id"`
	CitizenID      sql.NullString `json:"citizen_id,omitempty" db:"citizen_id"`
	StudentID      sql.NullString `json:"student_id,omitempty" db:"student_id"`
	Faculty        sql.NullString `json:"faculty,omitempty" db:"faculty"`
	Department     sql.NullString `json:"department,omitempty" db:"department"`
	Major          sql.NullString `json:"major,omitempty" db:"major"`
	YearLevel      sql.NullInt32  `json:"year_level,omitempty" db:"year_level"`
	AdmissionType  sql.NullString `json:"admission_type,omitempty" db:"admission_type"`
	AdmissionDetails sql.NullString `json:"admission_details,omitempty" db:"admission_details"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationAddress represents addresses
type ApplicationAddress struct {
	AddressID     string         `json:"address_id" db:"address_id"`
	ApplicationID int            `json:"application_id" db:"application_id"`
	AddressType   string         `json:"address_type" db:"address_type"`
	HouseNumber   sql.NullString `json:"house_number,omitempty" db:"house_number"`
	VillageNumber sql.NullString `json:"village_number,omitempty" db:"village_number"`
	Alley         sql.NullString `json:"alley,omitempty" db:"alley"`
	Road          sql.NullString `json:"road,omitempty" db:"road"`
	Subdistrict   sql.NullString `json:"subdistrict,omitempty" db:"subdistrict"`
	District      sql.NullString `json:"district,omitempty" db:"district"`
	Province      sql.NullString `json:"province,omitempty" db:"province"`
	PostalCode    sql.NullString `json:"postal_code,omitempty" db:"postal_code"`
	AddressLine1  sql.NullString `json:"address_line1,omitempty" db:"address_line1"`
	AddressLine2  sql.NullString `json:"address_line2,omitempty" db:"address_line2"`
	Latitude      sql.NullFloat64 `json:"latitude,omitempty" db:"latitude"`
	Longitude     sql.NullFloat64 `json:"longitude,omitempty" db:"longitude"`
	MapImageURL   sql.NullString `json:"map_image_url,omitempty" db:"map_image_url"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationEducationHistory represents education history
type ApplicationEducationHistory struct {
	HistoryID       string         `json:"history_id" db:"history_id"`
	ApplicationID   int            `json:"application_id" db:"application_id"`
	EducationLevel  string         `json:"education_level" db:"education_level"`
	SchoolName      string         `json:"school_name" db:"school_name"`
	SchoolProvince  sql.NullString `json:"school_province,omitempty" db:"school_province"`
	GPA             sql.NullFloat64 `json:"gpa,omitempty" db:"gpa"`
	GraduationYear  sql.NullString `json:"graduation_year,omitempty" db:"graduation_year"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationFamilyMember represents family member information
type ApplicationFamilyMember struct {
	MemberID         string         `json:"member_id" db:"member_id"`
	ApplicationID    int            `json:"application_id" db:"application_id"`
	Relationship     string         `json:"relationship" db:"relationship"`
	Title            sql.NullString `json:"title,omitempty" db:"title"`
	FirstName        string         `json:"first_name" db:"first_name"`
	LastName         string         `json:"last_name" db:"last_name"`
	Age              sql.NullInt32  `json:"age,omitempty" db:"age"`
	LivingStatus     sql.NullString `json:"living_status,omitempty" db:"living_status"`
	Occupation       sql.NullString `json:"occupation,omitempty" db:"occupation"`
	Position         sql.NullString `json:"position,omitempty" db:"position"`
	Workplace        sql.NullString `json:"workplace,omitempty" db:"workplace"`
	WorkplaceProvince sql.NullString `json:"workplace_province,omitempty" db:"workplace_province"`
	MonthlyIncome    sql.NullFloat64 `json:"monthly_income,omitempty" db:"monthly_income"`
	Phone            sql.NullString `json:"phone,omitempty" db:"phone"`
	Notes            sql.NullString `json:"notes,omitempty" db:"notes"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationAsset represents assets and liabilities
type ApplicationAsset struct {
	AssetID       string         `json:"asset_id" db:"asset_id"`
	ApplicationID int            `json:"application_id" db:"application_id"`
	AssetType     string         `json:"asset_type" db:"asset_type"`
	Category      sql.NullString `json:"category,omitempty" db:"category"`
	Description   sql.NullString `json:"description,omitempty" db:"description"`
	Value         sql.NullFloat64 `json:"value,omitempty" db:"value"`
	MonthlyCost   sql.NullFloat64 `json:"monthly_cost,omitempty" db:"monthly_cost"`
	Notes         sql.NullString `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationGuardian represents guardians (not parents)
type ApplicationGuardian struct {
	GuardianID      string         `json:"guardian_id" db:"guardian_id"`
	ApplicationID   int            `json:"application_id" db:"application_id"`
	Title           sql.NullString `json:"title,omitempty" db:"title"`
	FirstName       string         `json:"first_name" db:"first_name"`
	LastName        string         `json:"last_name" db:"last_name"`
	Relationship    sql.NullString `json:"relationship,omitempty" db:"relationship"`
	Address         sql.NullString `json:"address,omitempty" db:"address"`
	Phone           sql.NullString `json:"phone,omitempty" db:"phone"`
	Occupation      sql.NullString `json:"occupation,omitempty" db:"occupation"`
	Position        sql.NullString `json:"position,omitempty" db:"position"`
	Workplace       sql.NullString `json:"workplace,omitempty" db:"workplace"`
	WorkplacePhone  sql.NullString `json:"workplace_phone,omitempty" db:"workplace_phone"`
	MonthlyIncome   sql.NullFloat64 `json:"monthly_income,omitempty" db:"monthly_income"`
	Debts           sql.NullFloat64 `json:"debts,omitempty" db:"debts"`
	DebtDetails     sql.NullString `json:"debt_details,omitempty" db:"debt_details"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationSibling represents siblings information
type ApplicationSibling struct {
	SiblingID          string         `json:"sibling_id" db:"sibling_id"`
	ApplicationID      int            `json:"application_id" db:"application_id"`
	SiblingOrder       int            `json:"sibling_order" db:"sibling_order"`
	Gender             sql.NullString `json:"gender,omitempty" db:"gender"`
	SchoolOrWorkplace  sql.NullString `json:"school_or_workplace,omitempty" db:"school_or_workplace"`
	EducationLevel     sql.NullString `json:"education_level,omitempty" db:"education_level"`
	IsWorking          bool           `json:"is_working" db:"is_working"`
	MonthlyIncome      sql.NullFloat64 `json:"monthly_income,omitempty" db:"monthly_income"`
	Notes              sql.NullString `json:"notes,omitempty" db:"notes"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationLivingSituation represents living situation
type ApplicationLivingSituation struct {
	LivingID        string         `json:"living_id" db:"living_id"`
	ApplicationID   int            `json:"application_id" db:"application_id"`
	LivingWith      string         `json:"living_with" db:"living_with"`
	LivingDetails   sql.NullString `json:"living_details,omitempty" db:"living_details"`
	FrontHouseImage sql.NullString `json:"front_house_image,omitempty" db:"front_house_image"`
	SideHouseImage  sql.NullString `json:"side_house_image,omitempty" db:"side_house_image"`
	BackHouseImage  sql.NullString `json:"back_house_image,omitempty" db:"back_house_image"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationFinancialInfo represents financial information
type ApplicationFinancialInfo struct {
	FinancialID      string         `json:"financial_id" db:"financial_id"`
	ApplicationID    int            `json:"application_id" db:"application_id"`
	MonthlyAllowance sql.NullFloat64 `json:"monthly_allowance,omitempty" db:"monthly_allowance"`
	DailyTravelCost  sql.NullFloat64 `json:"daily_travel_cost,omitempty" db:"daily_travel_cost"`
	MonthlyDormCost  sql.NullFloat64 `json:"monthly_dorm_cost,omitempty" db:"monthly_dorm_cost"`
	OtherMonthlyCosts sql.NullFloat64 `json:"other_monthly_costs,omitempty" db:"other_monthly_costs"`
	HasIncome        bool           `json:"has_income" db:"has_income"`
	IncomeSource     sql.NullString `json:"income_source,omitempty" db:"income_source"`
	MonthlyIncome    sql.NullFloat64 `json:"monthly_income,omitempty" db:"monthly_income"`
	FinancialNotes   sql.NullString `json:"financial_notes,omitempty" db:"financial_notes"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationScholarshipHistory represents scholarship history
type ApplicationScholarshipHistory struct {
	HistoryID        string         `json:"history_id" db:"history_id"`
	ApplicationID    int            `json:"application_id" db:"application_id"`
	ScholarshipName  sql.NullString `json:"scholarship_name,omitempty" db:"scholarship_name"`
	ScholarshipType  sql.NullString `json:"scholarship_type,omitempty" db:"scholarship_type"`
	Amount           sql.NullFloat64 `json:"amount,omitempty" db:"amount"`
	AcademicYear     sql.NullString `json:"academic_year,omitempty" db:"academic_year"`
	HasStudentLoan   bool           `json:"has_student_loan" db:"has_student_loan"`
	LoanType         sql.NullString `json:"loan_type,omitempty" db:"loan_type"`
	LoanYear         sql.NullString `json:"loan_year,omitempty" db:"loan_year"`
	LoanAmount       sql.NullFloat64 `json:"loan_amount,omitempty" db:"loan_amount"`
	Notes            sql.NullString `json:"notes,omitempty" db:"notes"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationActivity represents activities and special abilities
type ApplicationActivity struct {
	ActivityID    string         `json:"activity_id" db:"activity_id"`
	ApplicationID int            `json:"application_id" db:"application_id"`
	ActivityType  string         `json:"activity_type" db:"activity_type"`
	ActivityName  sql.NullString `json:"activity_name,omitempty" db:"activity_name"`
	Description   sql.NullString `json:"description,omitempty" db:"description"`
	Achievement   sql.NullString `json:"achievement,omitempty" db:"achievement"`
	AwardLevel    sql.NullString `json:"award_level,omitempty" db:"award_level"`
	Year          sql.NullString `json:"year,omitempty" db:"year"`
	EvidenceURL   sql.NullString `json:"evidence_url,omitempty" db:"evidence_url"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationReference represents references
type ApplicationReference struct {
	ReferenceID   string         `json:"reference_id" db:"reference_id"`
	ApplicationID int            `json:"application_id" db:"application_id"`
	Title         sql.NullString `json:"title,omitempty" db:"title"`
	FirstName     string         `json:"first_name" db:"first_name"`
	LastName      string         `json:"last_name" db:"last_name"`
	Relationship  sql.NullString `json:"relationship,omitempty" db:"relationship"`
	Address       sql.NullString `json:"address,omitempty" db:"address"`
	Phone         sql.NullString `json:"phone,omitempty" db:"phone"`
	Email         sql.NullString `json:"email,omitempty" db:"email"`
	Notes         sql.NullString `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationHealthInfo represents health information
type ApplicationHealthInfo struct {
	HealthID           string         `json:"health_id" db:"health_id"`
	ApplicationID      int            `json:"application_id" db:"application_id"`
	HasHealthIssues    bool           `json:"has_health_issues" db:"has_health_issues"`
	HealthCondition    sql.NullString `json:"health_condition,omitempty" db:"health_condition"`
	HealthDetails      sql.NullString `json:"health_details,omitempty" db:"health_details"`
	AffectsStudy       bool           `json:"affects_study" db:"affects_study"`
	StudyImpactDetails sql.NullString `json:"study_impact_details,omitempty" db:"study_impact_details"`
	MonthlyMedicalCost sql.NullFloat64 `json:"monthly_medical_cost,omitempty" db:"monthly_medical_cost"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationFundingNeeds represents funding needs
type ApplicationFundingNeeds struct {
	NeedID           string         `json:"need_id" db:"need_id"`
	ApplicationID    int            `json:"application_id" db:"application_id"`
	TuitionSupport   sql.NullFloat64 `json:"tuition_support,omitempty" db:"tuition_support"`
	MonthlySupport   sql.NullFloat64 `json:"monthly_support,omitempty" db:"monthly_support"`
	BookSupport      sql.NullFloat64 `json:"book_support,omitempty" db:"book_support"`
	DormSupport      sql.NullFloat64 `json:"dorm_support,omitempty" db:"dorm_support"`
	OtherSupport     sql.NullFloat64 `json:"other_support,omitempty" db:"other_support"`
	OtherDetails     sql.NullString `json:"other_details,omitempty" db:"other_details"`
	TotalRequested   sql.NullFloat64 `json:"total_requested,omitempty" db:"total_requested"`
	NecessityReason  sql.NullString `json:"necessity_reason,omitempty" db:"necessity_reason"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationHouseDocument represents house documents
type ApplicationHouseDocument struct {
	DocID        string         `json:"doc_id" db:"doc_id"`
	ApplicationID int           `json:"application_id" db:"application_id"`
	DocumentType string         `json:"document_type" db:"document_type"`
	DocumentURL  string         `json:"document_url" db:"document_url"`
	FileName     sql.NullString `json:"file_name,omitempty" db:"file_name"`
	FileSize     sql.NullInt32  `json:"file_size,omitempty" db:"file_size"`
	MimeType     sql.NullString `json:"mime_type,omitempty" db:"mime_type"`
	Description  sql.NullString `json:"description,omitempty" db:"description"`
	Verified     bool           `json:"verified" db:"verified"`
	VerifiedBy   sql.NullString `json:"verified_by,omitempty" db:"verified_by"`
	VerifiedAt   sql.NullTime   `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}

// ApplicationIncomeCertificate represents income certificates
type ApplicationIncomeCertificate struct {
	CertID          string         `json:"cert_id" db:"cert_id"`
	ApplicationID   int            `json:"application_id" db:"application_id"`
	OwnerName       string         `json:"owner_name" db:"owner_name"`
	Relationship    sql.NullString `json:"relationship,omitempty" db:"relationship"`
	IncomeType      sql.NullString `json:"income_type,omitempty" db:"income_type"`
	MonthlyIncome   sql.NullFloat64 `json:"monthly_income,omitempty" db:"monthly_income"`
	CertifiedBy     sql.NullString `json:"certified_by,omitempty" db:"certified_by"`
	CertifierPosition sql.NullString `json:"certifier_position,omitempty" db:"certifier_position"`
	CertifierIDCard sql.NullString `json:"certifier_id_card,omitempty" db:"certifier_id_card"`
	CertificateURL  sql.NullString `json:"certificate_url,omitempty" db:"certificate_url"`
	IDCardCopyURL   sql.NullString `json:"id_card_copy_url,omitempty" db:"id_card_copy_url"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
}

// CompleteApplicationForm represents the complete application with all details
type CompleteApplicationForm struct {
	Application         *ScholarshipApplication           `json:"application"`
	PersonalInfo        *ApplicationPersonalInfo          `json:"personal_info,omitempty"`
	Addresses           []ApplicationAddress              `json:"addresses,omitempty"`
	EducationHistory    []ApplicationEducationHistory     `json:"education_history,omitempty"`
	FamilyMembers       []ApplicationFamilyMember         `json:"family_members,omitempty"`
	Assets              []ApplicationAsset                `json:"assets,omitempty"`
	Guardians           []ApplicationGuardian             `json:"guardians,omitempty"`
	Siblings            []ApplicationSibling              `json:"siblings,omitempty"`
	LivingSituation     *ApplicationLivingSituation       `json:"living_situation,omitempty"`
	FinancialInfo       *ApplicationFinancialInfo         `json:"financial_info,omitempty"`
	ScholarshipHistory  []ApplicationScholarshipHistory   `json:"scholarship_history,omitempty"`
	Activities          []ApplicationActivity             `json:"activities,omitempty"`
	References          []ApplicationReference            `json:"references,omitempty"`
	HealthInfo          *ApplicationHealthInfo            `json:"health_info,omitempty"`
	FundingNeeds        *ApplicationFundingNeeds          `json:"funding_needs,omitempty"`
	HouseDocuments      []ApplicationHouseDocument        `json:"house_documents,omitempty"`
	IncomeCertificates  []ApplicationIncomeCertificate    `json:"income_certificates,omitempty"`
}
