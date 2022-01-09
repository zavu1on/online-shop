package schemas

type AddGoodSchema struct {
	GoodID int `json:"good_id" binding:"required"`
}