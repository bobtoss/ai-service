package doc

import (
	"ai-service/internal/service/doc/models"
	"ai-service/internal/service/document"
	"ai-service/internal/util/doc"
	"ai-service/internal/util/errors"
	"ai-service/internal/util/middleware"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// SaveDoc
//
// @Description Save document
// @Summary	Save document in milvus
// @Tags doc
// @Accept json
// @Produce	json
// @Param		request	body		models.SaveDocResponse	true	"body param"
// @Success	200				{object}		models.SaveDocResponse
// @Router /api/v1/document 	[post]
func (d *docService) SaveDoc(c echo.Context) error {
	ctx := c.Request().Context()
	uid, ok := middleware.UserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "user not found")
	}
	var dataReq models.SaveDoc
	if err := c.Bind(&dataReq); err != nil {
		return errors.NewBadRequestErrorRsp(err.Error())
	}

	docID := uuid.New().String()
	encodedString, err := doc.DecodeBase64ToFileAndRead(dataReq.Document)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}
	embeddings, err := d.llm.Embed(encodedString)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}
	err = d.repository.Vector.SaveDoc(ctx, uid, docID, encodedString, embeddings)
	if err != nil {
		fmt.Println(err)
		return errors.NewInternalErrorRsp(err.Error())
	}

	pgModel := document.Document{
		DocumentID:   docID,
		UserID:       uid,
		DocumentName: dataReq.Name,
	}
	err = d.postgres.Create(ctx, &pgModel)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}

	response := models.SaveDocResponse{Status: true}
	return c.JSON(http.StatusOK, response)
}

// ListDoc
//
// @Description List documents
// @Summary	List documents in milvus
// @Tags doc
// @Accept json
// @Produce	json
// @Router /api/v1/document 	[get]
func (d *docService) ListDoc(c echo.Context) error {
	ctx := c.Request().Context()
	uid, ok := middleware.UserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "user not found")
	}

	list, err := d.postgres.ListByUser(ctx, uid)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

// DeleteDoc
//
// @Description Delete document
// @Summary	Delete document from milvus
// @Tags doc
// @Accept json
// @Produce	json
// @Success	200
// @Router /api/v1/doc/{id} 	[delete]
func (d *docService) DeleteDoc(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	uid, ok := middleware.UserIDFromContext(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, "user not found")
	}

	var dataReq models.DeleteDoc
	if err := c.Bind(&dataReq); err != nil {
		return errors.NewBadRequestErrorRsp(err.Error())
	}
	err := d.repository.Vector.DeleteDoc(ctx, uid, id)
	if err != nil {
		return errors.NewInternalErrorRsp(err.Error())
	}
	err = d.postgres.Delete(ctx, id)
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
//	@Success	200				{object} models.SaveDocResponse
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
