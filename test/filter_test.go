package utils

import (
	"encoding/json"
	"fmt"
	flattener "github.com/anshal21/json-flattener"
	"github.com/go-playground/validator/v10"
	"testing"
)

type testFilter struct {
	Body struct {
		Ip    string `validate:"ip" json:"ip"`
		Owner string `validate:"eq=fundawang" json:"owner"`
	} `json:"body"`

	Header struct {
		Priority int `validate:"lte=5" json:"priority"`
	} `json:"headers"`
}

func TestFilter(t *testing.T) {
	validate := validator.New()
	err := validate.Var("fundawang", "eq=fundawan")
	fmt.Println(err)
	err = validate.Var("bluechio", "eq=bluechi")
	fmt.Println(err)
	err = validate.Var("169.59.160.68", "ip")
	fmt.Println(err)
	err = validate.Var("aaa", "len>3")
	fmt.Println(err)

}

func TestParse(t *testing.T) {
	mod_filter := `{
"body.build": "eq=93853"
}`

	json1 := `{
  "body": {
    "build": 93853,
    "chroot": "openeuler-22.03_LTS-x86_64",
    "copr": "libsys",
    "ip": "169.59.160.68",
    "owner": "fundawang",
    "pid": 3601173,
    "pkg": "bluechi",
    "status": 1,
    "user": "fundawang",
    "version": "0.8.0-0.202404120704.git77f8733",
    "what": "build end: user:packit copr:eclipse-bluechi-bluechi-872 build:7301454 pkg:bluechi version:0.8.0-0.202404120704.git77f8733 ip:169.59.160.68 pid:3601173 status:1",
    "who": "backend.worker-rpm_build_worker:7301454-fedora-rawhide-s390x"
  },
  "headers": {
    "fedora_messaging_schema": "copr.build.end",
    "fedora_messaging_severity": 20,
    "fedora_messaging_user_packit": true,
    "priority": 0,
    "sent-at": "2024-04-12T07:07:51+00:00"
  },
  "id": "243634a7-aa46-4c53-b669-f9d8366eb350",
  "queue": null,
  "topic": "org.fedoraproject.prod.copr.build.end"
}`

	map1 := make(map[string]interface{})
	mode_filter1 := make(map[string]string)

	json.Unmarshal([]byte(mod_filter), &mode_filter1)
	flattnedJSON, _ := flattener.FlattenJSON(json1, flattener.DotSeparator)
	json.Unmarshal([]byte(flattnedJSON), &map1)
	validate := validator.New()
	for s, i := range mode_filter1 {
		err := validate.Var(map1[s], i)
		if err != nil {
			break
		}
		fmt.Println(i)
		fmt.Println(err)
	}
}
