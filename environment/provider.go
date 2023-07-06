// Copyright 2019 Form3 Financial Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package githubdeploymentenvironment

import (
	"context"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/oauth2"
)

const (
	resourceFileName = "github-deployment-environment_branch_policy"
)

const (
	githubTokenKey = "github_token"
	ownerKey       = "github_owner"
)

type providerConfiguration struct {
	githubClient *github.Client
	owner        string
}

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: d.Get(githubTokenKey).(string),
				},
			)
			tc := oauth2.NewClient(context.Background(), ts)
			return &providerConfiguration{
				githubClient: github.NewClient(tc),
				owner:        d.Get(ownerKey).(string),
			}, nil
		},
		ResourcesMap: map[string]*schema.Resource{
			resourceFileName: resourceGithubRepositoryEnvironmentDeploymentPolicy(),
		},
		Schema: map[string]*schema.Schema{
			githubTokenKey: {
				Type:        schema.TypeString,
				DefaultFunc: defaultFuncForKey(githubTokenKey),
				Required:    true,
				Sensitive:   true,
				Description: "A GitHub authorisation token with permissions to manage repositories.",
			},
			ownerKey: {
				Type:        schema.TypeString,
				DefaultFunc: defaultFuncForKey(ownerKey),
				Required:    true,
				Sensitive:   false,
				Description: "Github Organization",
			},
		},
	}
}

func defaultFuncForKey(v string) schema.SchemaDefaultFunc {
	return schema.EnvDefaultFunc(strings.ToUpper(v), "")
}
