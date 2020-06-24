package server

type AuthenticationRequest struct {
	ApiVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
	Spec       struct {
		Token string `json:"token"`
	}
	Status struct {
		User map[string]interface{} `json:"user"`
	}
}

type AuthenticationResponse struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Status     Status `json:"status"`
}

type Status struct {
	Authenticated bool `json:"authenticated"`
	User          User `json:"user"`
}

type User struct {
	Groups   []string `json:"groups"`
	UID      string   `json:"uid"`
	Username string   `json:"username"`
}
