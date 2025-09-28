package response

import "ai-software-copyright-server/internal/application/model/table"

type FlashPictureBrowseResponse struct {
	*table.FlashPicture
	UseVisits int `json:"useVisits"`
}

type FlashPictureVisitsResponse struct {
	UserBuyResponse
	*table.FlashPicture
	UseVisits int `json:"useVisits"`
}
