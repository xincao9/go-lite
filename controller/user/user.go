package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golite/constant"
	"golite/internal/util"
	"golite/model/user"
	userService "golite/service/user"
	"net/http"
	"strconv"
	"time"
)

func Route(engine *gin.Engine) {
	save := func(c *gin.Context) {
		u := &user.User{}
		if err := c.ShouldBindJSON(u); err != nil {
			util.RenderJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		err := userService.U.Login(u) // 登录校验
		if err != nil {
			util.RenderJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		u.Token = uuid.New().String()
		u.Expire = time.Now().Add(time.Hour * time.Duration(constant.SessionExpireHour))
		err = userService.U.Save(u) // 更新登录信息
		if err != nil {
			util.RenderJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.SetCookie(constant.Token, u.Token, (int)((time.Hour*time.Duration(constant.SessionExpireHour))/time.Second), "/", "*", false, false)
		util.RenderJSONDetail(c, http.StatusOK, constant.Success, u)
	}
	engine.PUT("/session", save)
	engine.POST("/session", save)
}

func AuthenticationRoute(engine *gin.RouterGroup) {
	engine.DELETE("/session/:id", func(c *gin.Context) {
		su, ok := c.Get(constant.SessionUser)
		if ok == false {
			util.RenderJSON(c, http.StatusInternalServerError, constant.SystemError)
			return
		}
		u := su.(*user.User)
		u, err := userService.U.GetUserByUsername(u.Username)
		if err != nil {
			util.RenderJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		u.Expire = time.Now()
		err = userService.U.Save(u)
		if err != nil {
			util.RenderJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		util.RenderJSON(c, http.StatusOK, constant.Success)
	})
	engine.GET("/user", func(c *gin.Context) {
		su, ok := c.Get(constant.SessionUser)
		if ok == false {
			util.RenderJSON(c, http.StatusInternalServerError, constant.SystemError)
			return
		}
		u := su.(*user.User)
		util.RenderJSONDetail(c, http.StatusOK, constant.Success, u)
	})
	engine.GET("/user/:id", func(c *gin.Context) {
		n := c.Param("id")
		if n == "" {
			util.RenderJSON(c, http.StatusBadRequest, "id is required")
			return
		}
		id, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			util.RenderJSON(c, http.StatusBadRequest, err.Error())
			return
		}
		u, err := userService.U.GetUserById(id)
		if err != nil {
			util.RenderJSON(c, http.StatusInternalServerError, err.Error())
			return
		}
		util.RenderJSONDetail(c, http.StatusOK, constant.Success, u)
	})
}
