package rental_info

import (
	"bytes"
	"encoding/json"
	"frog/module/common/constant"
	"frog/module/common/model/api_models"
	"frog/module/common/model/db_models"
	"frog/module/common/model/es_model"
	comTools "frog/module/common/tools"
	"frog/module/main_service/internal/config"
	"frog/module/main_service/internal/log"
	"github.com/gin-gonic/gin"
)

/**
DescribeRentalInfo
*/

func GetRentalInfoListController() api_models.ApiInterface {
	return &RentalInfoListRequest{}
}

type RentalInfoListRequest struct {
	PageNumber int                 `json:"PageNumber" validate:"required,gt=0"`
	PageSize   int                 `json:"PageSize" validate:"required,gt=0"`
	Filters    []api_models.Filter `json:"Filters"`
}

type RentalInfoResp []RentalListItem

type RentalListItem struct {
	ID    int       `json:"ID"`
	Title string    `json:"Title"`
	Cover string    `json:"Cover"`
	Price []float64 `json:"Price"`
}

func (r *RentalInfoListRequest) GetResult(ctx *gin.Context) (interface{}, *api_models.APIError) {
	var (
		queryObject = map[string]interface{}{}
		reqId       = ctx.GetString(constant.CtxKeyRequestID)
		res         = make(RentalInfoResp, 0)
	)

	queryObject = map[string]interface{}{
		"query": nil,
	}

	if len(r.Filters) > 0 {
		matchCondition := make([]interface{}, 0)
		for i := 0; i < len(r.Filters); i++ {
			_, ok := db_models.GetRentalInfoCols()[r.Filters[i].Key]
			if !ok {
				return nil, &api_models.APIError{
					Code:    constant.CodeBadRequest,
					Message: "error filter params",
				}
			}

			switch r.Filters[i].Type {
			case constant.FilterTypeEQ:
				//db.Where(fmt.Sprintf("%s = ?", r.Filters[i].Key), r.Filters[i].Value)
				matchCondition = append(matchCondition, map[string]interface{}{
					"match": map[string]interface{}{
						r.Filters[i].Key: r.Filters[i].Value,
					},
				})
			}
		}

		queryObject["query"] = map[string]interface{}{
			"bool": map[string]interface{}{
				"must": matchCondition,
			},
		}
	} else {
		queryObject["query"] = map[string]interface{}{
			"match_all": struct{}{},
		}
	}

	queryJsonData, _ := json.Marshal(queryObject)

	resp, err := config.GetESCli().Search(
		config.GetESCli().Search.WithContext(ctx),
		config.GetESCli().Search.WithIndex(config.GetESIndexByConfig(es_model.RentalInfo{}.Index())),
		config.GetESCli().Search.WithBody(bytes.NewReader(queryJsonData)),
	)
	if err != nil {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	if resp != nil {
		defer resp.Body.Close()

		if resp.IsError() {
			log.Errorf(reqId, resp.String())
			return nil, &api_models.APIError{
				Code:    constant.CodeInternalError,
				Message: constant.MsgInternalError,
			}
		}

		esModels := make([]es_model.RentalInfo, 0)
		err := comTools.GetModelFromESResp(resp, &esModels)
		if err != nil {
			log.Errorf(reqId, err.Error())
			return nil, &api_models.APIError{
				Code:    constant.CodeInternalError,
				Message: constant.MsgInternalError,
			}
		}

		for _, model := range esModels {
			res = append(res, RentalListItem{
				ID:    int(model.ID),
				Title: model.Title,
				Cover: model.Cover,
				Price: model.Price,
			})
		}
	}

	values, err := comTools.Pagination(res, r.PageNumber, r.PageSize)
	if err != nil {
		log.Errorf(reqId, err.Error())
		return nil, &api_models.APIError{
			Code:    constant.CodeInternalError,
			Message: constant.MsgInternalError,
		}
	}

	return values.Interface(), nil
}
