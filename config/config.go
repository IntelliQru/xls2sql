package config

import (
	"io/ioutil"
	"strings"
	"fmt"
	"bytes"
)

type configValue struct {
	line int
	value string
}

type Config struct {
	configFileName string
	lines          []string
	values         map[string]*configValue
}

func NewConfig() *Config  {
	return &Config{
		values: make(map[string]*configValue),
	}
}

func (c *Config) LoadConfig(fileName string) {

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	c.configFileName = fileName

	c.lines = strings.Split(string(b), "\n")
	for i := 0; i < len(c.lines); i++ {
		c.lines[i] = strings.TrimSpace(c.lines[i])
		if len(c.lines[i]) == 0 || c.lines[i][0] == ';' {
			continue
		}

		idx := strings.Index(c.lines[i], "=")
		if idx < 0 {
			panic(fmt.Sprintf("Invalid config syntax at line %d: %s", i, c.lines[i]))
		}

		key := strings.TrimSpace(c.lines[i][:idx])
		value := strings.TrimSpace(c.lines[i][idx+1:])
		c.values[key] = &configValue{
			line: i,
			value: value,
		}

		c.lines[i] = key + "=" + value

	}

}

func (c Config) GetWithFound(key string) (value string, isFound bool){
	val, b := c.values[key]

	isFound = b
	if b {
		value = val.value
	}
	return
}

func (c Config) Get(key string) (value string){
	val, b := c.values[key]

	if b {
		value = val.value
	}
	return
}

func (c *Config) Set(key, value string) {

	cv, isFound := c.values[key]
	var line int
	if isFound {
		line = cv.line
		cv.value = value
	} else {
		line = len(c.lines)
		c.lines = append(c.lines)
		c.values[key] = &configValue{
			line: line,
			value: value,
		}
	}

	c.lines[line] = key + "=" + value

}

func (c Config) Save() {

	b := bytes.NewBuffer(nil)

	bFirst := true

	for _, line := range c.lines {
		if !bFirst {
			b.WriteString("\n")
		}
		b.WriteString(line)
	}

	err := ioutil.WriteFile(c.configFileName, b.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
}