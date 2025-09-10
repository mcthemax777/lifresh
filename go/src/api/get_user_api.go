package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"lifresh/db"
	"lifresh/define"
	"lifresh/internal/transport/http/dto/request"
	"lifresh/internal/transport/http/dto/response"
	"lifresh/internal/transport/http/handler"
	"lifresh/models"
)

type GetUserHandler struct {
}

func (h GetUserHandler) ApiCall(c *gin.Context) {
	handler.ApiCall(c, h.process)
}

func (h GetUserHandler) process(reqBody []byte) ([]byte, error) {

	var req request.GetUserReq
	err := json.Unmarshal(reqBody, &req)

	if err != nil {
		return handler.ResponseToByteArray(response.CreateFailResponse(201, "invalid_json")), err
	}

	var userID define.SnowflakeID
	err = userID.UnmarshalJSON([]byte(req.UserID))
	if err != nil {
		return nil, err
	}

	var root models.Folder

	user, folderList, planList, err := db.DBHandlerSG.GetUserByUID(userID)

	folderMap := make(map[define.SnowflakeID]*models.Folder)
	for _, folder := range folderList {
		folderMap[folder.ID] = &folder
	}

	for _, folder := range folderList {
		if folder.ParentID == nil {
			root = folder
		} else {
			parent := folderMap[*folder.ParentID]
			parent.ChildrenFolders = append(parent.ChildrenFolders, folder)
		}
	}

	for _, plan := range planList {
		if plan.ParentID == nil {
			// 부모 폴더가 없는게 말이 안됨. error
			return handler.ResponseToByteArray(response.CreateFailResponse(202, "no parent")), err
		} else {
			parent := folderMap[*plan.ParentID]
			parent.ChildrenFolders = append(parent.ChildrenFolders, plan)
		}
	}

	user.Root = &root

	//전송할 데이터 만들기
	res := response.CreateSuccessResponse(response.GET_ITEM_RES)

	getUserRes := res.(*response.GetFolderRes)
	getUserRes.Folder = account

	return handler.ResponseToByteArray(loginRes), nil
}
