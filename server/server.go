package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v24/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func Start(baseUrl string, uploadUrl string, org string) error {
	log.Printf("[INFO] START: baseUrl: %s, uploadUrl: %s, org: %s", baseUrl, uploadUrl, org)
	http.HandleFunc("/webhook", func(rw http.ResponseWriter, req *http.Request) {
		log.Println("[DEBUG] received")
		decoder := json.NewDecoder(req.Body)

		var areq AuthenticationRequest
		err := decoder.Decode(&areq)
		if err != nil {
			http.Error(rw, "Failed to decode request body.", 401)
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: areq.Spec.Token},
		)
		tc := oauth2.NewClient(ctx, ts)

		user, err := getUserInfo(baseUrl, areq.Spec.Token)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Failed to get user info: %s", err.Error()), 401)
		}

		client, err := github.NewEnterpriseClient(baseUrl, uploadUrl, tc)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Failed to create github client: %s", err.Error()), 401)
		}
		teams, err := getTeams(ctx, client)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Failed to get teams: %s", err.Error()), 401)
		}

		aresp := &AuthenticationResponse{
			ApiVersion: "authentication.k8s.io/v1beta1",
			Kind:       "TokenReview",
			Status: Status{
				Authenticated: true,
				User: User{
					Username: *user.Login,
					Groups:   teams[org],
				},
			},
		}
		respBytes, err := json.Marshal(aresp)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Failed to marshal: %s", err.Error()), 401)
		}
		fmt.Fprint(rw, string(respBytes))
	})

	err := http.ListenAndServe("localhost:8443", nil)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

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

func getUserInfo(github_base_url string, token string) (github.User, error) {
	var u github.User
	req, _ := http.NewRequest("GET", github_base_url+"/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return u, err
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func getTeams(ctx context.Context, client *github.Client) (map[string][]string, error) {
	listOpt := &github.ListOptions{
		PerPage: 100,
	}
	var teams []*github.Team
	resp := map[string][]string{}

	for {
		tmpTeams, resp, err := client.Teams.ListUserTeams(ctx, listOpt)
		if err != nil {
			log.Fatal(err)
		}

		teams = append(teams, tmpTeams...)
		if resp.NextPage == 0 {
			break
		}
		listOpt.Page = resp.NextPage
	}

	for _, team := range teams {
		if resp[*team.Organization.Login] == nil {
			resp[*team.Organization.Login] = []string{}
		}
		resp[*team.Organization.Login] = append(resp[*team.Organization.Login], *team.Name)
	}

	return resp, nil
}
