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
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func resourceGithubRepositoryEnvironmentDeploymentPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceGithubRepositoryEnvironmentDeploymentPolicyCreate,
		Read:   resourceGithubRepositoryEnvironmentDeploymentPolicyRead,
		Update: resourceGithubRepositoryEnvironmentDeploymentPolicyUpdate,
		Delete: resourceGithubRepositoryEnvironmentDeploymentPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the repository. The name is not case sensitive.",
			},
			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the environment.",
			},
			"branch_pattern": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The name pattern that branches must match in order to deploy to the environment.",
			},
		},
	}

}

func resourceGithubRepositoryEnvironmentDeploymentPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerConfiguration).githubClient
	ctx := context.Background()

	owner := meta.(*providerConfiguration).owner
	repoName := d.Get("repository").(string)
	envName := d.Get("environment").(string)
	branchPattern := d.Get("branch_pattern").(string)
	escapedEnvName := url.QueryEscape(envName)

	createData := github.DeploymentBranchPolicyRequest{
		Name: github.String(branchPattern),
	}

	resultKey, _, err := client.Repositories.CreateDeploymentBranchPolicy(ctx, owner, repoName, escapedEnvName, &createData)
	if err != nil {
		return err
	}

	d.SetId(buildThreePartID(repoName, escapedEnvName, strconv.FormatInt(resultKey.GetID(), 10)))
	return resourceGithubRepositoryEnvironmentDeploymentPolicyRead(d, meta)
}

func resourceGithubRepositoryEnvironmentDeploymentPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerConfiguration).githubClient

	owner := meta.(*providerConfiguration).owner
	repoName, envName, branchPolicyIdString, err := parseThreePartID(d.Id(), "repository", "environment", "branchPolicyId")
	if err != nil {
		return err
	}

	branchPolicyId, err := strconv.ParseInt(branchPolicyIdString, 10, 64)
	if err != nil {
		return err
	}

	escapedEnvName := url.QueryEscape(envName)

	branchPolicy, _, err := client.Repositories.GetDeploymentBranchPolicy(context.Background(), owner, repoName, escapedEnvName, branchPolicyId)
	if err != nil {
		if ghErr, ok := err.(*github.ErrorResponse); ok {
			if ghErr.Response.StatusCode == http.StatusNotModified {
				return nil
			}
			if ghErr.Response.StatusCode == http.StatusNotFound {
				fmt.Printf("[INFO] Removing branch deployment policy for %s/%s/%s from state because it no longer exists in GitHub",
					owner, repoName, envName)
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("branch_pattern", branchPolicy.GetName())
	return nil
}

func resourceGithubRepositoryEnvironmentDeploymentPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerConfiguration).githubClient
	ctx := context.Background()

	owner := meta.(*providerConfiguration).owner
	repoName := d.Get("repository").(string)
	envName := d.Get("environment").(string)
	branchPattern := d.Get("branch_pattern").(string)
	escapedEnvName := url.QueryEscape(envName)
	_, _, branchPolicyIdString, err := parseThreePartID(d.Id(), "repository", "environment", "branchPolicyId")
	if err != nil {
		return err
	}

	branchPolicyId, err := strconv.ParseInt(branchPolicyIdString, 10, 64)
	if err != nil {
		return err
	}

	updateData := github.DeploymentBranchPolicyRequest{
		Name: github.String(branchPattern),
	}

	resultKey, _, err := client.Repositories.UpdateDeploymentBranchPolicy(ctx, owner, repoName, escapedEnvName, branchPolicyId, &updateData)
	if err != nil {
		return err
	}
	d.SetId(buildThreePartID(repoName, escapedEnvName, strconv.FormatInt(resultKey.GetID(), 10)))
	return resourceGithubRepositoryEnvironmentDeploymentPolicyRead(d, meta)
}

func resourceGithubRepositoryEnvironmentDeploymentPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerConfiguration).githubClient
	ctx := context.Background()

	owner := meta.(*providerConfiguration).owner
	repoName, envName, branchPolicyIdString, err := parseThreePartID(d.Id(), "repository", "environment", "branchPolicyId")
	if err != nil {
		return err
	}

	branchPolicyId, err := strconv.ParseInt(branchPolicyIdString, 10, 64)
	if err != nil {
		return err
	}

	escapedEnvName := url.QueryEscape(envName)

	_, err = client.Repositories.DeleteDeploymentBranchPolicy(ctx, owner, repoName, escapedEnvName, branchPolicyId)
	if err != nil {
		return err
	}

	return nil
}

// return the pieces of id `left:center:right` as left, center, right
func parseThreePartID(id, left, center, right string) (string, string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("unexpected ID format (%q). Expected %s:%s:%s", id, left, center, right)
	}

	return parts[0], parts[1], parts[2], nil
}

// format the strings into an id `a:b:c`
func buildThreePartID(a, b, c string) string {
	return fmt.Sprintf("%s:%s:%s", a, b, c)
}
