package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// Generate generates code from all templates
func Generate(projectPath string, projectName string, repoPath string) error {
	var err error

	parameters := Parameters{
		ProjectName: projectName,
		RepoPath:    repoPath,
	}

	for _, element := range elements {
		if err = generateElement(projectPath, element, parameters); err != nil {
			return err
		}
	}

	// generate main separately, because of special file path
	main.FilePath = filepath.Join(main.FilePath, projectName)
	return generateElement(projectPath, main, parameters)
}

// generateElement generates code from one template
func generateElement(projectPath string, element Element, parameters Parameters) error {
	// build template
	generatorTemplate, err := template.New(fmt.Sprintf("%s template", element.FileName)).Parse(element.Template)
	if err != nil {
		return err
	}

	// create file directory
	if err := os.MkdirAll(filepath.Join(projectPath, element.FilePath), os.ModePerm); err != nil {
		return err
	}

	// create file
	file, err := os.Create(filepath.Join(projectPath, element.FilePath, element.FileName))
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	return generatorTemplate.Execute(file, parameters)
}
