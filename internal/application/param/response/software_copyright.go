package response

type SCRequirementItemResponse struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Operation string `json:"operation"`
}

type SCRequestInfoResponse struct {
	Purpose  string `json:"purpose"`  //开发目的
	Oriented string `json:"oriented"` //面向领域/行业
	Function string `json:"function"` //软件的主要功能
	Feature  string `json:"feature"`  //软件的技术特点
}
