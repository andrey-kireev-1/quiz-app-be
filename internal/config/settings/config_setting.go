package settings

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"reflect"
	"strings"
)

type ConfigEr interface {
	SetDefault() error
}

const (
	TimeFormatGMDHISM = "2006-01-02 15:04:05.999999-07:00"
	TimeFormatGMDHIS  = "2006-01-02 15:04:05"
	TimeFormatDMGHIS  = "02.01.2006 15:04:05"
	TimeFormatHIS     = "15:04:05"
	TimeFormatDMG     = "02.01.2006"
)

func GetDefault(cfg ConfigEr, fileConf, envPrefix string) error {
	if err := Get(cfg, fileConf, envPrefix); err != nil {
		return err
	}
	return cfg.SetDefault()
}

func Get(cfg interface{}, fileConf, envPrefix string) error {
	vip := viper.New()
	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if envPrefix != "" {
		vip.SetEnvPrefix(envPrefix)
	}
	bindEnvs(vip, cfg)
	if fileConf != "" {
		// search config
		d, _ := os.Getwd()
		dd := string(os.PathSeparator) + "etc" + string(os.PathSeparator)
		for strings.Count(d, string(os.PathSeparator)) > 1 {
			if _, err := os.Stat(d + dd + fileConf); err != nil {
				d = path.Dir(d)
			} else {
				break
			}
		}
		fileConf = d + dd + fileConf
		// yaml
		data, err := os.ReadFile(fileConf)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return err
		}
		vip.SetConfigFile(fileConf)
		if err := vip.ReadInConfig(); err != nil {
			return err
		}
	}
	return vip.Unmarshal(cfg)
}

func bindEnvs(cfg *viper.Viper, cfgStruct interface{}, parts ...string) {
	ifv := reflect.ValueOf(cfgStruct)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}
	for i := 0; i < ifv.NumField(); i++ {
		v := ifv.Field(i)
		t := ifv.Type().Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			if tv, ok = t.Tag.Lookup("yaml"); !ok {
				continue
			}
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(cfg, v.Interface(), append(parts, tv)...)
		default:
			envVar := strings.Join(append(parts, tv), ".")
			err := cfg.BindEnv(envVar)
			if err != nil {
				panic(err)
			}
		}
	}
}
