package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v24/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

func Start(c *cli.Context) error {
	fmt.Println("[DEBUG] Start")
	http.HandleFunc("/webhook", func(rw http.ResponseWriter, req *http.Request) {
		log.Println("[DEBUG] received")
		decoder := json.NewDecoder(req.Body)

		var areq AuthenticationRequest
		err := decoder.Decode(&areq)
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: areq.Spec.Token},
		)
		tc := oauth2.NewClient(ctx, ts)

		user, err := getUserInfo(c.String("github-base-url"), areq.Spec.Token)
		if err != nil {
			log.Fatal(err)
		}

		client, err := github.NewEnterpriseClient(c.String("github-base-url"), c.String("github-upload-url"), tc)
		if err != nil {
			log.Fatal(err)
		}
		teams, err := getTeams(ctx, client)
		if err != nil {
			log.Fatal(err)
		}

		aresp := &AuthenticationResponse{
			ApiVersion: "authentication.k8s.io/v1beta1",
			Kind:       "TokenReview",
			Status: Status{
				Authenticated: true,
				User: User{
					Username: *user.Login,
					Groups:   teams[c.String("team")],
				},
			},
		}
		respBytes, err := json.Marshal(aresp)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[DEBUG] resp: %s\n", string(respBytes))
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
