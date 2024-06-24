package configuration

import (
	"strings"

	"github.com/spf13/viper"
)

var envKeyReplacer *strings.Replacer = strings.NewReplacer(".", "_")

var (
	Application *viper.Viper
	Database    *viper.Viper
	LDAP        *viper.Viper
	Email       *viper.Viper
)

func Initialize() {
	Application = viper.New()
	Application.SetEnvPrefix("samuel")
	Application.SetEnvKeyReplacer(envKeyReplacer)
	Application.AutomaticEnv()
	Application.SetDefault("port", 5000)
	Application.SetDefault("connections", 3)

	Database = viper.New()
	Database.SetEnvPrefix("samuel_database")
	Database.SetEnvKeyReplacer(envKeyReplacer)
	Database.AutomaticEnv()
	Database.SetDefault("driver", "mysql")
	Database.SetDefault("port", 3306)

	LDAP = viper.New()
	LDAP.SetEnvPrefix("samuel_ldap")
	LDAP.SetEnvKeyReplacer(envKeyReplacer)
	LDAP.AutomaticEnv()
	LDAP.SetDefault("port", 389)

	Email = viper.New()
	Email.SetEnvPrefix("samuel_email")
	Email.SetEnvKeyReplacer(envKeyReplacer)
	Email.AutomaticEnv()
	Email.SetDefault("port", 587)
}
