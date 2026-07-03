package handler

import (
	"github.com/Olive1117/gin-blog/internal/model"
	"github.com/Olive1117/gin-blog/internal/model/convert"
	"github.com/Olive1117/gin-blog/internal/service"
	"github.com/Olive1117/gin-blog/pkg/errs"
	"github.com/Olive1117/gin-blog/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type friendlinkHandler struct {
	service service.FriendLinkService
}

// Create implements [FriendLinkHandler].
func (f *friendlinkHandler) Create(c *gin.Context) {
	cx := c.Request.Context()
	var friendlinkDTO model.FriendLinkDTO
	if err := c.ShouldBindJSON(&friendlinkDTO); err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	friendlink := *convert.FriendLinkFromDTO(&friendlinkDTO)
	if err := f.service.Create(cx, &friendlink); err != nil {
		Fail(c, err)
	}
	friendlinkVO := convert.FriendLinkToVO(&friendlink)
	Success(c, friendlinkVO)
}

// Delete implements [FriendLinkHandler].
func (f *friendlinkHandler) Delete(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "删除友链", logger.Int64("友链ID", id))
	if err := f.service.Delete(cx, id); err != nil {
		Fail(c, err)
		return
	}
	Success(c, nil)
}

// Get implements [FriendLinkHandler].
func (f *friendlinkHandler) Get(c *gin.Context) {
	cx := c.Request.Context()
	id := cast.ToInt64(c.Param("id"))
	friendlink, err := f.service.Get(cx, id)
	if err != nil {
		Fail(c, err)
		return
	}
	friendlinkVO := convert.FriendLinkToVO(&friendlink)
	Success(c, friendlinkVO)
	panic("unimplemented")
}

// List implements [FriendLinkHandler].
func (f *friendlinkHandler) List(c *gin.Context) {
	cx := c.Request.Context()
	pageQue := model.PageQuery{Page: cast.ToInt(c.DefaultQuery("page", "1")), PageSize: cast.ToInt(c.DefaultQuery("page_size", "10"))}
	pageRes, err := f.service.List(cx, pageQue)
	if err != nil {
		Fail(c, err)
		return
	}
	friendlinkVOs := convert.MapSlice(pageRes.List, convert.FriendLinkToVO)
	Success(c, NewPageResponse(friendlinkVOs, pageRes.Total, pageQue))
}

// Update implements [FriendLinkHandler].
func (f *friendlinkHandler) Update(c *gin.Context) {
	cx := c.Request.Context()
	var friendlinkDTO model.FriendLinkDTO
	id := cast.ToInt64(c.Param("id"))
	logger.DebugContext(cx, "更新友链", logger.Int64("友链ID", id))
	if err := c.ShouldBindJSON(&friendlinkDTO); err != nil {
		logger.WarnContext(cx, "参数错误", logger.Err(err))
		Fail(c, errs.ErrInvalidParam)
		return
	}
	friendlink := *convert.FriendLinkFromDTO(&friendlinkDTO)
	if err := f.service.Update(cx, &friendlink, id); err != nil {
		Fail(c, err)
		return
	}
	friendlinkVO := convert.FriendLinkToVO(&friendlink)
	Success(c, friendlinkVO)
}

func NewFriendlinkHandler(service service.FriendLinkService) FriendLinkHandler {
	return &friendlinkHandler{
		service: service,
	}
}
