package response

import "ai-software-copyright-server/internal/application/model/table"

type UserLoginResponse struct {
	*table.User
	RewardCredits int    `json:"rewardCredits"`
	Message       string `json:"message,omitempty"`
}

type UserRewardResponse struct {
	RewardCredits  int    `json:"rewardCredits"`
	BalanceCredits int    `json:"balanceCredits"`
	Message        string `json:"message,omitempty"`
}

type UserBuyResponse struct {
	BuyCredits     int    `json:"buyCredits"`
	BalanceCredits int    `json:"balanceCredits"`
	BuyMessage     string `json:"buyMessage,omitempty"`
}

type UserBuyContentResponse struct {
	UserBuyResponse
	Content string `json:"content"`
}
