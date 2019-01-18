package loginStruct

type Request struct {
	Account 		string  `json:"account"`
	Password 		string  `json:"password"`
}