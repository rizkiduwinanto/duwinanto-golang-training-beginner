package err

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
