package cli

type Config struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	AccessToken string `yaml:"access_token"`
	Url         string `yaml:"url"`
}
