package rating

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/fact-check/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/meili"
	"github.com/factly/dega-server/util/slug"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create rating
// @Summary Create rating
// @Description Create rating
// @Tags Rating
// @ID add-rating
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param Rating body rating true "Rating Object"
// @Success 201 {object} model.Rating
// @Failure 400 {array} string
// @Router /fact-check/ratings [post]
func create(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	rating := &rating{}

	err = json.NewDecoder(r.Body).Decode(&rating)

	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(rating)

	if validationError != nil {
		loggerx.Error(errors.New("validation error"))
		errorx.Render(w, validationError)
		return
	}

	var ratingSlug string
	if rating.Slug != "" && slug.Check(rating.Slug) {
		ratingSlug = rating.Slug
	} else {
		ratingSlug = slug.Make(rating.Name)
	}

	// Check if rating with same name exist
	newRatingName := strings.ToLower(strings.TrimSpace(rating.Name))
	var ratingCount int
	config.DB.Model(&model.Rating{}).Where(&model.Rating{
		SpaceID: uint(sID),
	}).Where("name ILIKE ?", newRatingName).Count(&ratingCount)

	if ratingCount > 0 {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.CannotSaveChanges()))
		return
	}

	result := &model.Rating{
		Name:         rating.Name,
		Slug:         slug.Approve(ratingSlug, sID, config.DB.NewScope(&model.Rating{}).TableName()),
		Description:  rating.Description,
		MediumID:     rating.MediumID,
		SpaceID:      uint(sID),
		NumericValue: rating.NumericValue,
	}

	tx := config.DB.Begin()
	err = tx.Model(&model.Rating{}).Create(&result).Error

	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	tx.Model(&model.Rating{}).Preload("Medium").First(&result)

	// Insert into meili index
	meiliObj := map[string]interface{}{
		"id":            result.ID,
		"kind":          "rating",
		"name":          result.Name,
		"slug":          result.Slug,
		"description":   result.Description,
		"numeric_value": result.NumericValue,
		"space_id":      result.SpaceID,
	}

	err = meili.AddDocument(meiliObj)
	if err != nil {
		tx.Rollback()
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	tx.Commit()
	renderx.JSON(w, http.StatusCreated, result)
}
