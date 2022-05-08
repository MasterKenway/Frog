package rental_info

import (
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"

	"github.com/gin-gonic/gin"
)

func GetRentalInfoColsController() api_models.ApiInterface {
	return &RentalInfoColsRequest{}
}

type RentalInfoColsRequest struct {
}

func (req *RentalInfoColsRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		cols = make([]string, 0)
	)

	for key := range db_models.GetRentalInfoCols() {
		cols = append(cols, key)
	}

	return cols, nil
}
