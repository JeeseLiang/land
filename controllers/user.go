package controllers

import (
	"errors"
	"fmt"
	"land/dao/mysql"
	"land/logic"
	"land/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	p := new(models.SignUpForm)

	// var bodyBytes []byte
	// if c.Request.Body != nil {
	// 	bodyBytes, _ = io.ReadAll(c.Request.Body)
	// }
	// zap.L().Info("请求体", zap.ByteString("body", bodyBytes))
	// c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// zap.L().Info("body", zap.Any("body", c.Request.Body))

	// 绑定请求参数，校验参数合法性
	if err := c.ShouldBindJSON(&p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("注册参数无效", zap.Error(err), zap.Any("params:", p))
		// 判断err是否为validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResError(c, CodeInvalidParams)
			return
		}
		ResErrorWithMsg(c, CodeInvalidParams, removeStructName(errs.Translate(trans)))
		return
	}

	// 注册逻辑
	// 2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("注册逻辑处理失败", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResError(c, CodeUserExists)
			return
		}
		ResError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 获取请求参数，校验参数
	p := new(models.LoginForm)
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("登录参数无效", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResError(c, CodeInvalidParams)
			return
		}
		ResErrorWithMsg(c, CodeInvalidParams, removeStructName(errs.Translate(trans)))
		return
	}

	// 登录逻辑
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("登录逻辑处理失败", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResError(c, CodeUserNotFound)
			return
		}
		ResError(c, CodeUserPasswordError)
		return
	}

	// 3. 返回响应
	ResSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID), // ID大于2^53时，JSON序列化可能失真，2^64 > 2^53
		"user_name": user.Username,
		"token":     user.Token,
	})
}
