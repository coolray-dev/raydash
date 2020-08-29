package setting

import (
	"strings"

	"github.com/coolray-dev/raydash/modules/log"
	"github.com/coolray-dev/raydash/modules/utils"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config is the actual instance of config
var Config *viper.Viper

func init() {
	Config = viper.New()
	pflag.String("config", "config", "config file name")
	pflag.Parse()
	Config.BindPFlags(pflag.CommandLine)
	if Config.IsSet("config") {
		Config.SetConfigFile(Config.GetString("config"))
	} else {
		Config.SetConfigName("config")
		Config.SetConfigType("yaml")
		Config.AddConfigPath(utils.AbsPath(""))
		Config.AddConfigPath("/etc/raydash")
	}

	Config.SetEnvPrefix("RAYDASH")
	Config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	Config.AutomaticEnv()

	if err := Config.ReadInConfig(); err != nil {
		log.Log.WithError(err).Fatal("Error Reading Config")
	}
}
