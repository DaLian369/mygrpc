package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mygrpc/proto"

	"gopkg.in/yaml.v2"
)

func Init(path string, cfg *proto.ConfigSt) (err error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("err ", err)
		return err
	}
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		fmt.Println("err ", err)
		return err
	}
	data, _ := json.Marshal(cfg)
	fmt.Println("config: ", string(data))
	return
}
