package handler

import (
	"context"
	"time"

	"github.com/dwikie/sentra-payment-orchestrator/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationHandler struct {
	Pool *pgxpool.Pool
}

func NewOrganizationHandler(pool *pgxpool.Pool) *OrganizationHandler {
	return &OrganizationHandler{Pool: pool}
}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userId := c.GetInt64("user_id")

	var payload model.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(ctx)

	if err := h.createOrganization(userId, tx, ctx, payload); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create organization"})
		return
	}

	if err := tx.Commit(ctx); err != nil {
		c.JSON(500, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(201, gin.H{"message": "Organization created successfully"})
}

func (h *OrganizationHandler) createOrganization(userId int64, tx pgx.Tx, ctx context.Context, payload model.CreateOrganizationRequest) error {
	var orgId int64
	err := tx.QueryRow(ctx, `
		INSERT INTO organizations (org_name, org_type, org_email, org_phone, org_tax_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, payload.OrgName, payload.OrgType, payload.OrgEmail, payload.OrgPhone, payload.OrgTaxId, userId).Scan(&orgId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO organization_addresses (org_id, org_address, org_city, org_state, org_zip_code, org_country)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, orgId, payload.OrgAddress, payload.OrgCity, payload.OrgState, payload.OrgZipCode, payload.OrgCountry)
	if err != nil {
		return err
	}

	return nil
}
