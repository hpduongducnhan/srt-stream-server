package app

type HookRequest struct {
	Action   string `json:"action"`
	ClientID string `json:"client_id"`
	IP       string `json:"ip"`
	Vhost    string `json:"vhost"`
	App      string `json:"app"`
	Stream   string `json:"stream"`
	Param    string `json:"param"`
}

type HookResponse struct {
	Code    int    `json:"code"` // 0=deny, 1=allow
	Message string `json:"message,omitempty"`
}
