package handler

import (
	"net/http"
	"strconv"

	"xx-backend/internal/model"
	"xx-backend/internal/service"

	"github.com/gin-gonic/gin"
)

func GetUsers(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		search := c.Query("search")

		users, total, err := userService.GetUsers(page, pageSize, search)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取用户列表失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data": gin.H{
				"list":  users,
				"total": total,
				"page":  page,
				"size":  pageSize,
			},
		})
	}
}

func GetUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}

		user, err := userService.GetUser(int(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "用户不存在",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    user,
		})
	}
}

func CreateUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		if err := userService.CreateUser(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建用户失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "创建成功",
			"data":    user,
		})
	}
}

func UpdateUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		if err := userService.UpdateUser(int(id), updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新用户失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新成功",
		})
	}
}

func DeleteUser(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
			return
		}

		if err := userService.DeleteUser(int(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除用户失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除成功",
		})
	}
}

// 角色相关处理器
func GetRoles(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := userService.GetRoles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取角色列表失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    roles,
		})
	}
}

func CreateRole(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role model.Role
		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		if err := userService.CreateRole(&role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建角色失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "创建成功",
			"data":    role,
		})
	}
}

func UpdateRole(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的角色ID",
			})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		if err := userService.UpdateRole(uint(id), updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新角色失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新成功",
		})
	}
}

func DeleteRole(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的角色ID",
			})
			return
		}

		if err := userService.DeleteRole(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除角色失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除成功",
		})
	}
}

// 菜单相关处理器
func GetMenus(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		menus, err := userService.GetMenus()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取菜单列表失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "获取成功",
			"data":    menus,
		})
	}
}

func CreateMenu(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu model.Menu
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		if err := userService.CreateMenu(&menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "创建菜单失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "创建成功",
			"data":    menu,
		})
	}
}

func UpdateMenu(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的菜单ID",
			})
			return
		}

		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "请求参数错误",
				"error":   err.Error(),
			})
			return
		}

		if err := userService.UpdateMenu(uint(id), updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "更新菜单失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新成功",
		})
	}
}

func DeleteMenu(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "无效的菜单ID",
			})
			return
		}

		if err := userService.DeleteMenu(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "删除菜单失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "删除成功",
		})
	}
}
