package rental_info

import (
	"bytes"
	"fmt"

	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/es_model"
	"frog/module/common/tools"
	"frog/module/main_service/internal/config"

	"github.com/gin-gonic/gin"
)

func GetRentalInfoDetailController() api_models.ApiInterface {
	return &DetailRequest{}
}

type DetailRequest struct {
	Id int `json:"ID"`
}

func (req *DetailRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	queryObject := []byte(fmt.Sprintf(`
{
    "query": {
"match": {
                        "id": %d
                    }
    }
}`, req.Id))

	resp, err := config.GetESCli().Search(
		config.GetESCli().Search.WithContext(ctx),
		config.GetESCli().Search.WithIndex(config.GetESIndexByConfig(es_model.RentalInfo{}.Index())),
		config.GetESCli().Search.WithBody(bytes.NewReader(queryObject)),
	)
	if err != nil {
		return nil, &api_models.APIError{
			Code:    constant.CodeBadRequest,
			Message: err.Error(),
		}
	}

	hits := make([]es_model.RentalInfo, 0)
	if resp != nil {
		defer resp.Body.Close()

		if resp.IsError() {
			return nil, &api_models.APIError{
				Code:    constant.CodeBadRequest,
				Message: resp.String(),
			}
		}

		err := tools.GetModelFromESResp(resp, &hits)
		if err != nil {
			return nil, &api_models.APIError{
				Code:    constant.CodeBadRequest,
				Message: err.Error(),
			}
		}
	}

	return hits[0], nil
}
