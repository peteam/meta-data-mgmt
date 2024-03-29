package config

import (
	"os"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

/*
 * Auto generated by Quickstart. 
 * Developer to replace the generated test cases with real ones.
 */
func init() {
	v := viper.New()
	_env := os.Getenv("APP_ENV_PROFILE")
	if len(_env) == 0 {
		_env = "dev"
	}
	v.SetDefault("env", "_env")
	v.BindEnv("env")
	env := v.GetString("env")
	Viper = &config{initConfig(v, env)}
}

/*
 * Auto generated by Quickstart. 
 * Developer to replace the generated test cases with real ones.
 */
func Test_initConfig(t *testing.T) {
	type args struct {
		v   *viper.Viper
		env string
	}
	tests := []struct {
		name string
		args args
		want *viper.Viper
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initConfig(tt.args.v, tt.args.env); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
