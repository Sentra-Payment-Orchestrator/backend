package model

type CreateOrganizationRequest struct {
	OrgName    string `json:"org_name" binding:"required"`
	OrgType    string `json:"org_type" binding:"required"`
	OrgEmail   string `json:"org_email" binding:"required,email"`
	OrgPhone   string `json:"org_phone" binding:"required"`
	OrgTaxId   string `json:"org_tax_id" binding:"required"`
	OrgAddress string `json:"org_address" binding:"required"`
	OrgCity    string `json:"org_city" binding:"required"`
	OrgState   string `json:"org_state" binding:"required"`
	OrgZipCode string `json:"org_zip_code" binding:"required"`
	OrgCountry string `json:"org_country" binding:"required"`
}
