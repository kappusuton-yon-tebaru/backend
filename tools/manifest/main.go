package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

const DEPLOYMENT_DIRECTORY = "./deployment"

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	bs, err := os.ReadFile(".env")
	Must(err)

	envs := map[string]string{}
	for _, env := range strings.Split(string(bs), "\n") {
		if len(strings.TrimSpace(env)) > 0 {
			splitted := strings.SplitN(env, "=", 2)
			envs[splitted[0]] = splitted[1]
		}
	}

	entries, err := os.ReadDir(DEPLOYMENT_DIRECTORY)
	Must(err)

	rg, err := regexp.Compile("(<@[A-Z_]+>)")
	Must(err)

	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".example.yaml") {
			filename := strings.TrimSuffix(entry.Name(), ".example.yaml")
			local := fmt.Sprintf("%s.local.yaml", filename)

			bs, err := os.ReadFile(path.Join(DEPLOYMENT_DIRECTORY, entry.Name()))
			Must(err)

			content := rg.ReplaceAllStringFunc(string(bs), func(s string) string {
				key := strings.TrimSuffix(strings.TrimPrefix(s, "<@"), ">")
				val, ok := envs[key]
				if !ok {
					panic(fmt.Sprintf("Environment variable %s not found", key))
				}

				return val
			})

			target := path.Join(DEPLOYMENT_DIRECTORY, local)

			f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
			Must(err)
			defer f.Close()

			_, err = f.Write([]byte(content))
			Must(err)

			fmt.Println(target)
		}
	}
}
