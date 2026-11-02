package handlers

import (
	"photo-web/internal/api"
	"strconv"

	"github.com/gin-gonic/gin"

	"photo-web/internal/service"
)

// HistoryHandler 历史记录处理器
type HistoryHandler struct {
	historyService *service.HistoryService
	uploadDir      string
}

// NewHistoryHandler 创建历史 Handler
func NewHistoryHandler(historyService *service.HistoryService, uploadDir string) *HistoryHandler {
	return &HistoryHandler{
		historyService: historyService,
		uploadDir:      uploadDir,
	}
}

// GetHistory 获取历史列表 (分页)
// GET /api/v1/history?page=1&page_size=20
func (h *HistoryHandler) GetHistory(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		api.Unauthorized(c, "")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := h.historyService.GetHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		api.InternalError(c, err.Error())
		return
	}

	api.Success(c, result)
}

// GetHistoryDetail 获取历史详情
// GET /api/v1/history/:id
func (h *HistoryHandler) GetHistoryDetail(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		api.Unauthorized(c, "")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		api.BadRequest(c, "无效的 ID")
		return
	}

	record, err := h.historyService.GetHistoryDetail(c.Request.Context(), id, userID)
	if err != nil {
		api.NotFound(c, err.Error())
		return
	}

	api.Success(c, record)
}

// DeleteHistory 删除历史记录
// DELETE /api/v1/history/:id
func (h *HistoryHandler) DeleteHistory(c *gin.Context) {
	userID, err := GetUserID(c)
	if err != nil {
		api.Unauthorized(c, "")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		api.BadRequest(c, "无效的 ID")
		return
	}

	if err := h.historyService.DeleteHistory(c.Request.Context(), id, userID, h.uploadDir); err != nil {
		api.Error(c, 400, api.CodeHistoryNotFound, err.Error())
		return
	}

	api.Success(c, nil)
}
