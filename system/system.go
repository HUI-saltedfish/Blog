package system

import (
	"os"

	"github.com/go-yaml/yaml"
)

type Configuration struct {
	SignupEnabled bool `yaml:"signup_enabled"`

	ALiYunAccessKey string `yaml:"aliyun_accesskey"`
	ALiYunSecretkey string `yaml:"aliyun_secretkey"`
	QiniuBucketName string `yaml:"aliyun_bucket_name"`
	ALiYunEndpoint  string `yaml:"aliyun_endpoint"`

	GithubClientID     string `yaml:"github_client_id"`
	GithubClientSecret string `yaml:"github_client_secret"`
	GithubAuthURL      string `yaml:"github_auth_url"`
	GithubRedirectURL  string `yaml:"github_redirect_url"`
	GithubTokenURL     string `yaml:"github_token_url"`
	GithubScope        string `yaml:"github_scope"`

	SmtpUsername string `yaml:"smtp_username"`
	SmtpPassword string `yaml:"smtp_password"`
	SmtpHost     string `yaml:"smtp_host"`

	SessSecret     string `yaml:"session_secret"`
	Domain         string `yaml:"domain"`
	Pubic          string `yaml:"public"`
	Addr           string `yaml:"addr"`
	BackupKey      string `yaml:"backup_key"`
	DSN            string `yaml:"dsn"`
	NotifyEmails   string `yaml:"notify_emails"`
	PageSize       int    `yaml:"page_size"`
	SmmsFileServer string `yaml:"smms_fileserver"`

	MaxShowTags     int `yaml:"max_show_tags"`
	MaxShowArchives int `yaml:"max_show_archives"`
	MaxShowRead     int `yaml:"max_show_read"`
	MaxShowComments int `yaml:"max_show_comments"`
	MaxShowLinks    int `yaml:"max_show_links"`
}

const (
	DEAULT_PAGE_SIZE = 10
)

var configuration *Configuration

func LoadConfiguration(configFilePath string) error {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	configuration = &Configuration{}
	err = yaml.Unmarshal(data, configuration)
	if err != nil {
		return err
	}

	if configuration.PageSize <= 0 {
		configuration.PageSize = DEAULT_PAGE_SIZE
	}

	return nil
}

func GetConfiguration() *Configuration {
	return configuration
}
