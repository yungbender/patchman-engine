package controllers

import (
	"app/base/database"
	"app/base/models"
	"app/base/utils"
	"app/manager/middlewares"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var BaselineSystemFields = database.MustGetQueryAttrs(&BaselineSystemsDBLookup{})
var BaselineSystemSelect = database.MustGetSelect(&BaselineSystemsDBLookup{})
var BaselineSystemOpts = ListOpts{
	Fields:         BaselineSystemFields,
	DefaultFilters: map[string]FilterData{},
	DefaultSort:    "-display_name",
	SearchFields:   []string{"sp.display_name"},
	TotalFunc:      CountRows,
}

type BaselineSystemsDBLookup struct {
	ID string `json:"id" csv:"id" query:"sp.inventory_id" gorm:"column:id"`
	BaselineSystemAttributes
}

type BaselineSystemAttributes struct {
	DisplayName string `json:"display_name" csv:"display_name" query:"sp.display_name" gorm:"column:display_name"`
}

type BaselineSystemItem struct {
	Attributes BaselineSystemAttributes `json:"attributes"`
	ID         string                   `json:"id"`
	Type       string                   `json:"type"`
}

type BaselineSystemInlineItem BaselineSystemsDBLookup

type BaselineSystemsResponse struct {
	Data  []BaselineSystemItem `json:"data"`
	Links Links                `json:"links"`
	Meta  ListMeta             `json:"meta"`
}

// @Summary Show me all systems belonging to a baseline
// @Description  Show me all systems applicable to a baseline
// @ID listBaselineSystems
// @Security RhIdentity
// @Accept   json
// @Produce  json
// @Param    baseline_id    path    int     true    "Baseline ID"
// @Param    limit          query   int     false   "Limit for paging, set -1 to return all"
// @Param    offset         query   int     false   "Offset for paging"
// @Param    sort           query   string  false   "Sort field"    Enums(id,name,config)
// @Param    search         query   string  false   "Find matching text"
// @Param    filter[display_name]           query   string  false "Filter"
// @Success 200 {object} BaselineSystemsResponse
// @Router /api/patch/v1/baselines/{baseline_id}/systems [get]
func BaselineSystemsListHandler(c *gin.Context) {
	account := c.GetInt(middlewares.KeyAccount)

	baselineID := c.Param("baseline_id")
	if baselineID == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse{Error: "baseline_id param not found"})
		return
	}

	var exists int64
	err := database.Db.Model(&models.Baseline{}).
		Where("id = ? ", baselineID).Count(&exists).Error
	if err != nil {
		LogAndRespError(c, err, "database error")
	}
	if exists == 0 {
		LogAndRespNotFound(c, errors.New("Baseline not found"), "Baseline not found")
		return
	}

	query := buildQueryBaselineSystems(account, baselineID)

	path := fmt.Sprintf("/api/patch/v1/baselines/%v/systems", baselineID)
	query, meta, links, err := ListCommon(query, c, path, BaselineSystemOpts)
	if err != nil {
		// Error handling and setting of result code & content is done in ListCommon
		return
	}

	var baselineSystems []BaselineSystemsDBLookup
	err = query.Find(&baselineSystems).Error
	if err != nil {
		LogAndRespError(c, err, err.Error())
	}

	data := buildBaselineSystemData(baselineSystems)
	var resp = BaselineSystemsResponse{
		Data:  data,
		Links: *links,
		Meta:  *meta,
	}
	c.JSON(http.StatusOK, &resp)
}

func buildQueryBaselineSystems(account int, baselineID string) *gorm.DB {
	query := database.Db.Table("system_platform AS sp").Select(BaselineSystemSelect).
		Joins("JOIN inventory.hosts ih ON ih.id = sp.inventory_id").
		Where("sp.rh_account_id = ? AND sp.baseline_id = ?", account, baselineID).
		Where("sp.stale = false")
	return query
}

func buildBaselineSystemData(baselineSystems []BaselineSystemsDBLookup) []BaselineSystemItem {
	data := make([]BaselineSystemItem, len(baselineSystems))
	for i := 0; i < len(baselineSystems); i++ {
		baselineSystemDB := baselineSystems[i]
		data[i] = BaselineSystemItem{
			Attributes: BaselineSystemAttributes{
				DisplayName: baselineSystemDB.BaselineSystemAttributes.DisplayName,
			},
			ID:   baselineSystemDB.ID,
			Type: "baseline_systems",
		}
	}
	return data
}