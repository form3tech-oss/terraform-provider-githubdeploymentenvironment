A Terraform provider for managing files in GitHub repositories.

## Installation

Download the relevant binary from [releases](https://github.com/form3tech-oss/terraform-provider-githubdeploymentenvironment/releases) and copy it to `$HOME/.terraform.d/plugins/`.

## Configuration

The following provider block variables are available for configuration:

| Name | Description |
| ---- | ----------- |
| `github_token` | A GitHub authorisation token with `repo` permissions. |
| `github_owner` | The name of the Github organization in which to use the provider. |

Alternatively, these values can be read from environment variables.

## Resources

### `github-deployment-environment_branch_policy`

The `github-deployment-environment_branch_policy` resource represents a deployment environment branch policy.

:warning: [Custom branch policies](https://registry.terraform.io/providers/integrations/github/latest/docs/resources/repository_environment#custom_branch_policies) must be enabled, otherwise the provider is gonna fail!

#### Attributes

| Name | Description |
| ---- | ----------- |
| `repository` | The name of the repository. |
| `environment` | The name of the environment. |
| `branch_pattern` | The branch pattern on which to restrict environment deployment on. |

#### Example

```hcl
resource "github-deployment-environment_branch_policy" "my-resource" {
    repository 	   = "my-repo"
    environment	   = "test-environment"
    branch_pattern = "master"
}
```
