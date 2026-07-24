package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DBHost       string `envconfig:"DB_HOST"       required:"true"`
	DBPort       int    `envconfig:"DB_PORT"       default:"3306"`
	DBUser       string `envconfig:"DB_USER"       required:"true"`
	DBPassword   string `envconfig:"DB_PASSWORD"   required:"true"`
	DBName       string `envconfig:"DB_NAME"       required:"true"`
	R2AccountID        string `envconfig:"R2_ACCOUNT_ID"        required:"true"`
	R2AccessKeyID      string `envconfig:"R2_ACCESS_KEY_ID"     required:"true"`
	R2SecretAccessKey  string `envconfig:"R2_SECRET_ACCESS_KEY" required:"true"`
	R2BucketName       string `envconfig:"R2_BUCKET_NAME"       required:"true"`
	R2PublicBaseURL    string `envconfig:"R2_PUBLIC_BASE_URL"   required:"true"`
	JWTSecret          string `envconfig:"JWT_SECRET"           required:"true"`
	BattleTokenSecret  string `envconfig:"BATTLE_TOKEN_SECRET"  required:"true"`
	Port               int    `envconfig:"PORT"                 default:"8085"`
}

func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
