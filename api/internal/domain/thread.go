package domain

type Thread struct {
	Id      int32  `json:"id" readonly:"true" example:"42"`
	Title   string `json:"title" example:"Davy Jones cache"`
	Author  string `json:"author" example:"j.sparrow"`
	Forum   string `json:"forum" example:"pirate-stories"`
	Message string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"`
	Votes   int32  `json:"votes"`
	Slug    string `json:"slug" example:"jones-cache"`
	Created int64  `json:"created"`
}

type ThreadUpdate struct {
	Title   string `json:"title" example:"Davy Jones cache"`
	Message string `json:"message" example:"An urgent need to reveal the hiding place of Davy Jones. Who is willing to help in this matter?"`
}
