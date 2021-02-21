# App Engine Template

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/ad712e2c9fa34594b4974922deda63d9)](https://app.codacy.com/app/idestis/app_engine_template?utm_source=github.com&utm_medium=referral&utm_content=idestis/app_engine_template&utm_campaign=Badge_Grade_Dashboard)

This tool allows you to quickly and easily create app.yaml

You configure your App Engine app's settings in the [app.yaml](https://cloud.google.com/appengine/docs/standard/go111/config/appref) file. 
Sometimes we need to configure sensitive env_variables and they should not be in the same 
repository as a code base.
This tool allows you to grab secrets from [Hashicorp Vault](https://www.vaultproject.io/)

Tool is based on package `text/template`, based on this you should understand the basics 
of this [package](https://godoc.org/text/template)

We support variables, pipelines with if-else and loops with range.

## Synopsis

`app_engine_template [-source=PATH] [-dest=PATH] [-exta_vars=PATH]`

## List command flags

`-source = PATH`

The path of source template from which we need to create 

## Dependencies

- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) - our loggin based on this package
- [github.com/hashicorp/vault/api](https://github.com/hashicorp/vault/tree/master/api) - Hashicorp Vault
    official library
- [gopkg.in/yaml.v2]() - Marshal/Unmarshal of yaml files

## Custom functions

You can easily add your own functions, which should help you with your issues.

Just write function and add mapping into 

`fmap := template.FuncMap{
 		"hashiVault": HashiVault,
 	}`
 	
 
### HashiVault

You can easily grab any field from Hashicorp Vault

#### Usage

    {{ "secret/customers/databases/example-com.password" | hashiVault }}
    
Argument is  `search string` and after dot is required `field name` and function name should be after pipe.  

## Build

Nothing special to build. Like any tool on Go language

### Windows

    GOOS=windows GOARCH=386 go build -o app_engine_template.exe
    GOOS=windows GOARCH=amd64 go build -o app_engine_template64.exe

### Linux

    GOOS=linux GOARCH=amd64 go build -o app_engine_template

### MacOS X

    GOOS=darwin GOARCH=amd64 go build -o app_engine_template

## Install

Put the binary to `/usr/local/bin` on *nix based machines. This will allow you to run 
it anywhere (based on your $PATH variables)

## Plan for future

- Add auto-scaling parameters
- We need to add ability to qualify your App Engine configuration
- Simplify the template. Maybe we will migrate to any other method to work with templates.
