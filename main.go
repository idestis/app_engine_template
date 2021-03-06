package main

import (
	"bytes"
	"flag"
	vaultApi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

type Vault struct {
	Result string
}

type References struct {
	// TODO: Add autoscalling variables
	Env           string `yaml:"env"`
	Service       string `yaml:"service"`
	Threadsafe    bool   `yaml:"threadsafe";template:"threadsafe"`
	Runtime       string `yaml:"runtime"`
	ApiVersion    string `yaml:"api_version"`
	RuntimeConfig []struct {
		DocumentRoot string `yaml:"document_root"`
	} `yaml:"runtime_config"`
	Handlers []struct {
		Url                      string            `yaml:"url"`
		Script                   string            `yaml:"script"`
		Secure                   string            `yaml:"secure"`
		RedirectHttpResponseCode int               `yaml:"redirect_http_response_code"`
		StaticDir                string            `yaml:"static_dir"`
		HttpHeaders              map[string]string `yaml:"http_headers"`
		Upload                   string            `yaml:"upload"`
		Login                    string            `yaml:"login"`
		AuthFailAction           string            `yaml:"auth_fail_action"`
		StaticFiles              string            `yaml:"static_files"`
		Expiration               string            `yaml:"expiration"`
		MimeType                 string            `yaml:"mime_type"`
	} `yaml:"handlers"`
	EnvVariables  map[string]string `yaml:"env_variables"`
	InstanceClass string            `yaml:"instance_class"`
	ErrorHandlers []struct {
		File      string `yaml:"file"`
		ErrorCode string `yaml:"error_code"`
	} `yaml:"error_handlers"`
	DefaultExpiration string `yaml:"default_expiration"`
	BetaSettings      struct {
		CloudSQLInstance string `yaml:"cloud_sql_instances"`
	} `yaml:"beta_settings"`
	Builtins []struct {
		Appstats  string `yaml:"appstats"`
		Deferred  string `yaml:"deferred"`
		RemoteApi string `yaml:"remote_api"`
	} `yaml:"builtins"`
	Includes        []string `yaml:"includes"`
	InboundServices []string `yaml:"inbound_services"`
	Libraries       []struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"libraries"`
	SkipFiles        []string `yaml:"skip_files"`
	Version          string   `yaml:"version"`
	AutomaticScaling []struct {
		TargetCpuUtilization        string `yaml:"target_cpu_utilization"`
		TargetThroughputUtilization string `yaml:"target_throughput_utilization"`
		MinInstances                int    `yaml:"min_instances"`
		MaxInstances                int    `yaml:"max_instances"`
		MinPendingLatency           string `yaml:"min_pending_latency"`
		MaxPendingLatency           string `yaml:"max_pending_latency"`
		MaxConcurrentRequests       int    `yaml:"max_concurrent_requests"`
		MinIdleInstances            int    `yaml:"min_idle_instances"`
		MaxIdleInstances            int    `yaml:"max_idle_instances"`
	} `yaml:"automatic_scaling"`
	BasicScaling []struct {
		MaxInstances int    `yaml:"max_instances"`
		IdleTimeout  string `yaml:"idle_timeout"`
	} `yaml:"basic_scaling"`
	ManualScaling []struct {
		Instances int `yaml:"instances"`
	} `yaml:"manual_scaling"`
}

func HashiVault(request string) interface{} {
	vaultAddr := os.Getenv("VAULT_ADDR")
	vaultToken := os.Getenv("VAULT_TOKEN")
	client, err := vaultApi.NewClient(&vaultApi.Config{
		Address: vaultAddr,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	keyfield := strings.Split(request, ".")
	client.SetToken(vaultToken)
	secretValues, err := client.Logical().Read(keyfield[0])
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	dataMap := secretValues.Data
	return dataMap[keyfield[1]]
}

func main() {
	source := flag.String("source", "", "Path to source template")
	dest := flag.String("dest", "./app.yaml", "Filename how to save output")
	extra := flag.String("extra_vars", "", "Path to file with extra variables")
	flag.Parse()

	if *source == "" {
		log.Error("-source is required to start")
		os.Exit(1)
	}

	fmap := template.FuncMap{
		"hashiVault": HashiVault,
	}
	var Vars References
	filename := path.Base(*source)
	yamlTemplate := template.Must(template.New(filename).Funcs(fmap).ParseFiles(*source))
	if *extra != "" {
		var tpl bytes.Buffer
		varsFile := path.Base(*extra)
		variables := template.Must(template.New(varsFile).Funcs(fmap).ParseFiles(*extra))
		err := variables.Execute(&tpl, "no data needed")
		vars, _ := ioutil.ReadAll(&tpl)
		err = yaml.Unmarshal(vars, &Vars)

		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
	f, err := os.Create(*dest)
	if err != nil {
		log.Error("Create file: ", err)
		return
	}
	err = yamlTemplate.Execute(f, Vars)
	if err != nil {
		log.Error("Execute template: ", err)
		return
	}
	f.Close()
}
