package doc

import (
	"ai-service/internal/repository"
	"ai-service/internal/service/doc/models"
	"ai-service/internal/util/errors"
	"github.com/labstack/echo"
	"net/http"
)

// SaveDoc
//
// @Description Save document
// @Summary	Save document in milvus
// @Tags doc
// @Accept json
// @Produce	json
// @Success	200				{object}		models.SaveDocResponse
// @Failure	500				{object}	status.SaveDoc
// @Router /api/v1/document 	[post]
func (d *docService) SaveDoc(c echo.Context) error {
	var dataReq models.SaveDoc
	if err := c.Bind(&dataReq); err != nil {
		return errors.NewBadRequestErrorRsp(err.Error())
	}
	response := models.SaveDocResponse{Status: true}
	return c.JSON(http.StatusOK, response)
}

// DeleteDoc
//
// @Description Delete document
// @Summary	Delete document from milvus
// @Tags doc
// @Accept json
// @Produce	json
// @Success	200				{object}
// @Failure	500				{object}	status.SaveDoc
// @Router /api/v1/document/{id} 	[delete]
func (d *docService) DeleteDoc(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	var dataReq models.DeleteDoc
	if err := c.Bind(&dataReq); err != nil {
		return errors.NewBadRequestErrorRsp(err.Error())
	}
	err := d.repository.Vector.DeleteDoc(ctx, id)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}
	return c.JSON(http.StatusOK, nil)
}

// UpdatePriority
//
//	@Description Update document priority
//	@Summary	Update document priority in milvus
//	@Tags		doc
//	@Accept		json
//	@Produce	json
//	@Success	200				{object}
//	@Failure	500				{object}	status.SaveDoc
//
// @Router		/api/v1/document/{id} 	[put]
func (d *docService) UpdatePriority(c echo.Context) error {
	var dataReq models.UpdatePriority
	if err := c.Bind(&dataReq); err != nil {
		return errors.NewBadRequestErrorRsp(err.Error())
	}
	response := models.SaveDocResponse{Status: true}
	return c.JSON(http.StatusOK, response)
}

func (d *docService) convert(base64 string) (string, error) {
	return base64, nil
}
