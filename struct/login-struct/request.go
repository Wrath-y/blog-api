package login_struct

type Request struct {
	Account 		string  `json:"account"`
	Password 		string  `json:"password"`
}