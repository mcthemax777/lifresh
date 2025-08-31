package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"lifresh/db"
	"lifresh/define"
	"lifresh/request"
	"lifresh/response"
)

type GetItemHandler struct {
}

func (h GetItemHandler) ApiCall(c *gin.Context) {
	ApiCall(c, h.process)
}

func (h GetItemHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetItemReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	var itemID define.SnowflakeID
	err = itemID.UnmarshalJSON([]byte(req.ItemID))
	if err != nil {
		return nil, err
	}

	//자신의 폴더
	if itemID == 0 {

	}
	account, err := db.DBHandlerSG.GetAccountByUID(uid)

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_ITEM_RES)

	getFolderRes := res.(*response.GetFolderRes)

	return ResponseToByteArray(getFolderRes), nil
}
