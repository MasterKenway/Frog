package controller

import (
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/main_service/internal/controller/auth/email"
	"frog/module/main_service/internal/controller/auth/login"
	"frog/module/main_service/internal/controller/auth/register"
	"frog/module/main_service/internal/controller/comment"
	"frog/module/main_service/internal/controller/rental_info"
)

var (
	ApiAdapter = map[string]api_models.ControllerAdapter{
		constant.ApiLogin:        login.GetLoginRequestAdapter,
		constant.ApiRegister:     register.GetRegisterAdapter,
		constant.ApiGetEmailCode: email.GetEmailCodeAdapter,

		constant.ApiRentalInfoDetail: rental_info.GetRentalInfoDetailController,
		constant.ApiRentalInfoList:   rental_info.GetRentalInfoListController,
		constant.ApiRentalInfoCreate: rental_info.GetRentalInfoCreateController,

		constant.ApiCommentCreate: comment.GetCommentCreateController,
	}
)
