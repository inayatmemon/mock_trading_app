package commonmodels

type ApiResponse struct {
	Code           int         `json:"code"`
	Message        string      `json:"message"`
	Data           interface{} `json:"data"`
	TotalCount     int64       `json:"totalCount"`
	TotalPageCount int64       `json:"totalPageCount"`
}
