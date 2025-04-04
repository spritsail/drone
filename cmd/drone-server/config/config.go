// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/dustin/go-humanize"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// IMPORTANT please do not add new configuration parameters unless it has
// been discussed on the mailing list. We are attempting to reduce the
// number of configuration parameters, and may reject pull requests that
// introduce new parameters. (mailing list https://community.harness.io)

// default runner hostname.
var hostname string

func init() {
	hostname, _ = os.Hostname()
	if hostname == "" {
		hostname = "localhost"
	}
}

type (
	// Config provides the system configuration.
	Config struct {
		License string `envconfig:"DRONE_LICENSE"`

		Authn        Authentication
		Agent        Agent
		Convert      Convert
		Cleanup      Cleanup
		Cron         Cron
		Cloning      Cloning
		Database     Database
		Docker       Docker
		HTTP         HTTP
		Jsonnet      Jsonnet
		Starlark     Starlark
		Logging      Logging
		Prometheus   Prometheus
		Proxy        Proxy
		Redis        Redis
		Registration Registration
		Registries   Registries
		Repository   Repository
		Runner       Runner
		RPC          RPC
		S3           S3
		Secrets      Secrets
		Server       Server
		Session      Session
		Status       Status
		Users        Users
		Validate     Validate
		Webhook      Webhook
		Yaml         Yaml

		// Remote configurations
		Bitbucket       Bitbucket
		Gitea           Gitea
		Github          Github
		GitLab          GitLab
		Gogs            Gogs
		Stash           Stash
		Gitee           Gitee
		IncomingWebhook IncomingWebhook
	}

	// Cloning provides the cloning configuration.
	Cloning struct {
		AlwaysAuth bool   `envconfig:"DRONE_GIT_ALWAYS_AUTH"`
		Username   string `envconfig:"DRONE_GIT_USERNAME"`
		Password   string `envconfig:"DRONE_GIT_PASSWORD"`
		Image      string `envconfig:"DRONE_GIT_IMAGE"`
		Pull       string `envconfig:"DRONE_GIT_IMAGE_PULL" default:"IfNotExists"`
	}

	Cleanup struct {
		Disabled bool          `envconfig:"DRONE_CLEANUP_DISABLED"`
		Interval time.Duration `envconfig:"DRONE_CLEANUP_INTERVAL"         default:"24h"`
		Running  time.Duration `envconfig:"DRONE_CLEANUP_DEADLINE_RUNNING" default:"24h"`
		Pending  time.Duration `envconfig:"DRONE_CLEANUP_DEADLINE_PENDING" default:"24h"`
		Buffer   time.Duration `envconfig:"DRONE_CLEANUP_BUFFER" default:"30m"`
	}

	// Cron provides the cron configuration.
	Cron struct {
		Disabled bool          `envconfig:"DRONE_CRON_DISABLED"`
		Interval time.Duration `envconfig:"DRONE_CRON_INTERVAL" default:"30m"`
	}

	// Database provides the database configuration.
	Database struct {
		Driver         string `envconfig:"DRONE_DATABASE_DRIVER"          default:"sqlite3"`
		Datasource     string `envconfig:"DRONE_DATABASE_DATASOURCE"      default:"core.sqlite"`
		Secret         string `envconfig:"DRONE_DATABASE_SECRET"`
		MaxConnections int    `envconfig:"DRONE_DATABASE_MAX_CONNECTIONS" default:"0"`

		// Feature flag
		LegacyBatch bool `envconfig:"DRONE_DATABASE_LEGACY_BATCH"`

		// Feature flag
		EncryptUserTable    bool `envconfig:"DRONE_DATABASE_ENCRYPT_USER_TABLE"`
		EncryptMixedContent bool `envconfig:"DRONE_DATABASE_ENCRYPT_MIXED_MODE"`
	}

	// Docker provides docker configuration
	Docker struct {
		Config string `envconfig:"DRONE_DOCKER_CONFIG"`
	}

	// Jsonnet configures the jsonnet plugin
	Jsonnet struct {
		Enabled     bool `envconfig:"DRONE_JSONNET_ENABLED"`
		ImportLimit int  `envconfig:"DRONE_JSONNET_IMPORT_LIMIT" default:"0"`
	}

	// Starlark configures the starlark plugin
	Starlark struct {
		Enabled   bool   `envconfig:"DRONE_STARLARK_ENABLED"`
		StepLimit uint64 `envconfig:"DRONE_STARLARK_STEP_LIMIT"`
		SizeLimit uint64 `envconfig:"DRONE_STARLARK_SIZE_LIMIT" default:"0"`
	}

	// License provides license configuration
	License struct {
		Key      string `envconfig:"DRONE_LICENSE"`
		Endpoint string `envconfig:"DRONE_LICENSE_ENDPOINT"`
	}

	// Logging provides the logging configuration.
	Logging struct {
		Debug  bool `envconfig:"DRONE_LOGS_DEBUG"`
		Trace  bool `envconfig:"DRONE_LOGS_TRACE"`
		Color  bool `envconfig:"DRONE_LOGS_COLOR"`
		Pretty bool `envconfig:"DRONE_LOGS_PRETTY"`
		Text   bool `envconfig:"DRONE_LOGS_TEXT"`
	}

	// Prometheus provides the prometheus configuration.
	Prometheus struct {
		EnableAnonymousAccess bool `envconfig:"DRONE_PROMETHEUS_ANONYMOUS_ACCESS" default:"false"`
		EnableHTTPMetrics     bool `envconfig:"DRONE_PROMETHEUS_HTTP_METRICS" default:"false"`
	}

	// Redis provides the redis configuration.
	Redis struct {
		ConnectionString string `envconfig:"DRONE_REDIS_CONNECTION"`
		Addr             string `envconfig:"DRONE_REDIS_ADDR"`
		Password         string `envconfig:"DRONE_REDIS_PASSWORD"`
		DB               int    `envconfig:"DRONE_REDIS_DB"`
	}

	// Repository provides the repository configuration.
	Repository struct {
		Filter     []string `envconfig:"DRONE_REPOSITORY_FILTER"`
		Visibility string   `envconfig:"DRONE_REPOSITORY_VISIBILITY"`
		Trusted    bool     `envconfig:"DRONE_REPOSITORY_TRUSTED"`

		// THIS SETTING IS INTERNAL USE ONLY AND SHOULD
		// NOT BE USED OR RELIED UPON IN PRODUCTION.
		Ignore []string `envconfig:"DRONE_REPOSITORY_IGNORE"`
	}

	// Registries provides the registry configuration.
	Registries struct {
		Endpoint   string `envconfig:"DRONE_REGISTRY_ENDPOINT"`
		Password   string `envconfig:"DRONE_REGISTRY_SECRET"`
		SkipVerify bool   `envconfig:"DRONE_REGISTRY_SKIP_VERIFY"`
	}

	// Secrets provides the secret configuration.
	Secrets struct {
		Endpoint   string `envconfig:"DRONE_SECRET_ENDPOINT"`
		Password   string `envconfig:"DRONE_SECRET_SECRET"`
		SkipVerify bool   `envconfig:"DRONE_SECRET_SKIP_VERIFY"`
	}

	// RPC provides the rpc configuration.
	RPC struct {
		Server string `envconfig:"DRONE_RPC_SERVER"`
		Secret string `envconfig:"DRONE_RPC_SECRET"`
		Debug  bool   `envconfig:"DRONE_RPC_DEBUG"`
		Host   string `envconfig:"DRONE_RPC_HOST"`
		Proto  string `envconfig:"DRONE_RPC_PROTO"`
		// Hosts  map[string]string `envconfig:"DRONE_RPC_EXTRA_HOSTS"`
	}

	Agent struct {
		Disabled bool `envconfig:"DRONE_AGENTS_DISABLED"`
	}

	// Runner provides the runner configuration.
	Runner struct {
		Local      bool              `envconfig:"DRONE_RUNNER_LOCAL"`
		Image      string            `envconfig:"DRONE_RUNNER_IMAGE"    default:"drone/controller:1"`
		Platform   string            `envconfig:"DRONE_RUNNER_PLATFORM" default:"linux/amd64"`
		OS         string            `envconfig:"DRONE_RUNNER_OS"`
		Arch       string            `envconfig:"DRONE_RUNNER_ARCH"`
		Kernel     string            `envconfig:"DRONE_RUNNER_KERNEL"`
		Variant    string            `envconfig:"DRONE_RUNNER_VARIANT"`
		Machine    string            `envconfig:"DRONE_RUNNER_NAME"`
		Capacity   int               `envconfig:"DRONE_RUNNER_CAPACITY" default:"2"`
		Labels     map[string]string `envconfig:"DRONE_RUNNER_LABELS"`
		Volumes    []string          `envconfig:"DRONE_RUNNER_VOLUMES"`
		Networks   []string          `envconfig:"DRONE_RUNNER_NETWORKS"`
		Devices    []string          `envconfig:"DRONE_RUNNER_DEVICES"`
		Privileged []string          `envconfig:"DRONE_RUNNER_PRIVILEGED_IMAGES"`
		Environ    map[string]string `envconfig:"DRONE_RUNNER_ENVIRON"`
		Limits     struct {
			MemSwapLimit Bytes  `envconfig:"DRONE_LIMIT_MEM_SWAP"`
			MemLimit     Bytes  `envconfig:"DRONE_LIMIT_MEM"`
			ShmSize      Bytes  `envconfig:"DRONE_LIMIT_SHM_SIZE"`
			CPUQuota     int64  `envconfig:"DRONE_LIMIT_CPU_QUOTA"`
			CPUShares    int64  `envconfig:"DRONE_LIMIT_CPU_SHARES"`
			CPUSet       string `envconfig:"DRONE_LIMIT_CPU_SET"`
		}
	}

	// Server provides the server configuration.
	Server struct {
		Addr  string `envconfig:"-"`
		Host  string `envconfig:"DRONE_SERVER_HOST" default:"localhost:8080"`
		Port  string `envconfig:"DRONE_SERVER_PORT" default:":8080"`
		Proto string `envconfig:"DRONE_SERVER_PROTO" default:"http"`
		Pprof bool   `envconfig:"DRONE_PPROF_ENABLED"`
		Acme  bool   `envconfig:"DRONE_TLS_AUTOCERT"`
		Email string `envconfig:"DRONE_TLS_EMAIL"`
		Cert  string `envconfig:"DRONE_TLS_CERT"`
		Key   string `envconfig:"DRONE_TLS_KEY"`
	}

	// Proxy provides proxy server configuration.
	Proxy struct {
		Addr  string `envconfig:"-"`
		Host  string `envconfig:"DRONE_SERVER_PROXY_HOST"`
		Proto string `envconfig:"DRONE_SERVER_PROXY_PROTO"`
	}

	// Registration configuration.
	Registration struct {
		Closed bool `envconfig:"DRONE_REGISTRATION_CLOSED"`
	}

	// Authentication Controller configuration
	Authentication struct {
		Endpoint   string `envconfig:"DRONE_ADMISSION_PLUGIN_ENDPOINT"`
		Secret     string `envconfig:"DRONE_ADMISSION_PLUGIN_SECRET"`
		SkipVerify bool   `envconfig:"DRONE_ADMISSION_PLUGIN_SKIP_VERIFY"`
	}

	// Session provides the session configuration.
	Session struct {
		Timeout time.Duration `envconfig:"DRONE_COOKIE_TIMEOUT" default:"720h"`
		Secret  string        `envconfig:"DRONE_COOKIE_SECRET"`
		Secure  bool          `envconfig:"DRONE_COOKIE_SECURE"`
	}

	// Status provides status configurations.
	Status struct {
		Disabled bool   `envconfig:"DRONE_STATUS_DISABLED"`
		Name     string `envconfig:"DRONE_STATUS_NAME"`
	}

	// Users provides the user configuration.
	Users struct {
		Create UserCreate    `envconfig:"DRONE_USER_CREATE"`
		Filter []string      `envconfig:"DRONE_USER_FILTER"`
		MinAge time.Duration `envconfig:"DRONE_MIN_AGE"`
	}

	// Webhook provides the webhook configuration.
	Webhook struct {
		Events     []string `envconfig:"DRONE_WEBHOOK_EVENTS"`
		Endpoint   []string `envconfig:"DRONE_WEBHOOK_ENDPOINT"`
		Secret     string   `envconfig:"DRONE_WEBHOOK_SECRET"`
		SkipVerify bool     `envconfig:"DRONE_WEBHOOK_SKIP_VERIFY"`
	}

	// Yaml provides the yaml webhook configuration.
	Yaml struct {
		Endpoint   string        `envconfig:"DRONE_YAML_ENDPOINT"`
		Secret     string        `envconfig:"DRONE_YAML_SECRET"`
		SkipVerify bool          `envconfig:"DRONE_YAML_SKIP_VERIFY"`
		Timeout    time.Duration `envconfig:"DRONE_YAML_TIMEOUT" default:"1m"`
	}

	// Convert provides the converter webhook configuration.
	Convert struct {
		Extension  string        `envconfig:"DRONE_CONVERT_PLUGIN_EXTENSION"`
		Endpoint   string        `envconfig:"DRONE_CONVERT_PLUGIN_ENDPOINT"`
		Secret     string        `envconfig:"DRONE_CONVERT_PLUGIN_SECRET"`
		SkipVerify bool          `envconfig:"DRONE_CONVERT_PLUGIN_SKIP_VERIFY"`
		CacheSize  int           `envconfig:"DRONE_CONVERT_PLUGIN_CACHE_SIZE" default:"10"`
		Timeout    time.Duration `envconfig:"DRONE_CONVERT_TIMEOUT" default:"1m"`

		// this flag can be removed once we solve for
		// https://github.com/harness/drone/pull/2994#issuecomment-795955312
		Multi bool `envconfig:"DRONE_CONVERT_MULTI"`
	}

	// Validate provides the validation webhook configuration.
	Validate struct {
		Endpoint   string        `envconfig:"DRONE_VALIDATE_PLUGIN_ENDPOINT"`
		Secret     string        `envconfig:"DRONE_VALIDATE_PLUGIN_SECRET"`
		SkipVerify bool          `envconfig:"DRONE_VALIDATE_PLUGIN_SKIP_VERIFY"`
		Timeout    time.Duration `envconfig:"DRONE_VALIDATE_TIMEOUT" default:"1m"`
	}

	//
	// Source code management.
	//

	// Bitbucket provides the bitbucket client configuration.
	Bitbucket struct {
		ClientID     string `envconfig:"DRONE_BITBUCKET_CLIENT_ID"`
		ClientSecret string `envconfig:"DRONE_BITBUCKET_CLIENT_SECRET"`
		SkipVerify   bool   `envconfig:"DRONE_BITBUCKET_SKIP_VERIFY"`
		Debug        bool   `envconfig:"DRONE_BITBUCKET_DEBUG"`
	}

	// Gitea provides the gitea client configuration.
	Gitea struct {
		Server       string   `envconfig:"DRONE_GITEA_SERVER"`
		ClientID     string   `envconfig:"DRONE_GITEA_CLIENT_ID"`
		ClientSecret string   `envconfig:"DRONE_GITEA_CLIENT_SECRET"`
		RedirectURL  string   `envconfig:"DRONE_GITEA_REDIRECT_URL"`
		SkipVerify   bool     `envconfig:"DRONE_GITEA_SKIP_VERIFY"`
		Scope        []string `envconfig:"DRONE_GITEA_SCOPE" default:"repo,repo:status,user:email,read:org"`
		Debug        bool     `envconfig:"DRONE_GITEA_DEBUG"`
	}

	// Github provides the github client configuration.
	Github struct {
		Server       string   `envconfig:"DRONE_GITHUB_SERVER" default:"https://github.com"`
		APIServer    string   `envconfig:"DRONE_GITHUB_API_SERVER"`
		ClientID     string   `envconfig:"DRONE_GITHUB_CLIENT_ID"`
		ClientSecret string   `envconfig:"DRONE_GITHUB_CLIENT_SECRET"`
		SkipVerify   bool     `envconfig:"DRONE_GITHUB_SKIP_VERIFY"`
		Scope        []string `envconfig:"DRONE_GITHUB_SCOPE" default:"repo,repo:status,user:email,read:org"`
		RateLimit    int      `envconfig:"DRONE_GITHUB_USER_RATELIMIT"`
		Debug        bool     `envconfig:"DRONE_GITHUB_DEBUG"`
	}

	// Gitee providers the gitee client configuration.
	Gitee struct {
		Server       string   `envconfig:"DRONE_GITEE_SERVER" default:"https://gitee.com"`
		APIServer    string   `envconfig:"DRONE_GITEE_API_SERVER" default:"https://gitee.com/api/v5"`
		ClientID     string   `envconfig:"DRONE_GITEE_CLIENT_ID"`
		ClientSecret string   `envconfig:"DRONE_GITEE_CLIENT_SECRET"`
		RedirectURL  string   `envconfig:"DRONE_GITEE_REDIRECT_URL"`
		SkipVerify   bool     `envconfig:"DRONE_GITEE_SKIP_VERIFY"`
		Scope        []string `envconfig:"DRONE_GITEE_SCOPE" default:"user_info,projects,pull_requests,hook"`
		Debug        bool     `envconfig:"DRONE_GITEE_DEBUG"`
	}

	// GitLab provides the gitlab client configuration.
	GitLab struct {
		Server       string `envconfig:"DRONE_GITLAB_SERVER" default:"https://gitlab.com"`
		ClientID     string `envconfig:"DRONE_GITLAB_CLIENT_ID"`
		ClientSecret string `envconfig:"DRONE_GITLAB_CLIENT_SECRET"`
		SkipVerify   bool   `envconfig:"DRONE_GITLAB_SKIP_VERIFY"`
		Debug        bool   `envconfig:"DRONE_GITLAB_DEBUG"`
	}

	// Gogs provides the gogs client configuration.
	Gogs struct {
		Server     string `envconfig:"DRONE_GOGS_SERVER"`
		SkipVerify bool   `envconfig:"DRONE_GOGS_SKIP_VERIFY"`
		Debug      bool   `envconfig:"DRONE_GOGS_DEBUG"`
	}

	// Stash provides the stash client configuration.
	Stash struct {
		Server         string `envconfig:"DRONE_STASH_SERVER"`
		ConsumerKey    string `envconfig:"DRONE_STASH_CONSUMER_KEY"`
		ConsumerSecret string `envconfig:"DRONE_STASH_CONSUMER_SECRET"`
		PrivateKey     string `envconfig:"DRONE_STASH_PRIVATE_KEY"`
		SkipVerify     bool   `envconfig:"DRONE_STASH_SKIP_VERIFY"`
		Debug          bool   `envconfig:"DRONE_STASH_DEBUG"`
	}

	// S3 provides the storage configuration.
	S3 struct {
		Bucket    string `envconfig:"DRONE_S3_BUCKET"`
		Prefix    string `envconfig:"DRONE_S3_PREFIX"`
		Endpoint  string `envconfig:"DRONE_S3_ENDPOINT"`
		PathStyle bool   `envconfig:"DRONE_S3_PATH_STYLE"`
	}

	// HTTP provides http configuration.
	HTTP struct {
		AllowedHosts          []string          `envconfig:"DRONE_HTTP_ALLOWED_HOSTS"`
		HostsProxyHeaders     []string          `envconfig:"DRONE_HTTP_PROXY_HEADERS"`
		SSLRedirect           bool              `envconfig:"DRONE_HTTP_SSL_REDIRECT"`
		SSLTemporaryRedirect  bool              `envconfig:"DRONE_HTTP_SSL_TEMPORARY_REDIRECT"`
		SSLHost               string            `envconfig:"DRONE_HTTP_SSL_HOST"`
		SSLProxyHeaders       map[string]string `envconfig:"DRONE_HTTP_SSL_PROXY_HEADERS"`
		STSSeconds            int64             `envconfig:"DRONE_HTTP_STS_SECONDS"`
		STSIncludeSubdomains  bool              `envconfig:"DRONE_HTTP_STS_INCLUDE_SUBDOMAINS"`
		STSPreload            bool              `envconfig:"DRONE_HTTP_STS_PRELOAD"`
		ForceSTSHeader        bool              `envconfig:"DRONE_HTTP_STS_FORCE_HEADER"`
		BrowserXSSFilter      bool              `envconfig:"DRONE_HTTP_BROWSER_XSS_FILTER"    default:"true"`
		FrameDeny             bool              `envconfig:"DRONE_HTTP_FRAME_DENY"            default:"true"`
		ContentTypeNosniff    bool              `envconfig:"DRONE_HTTP_CONTENT_TYPE_NO_SNIFF"`
		ContentSecurityPolicy string            `envconfig:"DRONE_HTTP_CONTENT_SECURITY_POLICY"`
		ReferrerPolicy        string            `envconfig:"DRONE_HTTP_REFERRER_POLICY"`
	}

	IncomingWebhook struct {
		Events []string `envconfig:"DRONE_INCOMING_WEBHOOK_EVENTS" default:"branch,deployment,push,tag,pull_request"`
	}
)

// Environ returns the settings from the environment.
func Environ() (Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	defaultAddress(&cfg)
	defaultProxy(&cfg)
	defaultRunner(&cfg)
	defaultSession(&cfg)
	defaultCallback(&cfg)
	configureGithub(&cfg)
	if err := kubernetesServiceConflict(&cfg); err != nil {
		return cfg, err
	}
	return cfg, err
}

// String returns the configuration in string format.
func (c *Config) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

// IsGitHub returns true if the GitHub integration
// is activated.
func (c *Config) IsGitHub() bool {
	return c.Github.ClientID != ""
}

// IsGitHubEnterprise returns true if the GitHub
// integration is activated.
func (c *Config) IsGitHubEnterprise() bool {
	return c.IsGitHub() && !strings.HasPrefix(c.Github.Server, "https://github.com")
}

// IsGitLab returns true if the GitLab integration
// is activated.
func (c *Config) IsGitLab() bool {
	return c.GitLab.ClientID != ""
}

// IsGogs returns true if the Gogs integration
// is activated.
func (c *Config) IsGogs() bool {
	return c.Gogs.Server != ""
}

// IsGitea returns true if the Gitea integration
// is activated.
func (c *Config) IsGitea() bool {
	return c.Gitea.Server != ""
}

// IsGitee returns true if the Gitee integration
// is activated.
func (c *Config) IsGitee() bool {
	return c.Gitee.ClientID != ""
}

// IsBitbucket returns true if the Bitbucket Cloud
// integration is activated.
func (c *Config) IsBitbucket() bool {
	return c.Bitbucket.ClientID != ""
}

// IsStash returns true if the Atlassian Stash
// integration is activated.
func (c *Config) IsStash() bool {
	return c.Stash.Server != ""
}

func cleanHostname(hostname string) string {
	hostname = strings.ToLower(hostname)
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.TrimPrefix(hostname, "https://")

	return hostname
}

func defaultAddress(c *Config) {
	if c.Server.Key != "" || c.Server.Cert != "" || c.Server.Acme {
		c.Server.Port = ":443"
		c.Server.Proto = "https"
	}
	c.Server.Host = cleanHostname(c.Server.Host)
	c.Server.Addr = c.Server.Proto + "://" + c.Server.Host
}

func defaultProxy(c *Config) {
	if c.Proxy.Host == "" {
		c.Proxy.Host = c.Server.Host
	} else {
		c.Proxy.Host = cleanHostname(c.Proxy.Host)
	}
	if c.Proxy.Proto == "" {
		c.Proxy.Proto = c.Server.Proto
	}
	c.Proxy.Addr = c.Proxy.Proto + "://" + c.Proxy.Host
}

func defaultCallback(c *Config) {
	if c.RPC.Host == "" {
		c.RPC.Host = c.Server.Host
	}
	if c.RPC.Proto == "" {
		c.RPC.Proto = c.Server.Proto
	}
}

func defaultRunner(c *Config) {
	if c.Runner.Machine == "" {
		c.Runner.Machine = hostname
	}
	parts := strings.Split(c.Runner.Platform, "/")
	if len(parts) == 2 && c.Runner.OS == "" {
		c.Runner.OS = parts[0]
	}
	if len(parts) == 2 && c.Runner.Arch == "" {
		c.Runner.Arch = parts[1]
	}
}

func defaultSession(c *Config) {
	if c.Session.Secret == "" {
		c.Session.Secret = uniuri.NewLen(32)
	}
}

func configureGithub(c *Config) {
	if c.Github.APIServer != "" {
		return
	}
	if c.Github.Server == "https://github.com" {
		c.Github.APIServer = "https://api.github.com"
	} else {
		c.Github.APIServer = strings.TrimSuffix(c.Github.Server, "/") + "/api/v3"
	}
}

func kubernetesServiceConflict(c *Config) error {
	if strings.HasPrefix(c.Server.Port, "tcp://") {
		return errors.New("Invalid port configuration. See https://community.harness.io/t/drone-server-changing-ports-protocol/11400")
	}
	return nil
}

// Bytes stores number bytes (e.g. megabytes)
type Bytes int64

// Decode implements a decoder that parses a string representation
// of bytes into the number of bytes it represents.
func (b *Bytes) Decode(value string) error {
	v, err := humanize.ParseBytes(value)
	*b = Bytes(v)
	return err
}

// Int64 returns the int64 value of the Byte.
func (b *Bytes) Int64() int64 {
	return int64(*b)
}

// String returns the string value of the Byte.
func (b *Bytes) String() string {
	return fmt.Sprint(*b)
}

// UserCreate stores account information used to bootstrap
// the admin user account when the system initializes.
type UserCreate struct {
	Username string
	Machine  bool
	Admin    bool
	Token    string
}

// Decode implements a decoder that extracts user information
// from the environment variable string.
func (u *UserCreate) Decode(value string) error {
	for _, param := range strings.Split(value, ",") {
		parts := strings.Split(param, ":")
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		val := parts[1]
		switch key {
		case "username":
			u.Username = val
		case "token":
			u.Token = val
		case "machine":
			u.Machine = val == "true"
		case "admin":
			u.Admin = val == "true"
		}
	}
	return nil
}
