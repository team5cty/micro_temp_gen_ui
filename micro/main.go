package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

/*Yaml has 4 fields
-Module - name of module to be generated
-Port - Port number
-Database - url and schema of tables
-Endpoints array - list of endpoints
*/

type Yaml struct {
	Module   string `yaml:"module"`
	Kafka    string `yaml:"kafka"`
	Port     string `yaml:"port"`
	Database struct {
		Provider string  `yaml:"provider"`
		Url      string  `yaml:"url"`
		Models   []Model `yaml:"models"`
	} `yaml:"database"`

	Endpoints []Endpoint `yaml:"endpoints"`
}

type Model struct {
	Table  string            `yaml:"table"`
	Schema map[string]string `yaml:"schema"`
}

type Endpoint struct {
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Method string `yaml:"method"`
	Kafka  struct {
		Topic string `yaml:"topic"`
		Type  string `yaml:"type"`
	} `yaml:"kafka"`
	Table string `yaml:"table"`
	Key   struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
	} `yaml:"key"`
	Json struct {
		Type       string            `yaml:"type"`
		Properties map[string]string `yaml:"properties"`
	} `yaml:"json"`
}

func main() {

	//read example.yaml and parse it into yamlobject
	yamlfile, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to read input.yaml file: %s\n", err.Error())
		return
	}
	yamlobject := &Yaml{}
	err = yaml.Unmarshal(yamlfile, yamlobject)
	if err != nil {
		fmt.Printf("Failed to parse YAML file: %s\n", err.Error())
	}

	/*Create folder structure:
	output_module:
		handlers: handler file for each endpoint.
		prisma : Database/prisma schema
	*/
	os.Mkdir("output", os.ModePerm)
	os.Chdir("output")
	os.Mkdir("handlers", os.ModePerm)
	os.Mkdir("prisma", os.ModePerm)
	os.Mkdir("kafka", os.ModePerm)

	//Run commands to initialise go module with dependencies
	cmd := exec.Command("go", "mod", "init", yamlobject.Module)
	execute_with_stdout_stderr(cmd)
	cmd = exec.Command("go", "get", "github.com/gorilla/mux")
	execute_with_stdout_stderr(cmd)
	cmd = exec.Command("go", "get", "github.com/steebchen/prisma-client-go")
	execute_with_stdout_stderr(cmd)
	cmd = exec.Command("go", "get", "github.com/joho/godotenv")
	execute_with_stdout_stderr(cmd)
	cmd = exec.Command("go", "get", "github.com/shopspring/decimal")
	execute_with_stdout_stderr(cmd)
	cmd = exec.Command("go", "get", "github.com/segmentio/kafka-go")
	execute_with_stdout_stderr(cmd)
	os.Chdir("..")

	//There are three template files:-
	//main - for generating main.go
	//handler - for generating files inside handler file for endpoints
	//prisma - for generating schema.prisma

	//read main template and generate main.go inside module

	template_file_buffer, err := os.ReadFile(filepath.Join("templates", "main"))
	if err != nil {
		fmt.Printf("Failed to read template.go file: %s\n", err.Error())
	}

	template_output_buffer, err := os.Create(filepath.Join("output", "main.go"))
	if err != nil {
		fmt.Printf("Failed to create output.go file: %s\n", err.Error())
	}

	t := template.Must(template.New("main.go_template").
		Funcs(template.FuncMap{"tolower": strings.ToLower}).
		Parse(string(template_file_buffer)))

	err = t.Execute(template_output_buffer, yamlobject)

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err.Error())
	}

	//read prisma template and generate schema.prisma inside prisma folder

	template_file_buffer, err = os.ReadFile(filepath.Join("templates", "prisma"))
	if err != nil {
		fmt.Printf("Failed to read template.go file: %s\n", err.Error())
	}

	template_output_buffer, err = os.Create(filepath.Join("output", "prisma", "schema.prisma"))
	if err != nil {
		fmt.Printf("Failed to create output.go file: %s\n", err.Error())
	}

	t = template.Must(template.New("prisma_template").Parse(string(template_file_buffer)))
	err = t.Execute(template_output_buffer, yamlobject)
	if err != nil {
		fmt.Printf("Error executing template: %s\n", err.Error())
	}

	// Run prisma db push to sync schema and db
	fmt.Println("Syncing schema with database...")
	os.Chdir(filepath.Join("output", "prisma"))
	cmd = exec.Command("go", "run", "github.com/steebchen/prisma-client-go", "db", "push")
	execute_with_stdout_stderr(cmd)
	os.Chdir("..")
	os.Chdir("..")

	//Loop through each endpoints and for each,
	//generate a handler file in handlers folder

	for _, i := range yamlobject.Endpoints {
		template_file_buffer, err := os.ReadFile(filepath.Join("templates", "handlers"))
		if err != nil {
			fmt.Printf("Failed to read template.go file: %s\n", err.Error())
		}
		template_output_buffer, err := os.Create(filepath.Join("output", "handlers", i.Name+".go"))
		if err != nil {
			fmt.Printf("Failed to create output.go file: %s\n", err.Error())
		}

		t := template.Must(template.New("handler_template").
			Funcs(template.FuncMap{
				"tolower":      strings.ToLower,
				"title":        strings.Title,
				"hasPathParam": hasPathParam,
			}).
			Parse(string(template_file_buffer)))

		//Passing module name along with endpoint
		data := map[string]any{
			"endpoint": i,
			"module":   yamlobject.Module,
		}

		err = t.Execute(template_output_buffer, data)
		if err != nil {
			fmt.Printf("Error executing template: %s\n", err.Error())
		}
	}

	template_file_buffer, err = os.ReadFile(filepath.Join("templates", "kafka"))
	if err != nil {
		fmt.Printf("Failed to read kafka file: %s\n", err.Error())
	}

	template_output_buffer, err = os.Create(filepath.Join("output", "kafka", "kafka.go"))
	if err != nil {
		fmt.Printf("Failed to create output.go file: %s\n", err.Error())
	}

	t = template.Must(template.New("kafka_template").
		Funcs(template.FuncMap{"tolower": strings.ToLower}).
		Parse(string(template_file_buffer)))

	err = t.Execute(template_output_buffer, yamlobject)

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err.Error())
	}
}

//User defined functions

func execute_with_stdout_stderr(cmd *exec.Cmd) {
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(stderr.String())
		return
	}
	fmt.Println(out.String())
}

func hasPathParam(path string) bool {
	return strings.Contains(path, "{")
}
