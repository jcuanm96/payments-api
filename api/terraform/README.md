#  Terraform
https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/gcp-get-started#install-terraform

## Install:
1. `brew tap hashicorp/tap`<br/>
2. `brew install hashicorp/tap/terraform`<br/>
3. `brew update`<br/>
4. `brew upgrade hashicorp/tap/terraform`

## Terminal tab complete for terraform subcommands:
If Bash: `touch ~/.bashrc`<br/>
If Zsh: `touch ~/.zshrc`<br/><br/>

`terraform -install-autocomplete`<br/>

Might need to add following to .zshrc before terraform complete stuff:
`autoload -Uz compinit`<br/>
`compinit`

## Vscode:
Install hashicorp.terraform extension

## Terraform commands:
1. Run `terraform init`
2. Run `terraform fmt` to format 
3. Run `terraform validate` to validate (includes checking paths and lots of cool stuff)
4. Run `terraform plan` to see diff in state/config before applying

To actually create the resources and apply any changes (execution plan is smart and will only include the resources that are new/changed/deleted from previous state):
    `terrfaform apply`
    Type `yes` if all looks good
To show current state:
    `terraform show`
To destroy current state:
    `terraform destroy` (probably don't wanna run this one often)
    Terraform builds a dependency graph to destroy things in the right order

## Terraform files:
Resource definiton order does not matter so whatever is easiest to manage/visually appealing.

Variable Files:
    We can define variable files like `variables.tf`- useful for using the same config file but different values for different services (dev vs staging vs prod).
    Variables are accessed by the prefix `var.` in the config file.
    
    We then define a `terraform.tfvars` file that contains the values for any variables
    where we override the default or values that are secret/credentials.
    THIS SHOULD NOT BE CHECKED INTO GIT. It should be passed around just like our
    `.env.local.rc` and `creds.json` etc

Output:
    We can define output files like `outputs.tf`- useful for displaying a subset of the
    massive number of attributes that Terraform produces.

    After adding an output value, you must `terraform apply`. After that, you can see
    the output values for the current state with `terraform output`


## Useful Links
https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudbuild_trigger

https://cloud.google.com/build/docs/api/reference/rest/v1/projects.locations.triggers

https://cloud.google.com/architecture/managing-infrastructure-as-code

https://medium.com/@Paul_D/terraform-provisioning-gcp-cloud-run-cloud-endpoint-container-registry-using-cloud-build-to-6a2a734f8077

https://github.com/hashicorp/terraform-provider-google/issues/6635
