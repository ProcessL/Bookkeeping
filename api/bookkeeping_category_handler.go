package api

import (
	"strconv"

	"github.com/dotdancer/gogofly/global"
	"github.com/dotdancer/gogofly/model"
	"github.com/dotdancer/gogofly/model/common/response"
	"github.com/dotdancer/gogofly/service"
	"github.com/dotdancer/gogofly/service/dto"
	"github.com/dotdancer/gogofly/utils"
	"github.com/gin-gonic/gin"
)

// BookkeepingCategoryApi 结构体定义了分类管理的API处理器
type BookkeepingCategoryApi struct {
	Service service.BookkeepingCategoryService
}

// CreateCategory godoc
// @Tags BookkeepingCategory
// @Summary 创建分类
// @Description 用户创建一个新的收支分类
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   category_info body dto.CreateCategoryRequest true "分类信息"
// @Success 200 {object} response.Response{data=dto.CategoryResponse,msg=string} "创建成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/categories [post]
func (a *BookkeepingCategoryApi) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+utils.GetErrorMsg(req, err), c)
		return
	}

	userID := utils.GetUserID(c) // 从JWT获取用户ID
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	category, err := a.Service.CreateCategory(userID, req)
	if err != nil {
		response.FailWithMessage("创建分类失败: "+err.Error(), c)
		return
	}

	response.OkWithData(category, c)
}

// GetCategory godoc
// @Tags BookkeepingCategory
// @Summary 获取单个分类信息
// @Description 根据ID获取用户的单个分类信息，包含子分类
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path uint true "分类ID"
// @Success 200 {object} response.Response{data=dto.CategoryResponse,msg=string} "获取成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "分类不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/categories/{id} [get]
func (a *BookkeepingCategoryApi) GetCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的分类ID", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	category, err := a.Service.GetCategoryByID(userID, uint(categoryID))
	if err != nil {
		response.FailWithMessage("获取分类信息失败: "+err.Error(), c)
		return
	}

	response.OkWithData(category, c)
}

// UpdateCategory godoc
// @Tags BookkeepingCategory
// @Summary 更新分类信息
// @Description 根据ID更新用户的分类信息
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path uint true "分类ID"
// @Param   category_info body dto.UpdateCategoryRequest true "待更新的分类信息"
// @Success 200 {object} response.Response{data=dto.CategoryResponse,msg=string} "更新成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "分类不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/categories/{id} [put]
func (a *BookkeepingCategoryApi) UpdateCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的分类ID", c)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("请求参数错误: "+utils.GetErrorMsg(req, err), c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	updatedCategory, err := a.Service.UpdateCategory(userID, uint(categoryID), req)
	if err != nil {
		response.FailWithMessage("更新分类失败: "+err.Error(), c)
		return
	}

	response.OkWithData(updatedCategory, c)
}

// DeleteCategory godoc
// @Tags BookkeepingCategory
// @Summary 删除分类
// @Description 根据ID删除用户的分类
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   id path uint true "分类ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Failure 400 {object} response.Response{msg=string} "请求参数错误"
// @Failure 404 {object} response.Response{msg=string} "分类不存在"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/categories/{id} [delete]
func (a *BookkeepingCategoryApi) DeleteCategory(c *gin.Context) {
	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("无效的分类ID", c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	if err := a.Service.DeleteCategory(userID, uint(categoryID)); err != nil {
		response.FailWithMessage("删除分类失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("分类删除成功", c)
}

// ListCategories godoc
// @Tags BookkeepingCategory
// @Summary 获取分类列表 (层级)
// @Description 获取用户的所有分类，支持按类型过滤，并以层级结构返回顶级分类及其子分类
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   type query string false "分类类型 (income/expense)"
// @Param   parent_id query int false "父分类ID (查询指定父分类下的子分类，0表示顶级分类)"
// @Success 200 {object} response.Response{data=[]dto.CategoryResponse,msg=string} "获取成功"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/categories [get]
func (a *BookkeepingCategoryApi) ListCategories(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	categoryType := model.CategoryType(c.Query("type"))
	parentIDStr := c.Query("parent_id")
	var parentID *uint
	if parentIDStr != "" {
		pID, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err == nil {
			tempPID := uint(pID)
			parentID = &tempPID
		} else {
			global.Logger.Warn("Invalid parent_id query parameter: " + parentIDStr)
		}
	}

	categories, err := a.Service.ListCategories(userID, categoryType, parentID)
	if err != nil {
		response.FailWithMessage("获取分类列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(categories, c)
}

// ListAllCategoriesFlat godoc
// @Tags BookkeepingCategory
// @Summary 获取所有分类列表 (扁平)
// @Description 获取用户的所有分类，以扁平列表形式返回，主要用于下拉选择框等场景
// @Accept  json
// @Produce  json
// @Param   x-token header string true "令牌"
// @Param   type query string false "分类类型 (income/expense)"
// @Success 200 {object} response.Response{data=[]dto.CategoryResponse,msg=string} "获取成功"
// @Failure 500 {object} response.Response{msg=string} "服务器内部错误"
// @Router /bk/categories/flat [get]
func (a *BookkeepingCategoryApi) ListAllCategoriesFlat(c *gin.Context) {
	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录或无法获取用户信息", c)
		return
	}

	categoryType := model.CategoryType(c.Query("type"))

	categories, err := a.Service.GetAllCategoriesFlat(userID, categoryType)
	if err != nil {
		response.FailWithMessage("获取分类列表失败: "+err.Error(), c)
		return
	}

	response.OkWithData(categories, c)
}
