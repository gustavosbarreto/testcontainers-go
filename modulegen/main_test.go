package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/testcontainers/testcontainers-go/modulegen/internal/dependabot"
	"github.com/testcontainers/testcontainers-go/modulegen/internal/mkdocs"
)

func TestExample(t *testing.T) {
	tests := []struct {
		name                  string
		example               Example
		expectedContainerName string
		expectedEntrypoint    string
		expectedTitle         string
	}{
		{
			name: "Module with title",
			example: Example{
				Name:      "mongoDB",
				IsModule:  true,
				Image:     "mongodb:latest",
				TitleName: "MongoDB",
			},
			expectedContainerName: "MongoDBContainer",
			expectedEntrypoint:    "RunContainer",
			expectedTitle:         "MongoDB",
		},
		{
			name: "Module without title",
			example: Example{
				Name:     "mongoDB",
				IsModule: true,
				Image:    "mongodb:latest",
			},
			expectedContainerName: "MongodbContainer",
			expectedEntrypoint:    "RunContainer",
			expectedTitle:         "Mongodb",
		},
		{
			name: "Example with title",
			example: Example{
				Name:      "mongoDB",
				IsModule:  false,
				Image:     "mongodb:latest",
				TitleName: "MongoDB",
			},
			expectedContainerName: "mongoDBContainer",
			expectedEntrypoint:    "runContainer",
			expectedTitle:         "MongoDB",
		},
		{
			name: "Example without title",
			example: Example{
				Name:     "mongoDB",
				IsModule: false,
				Image:    "mongodb:latest",
			},
			expectedContainerName: "mongodbContainer",
			expectedEntrypoint:    "runContainer",
			expectedTitle:         "Mongodb",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			example := test.example

			assert.Equal(t, "mongodb", example.Lower())
			assert.Equal(t, test.expectedTitle, example.Title())
			assert.Equal(t, test.expectedContainerName, example.ContainerName())
			assert.Equal(t, test.expectedEntrypoint, example.Entrypoint())
		})
	}
}

func TestExample_Validate(outer *testing.T) {
	outer.Parallel()

	tests := []struct {
		name        string
		example     Example
		expectedErr error
	}{
		{
			name: "only alphabetical characters in name/title",
			example: Example{
				Name:      "AmazingDB",
				TitleName: "AmazingDB",
			},
		},
		{
			name: "alphanumerical characters in name",
			example: Example{
				Name:      "AmazingDB4tw",
				TitleName: "AmazingDB",
			},
		},
		{
			name: "alphanumerical characters in title",
			example: Example{
				Name:      "AmazingDB",
				TitleName: "AmazingDB4tw",
			},
		},
		{
			name: "non-alphanumerical characters in name",
			example: Example{
				Name:      "Amazing DB 4 The Win",
				TitleName: "AmazingDB",
			},
			expectedErr: errors.New("invalid name: Amazing DB 4 The Win. Only alphanumerical characters are allowed (leading character must be a letter)"),
		},
		{
			name: "non-alphanumerical characters in title",
			example: Example{
				Name:      "AmazingDB",
				TitleName: "Amazing DB 4 The Win",
			},
			expectedErr: errors.New("invalid title: Amazing DB 4 The Win. Only alphanumerical characters are allowed (leading character must be a letter)"),
		},
		{
			name: "leading numerical character in name",
			example: Example{
				Name:      "1AmazingDB",
				TitleName: "AmazingDB",
			},
			expectedErr: errors.New("invalid name: 1AmazingDB. Only alphanumerical characters are allowed (leading character must be a letter)"),
		},
		{
			name: "leading numerical character in title",
			example: Example{
				Name:      "AmazingDB",
				TitleName: "1AmazingDB",
			},
			expectedErr: errors.New("invalid title: 1AmazingDB. Only alphanumerical characters are allowed (leading character must be a letter)"),
		},
	}

	for _, test := range tests {
		outer.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedErr, test.example.Validate())
		})
	}
}

func TestGenerateWrongExampleName(t *testing.T) {
	tmpCtx := NewContext(t.TempDir())
	examplesTmp := filepath.Join(tmpCtx.RootDir, "examples")
	examplesDocTmp := filepath.Join(tmpCtx.DocsDir(), "examples")
	githubWorkflowsTmp := tmpCtx.GithubWorkflowsDir()

	err := os.MkdirAll(examplesTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(examplesDocTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(githubWorkflowsTmp, 0o777)
	assert.Nil(t, err)

	err = copyInitialMkdocsConfig(t, tmpCtx)
	assert.Nil(t, err)

	tests := []struct {
		name string
	}{
		{name: " foo"},
		{name: "foo "},
		{name: "foo bar"},
		{name: "foo-bar"},
		{name: "foo/bar"},
		{name: "foo\\bar"},
		{name: "1foo"},
		{name: "foo1"},
		{name: "-foo"},
		{name: "foo-"},
	}

	for _, test := range tests {
		example := Example{
			Name:  test.name,
			Image: "docker.io/example/" + test.name + ":latest",
		}

		err = generate(example, tmpCtx)
		assert.Error(t, err)
	}
}

func TestGenerateWrongExampleTitle(t *testing.T) {
	tmpCtx := NewContext(t.TempDir())
	examplesTmp := filepath.Join(tmpCtx.RootDir, "examples")
	examplesDocTmp := filepath.Join(tmpCtx.DocsDir(), "examples")
	githubWorkflowsTmp := tmpCtx.GithubWorkflowsDir()

	err := os.MkdirAll(examplesTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(examplesDocTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(githubWorkflowsTmp, 0o777)
	assert.Nil(t, err)

	err = copyInitialMkdocsConfig(t, tmpCtx)
	assert.Nil(t, err)

	tests := []struct {
		title string
	}{
		{title: " fooDB"},
		{title: "fooDB "},
		{title: "foo barDB"},
		{title: "foo-barDB"},
		{title: "foo/barDB"},
		{title: "foo\\barDB"},
		{title: "1fooDB"},
		{title: "foo1DB"},
		{title: "-fooDB"},
		{title: "foo-DB"},
	}

	for _, test := range tests {
		example := Example{
			Name:      "foo",
			TitleName: test.title,
			Image:     "docker.io/example/foo:latest",
		}

		err = generate(example, tmpCtx)
		assert.Error(t, err)
	}
}

func TestGenerate(t *testing.T) {
	tmpCtx := NewContext(t.TempDir())
	examplesTmp := filepath.Join(tmpCtx.RootDir, "examples")
	examplesDocTmp := filepath.Join(tmpCtx.DocsDir(), "examples")
	githubWorkflowsTmp := tmpCtx.GithubWorkflowsDir()

	err := os.MkdirAll(examplesTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(examplesDocTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(githubWorkflowsTmp, 0o777)
	assert.Nil(t, err)

	err = copyInitialMkdocsConfig(t, tmpCtx)
	assert.Nil(t, err)

	originalConfig, err := mkdocs.ReadConfig(tmpCtx.MkdocsConfigFile())
	assert.Nil(t, err)

	err = copyInitialDependabotConfig(t, tmpCtx)
	assert.Nil(t, err)

	originalDependabotConfigUpdates, err := dependabot.GetUpdates(tmpCtx.DependabotConfigFile())
	assert.Nil(t, err)

	example := Example{
		Name:      "foodb4tw",
		TitleName: "FooDB4TheWin",
		IsModule:  false,
		Image:     "docker.io/example/foodb:latest",
	}
	exampleNameLower := example.Lower()

	err = generate(example, tmpCtx)
	assert.Nil(t, err)

	exampleDirPath := filepath.Join(examplesTmp, exampleNameLower)

	exampleDirFileInfo, err := os.Stat(exampleDirPath)
	assert.Nil(t, err) // error nil implies the file exist
	assert.True(t, exampleDirFileInfo.IsDir())

	exampleDocFile := filepath.Join(examplesDocTmp, exampleNameLower+".md")
	_, err = os.Stat(exampleDocFile)
	assert.Nil(t, err) // error nil implies the file exist

	mainWorkflowFile := filepath.Join(githubWorkflowsTmp, "ci.yml")
	_, err = os.Stat(mainWorkflowFile)
	assert.Nil(t, err) // error nil implies the file exist

	assertExampleDocContent(t, example, exampleDocFile)
	assertExampleGithubWorkflowContent(t, example, mainWorkflowFile)

	generatedTemplatesDir := filepath.Join(examplesTmp, exampleNameLower)
	assertExampleTestContent(t, example, filepath.Join(generatedTemplatesDir, exampleNameLower+"_test.go"))
	assertExampleContent(t, example, filepath.Join(generatedTemplatesDir, exampleNameLower+".go"))
	assertGoModContent(t, example, originalConfig.Extra.LatestVersion, filepath.Join(generatedTemplatesDir, "go.mod"))
	assertMakefileContent(t, example, filepath.Join(generatedTemplatesDir, "Makefile"))
	assertMkdocsExamplesNav(t, example, originalConfig, tmpCtx)
	assertDependabotExamplesUpdates(t, example, originalDependabotConfigUpdates, tmpCtx)
}

func TestGenerateModule(t *testing.T) {
	tmpCtx := NewContext(t.TempDir())
	modulesTmp := filepath.Join(tmpCtx.RootDir, "modules")
	modulesDocTmp := filepath.Join(tmpCtx.DocsDir(), "modules")
	githubWorkflowsTmp := tmpCtx.GithubWorkflowsDir()

	err := os.MkdirAll(modulesTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(modulesDocTmp, 0o777)
	assert.Nil(t, err)
	err = os.MkdirAll(githubWorkflowsTmp, 0o777)
	assert.Nil(t, err)

	err = copyInitialMkdocsConfig(t, tmpCtx)
	assert.Nil(t, err)

	originalConfig, err := mkdocs.ReadConfig(tmpCtx.MkdocsConfigFile())
	assert.Nil(t, err)

	err = copyInitialDependabotConfig(t, tmpCtx)
	assert.Nil(t, err)

	originalDependabotConfigUpdates, err := dependabot.GetUpdates(tmpCtx.DependabotConfigFile())
	assert.Nil(t, err)

	example := Example{
		Name:      "foodb",
		TitleName: "FooDB",
		IsModule:  true,
		Image:     "docker.io/example/foodb:latest",
	}
	exampleNameLower := example.Lower()

	err = generate(example, tmpCtx)
	assert.Nil(t, err)

	exampleDirPath := filepath.Join(modulesTmp, exampleNameLower)

	exampleDirFileInfo, err := os.Stat(exampleDirPath)
	assert.Nil(t, err) // error nil implies the file exist
	assert.True(t, exampleDirFileInfo.IsDir())

	exampleDocFile := filepath.Join(modulesDocTmp, exampleNameLower+".md")
	_, err = os.Stat(exampleDocFile)
	assert.Nil(t, err) // error nil implies the file exist

	mainWorkflowFile := filepath.Join(githubWorkflowsTmp, "ci.yml")
	_, err = os.Stat(mainWorkflowFile)
	assert.Nil(t, err) // error nil implies the file exist

	assertExampleDocContent(t, example, exampleDocFile)
	assertExampleGithubWorkflowContent(t, example, mainWorkflowFile)

	generatedTemplatesDir := filepath.Join(modulesTmp, exampleNameLower)
	assertExampleTestContent(t, example, filepath.Join(generatedTemplatesDir, exampleNameLower+"_test.go"))
	assertExampleContent(t, example, filepath.Join(generatedTemplatesDir, exampleNameLower+".go"))
	assertGoModContent(t, example, originalConfig.Extra.LatestVersion, filepath.Join(generatedTemplatesDir, "go.mod"))
	assertMakefileContent(t, example, filepath.Join(generatedTemplatesDir, "Makefile"))
	assertMkdocsExamplesNav(t, example, originalConfig, tmpCtx)
	assertDependabotExamplesUpdates(t, example, originalDependabotConfigUpdates, tmpCtx)
}

// assert content in the Examples nav from mkdocs.yml
func assertDependabotExamplesUpdates(t *testing.T, example Example, originalConfigUpdates dependabot.Updates, tmpCtx *Context) {
	examples, err := dependabot.GetUpdates(tmpCtx.DependabotConfigFile())
	assert.Nil(t, err)

	assert.Equal(t, len(originalConfigUpdates)+1, len(examples))

	// the example should be in the dependabot updates
	found := false
	for _, ex := range examples {
		directory := "/" + example.ParentDir() + "/" + example.Lower()
		if directory == ex.Directory {
			found = true
		}
	}

	assert.True(t, found)

	// first item is the github-actions module
	assert.Equal(t, "/", examples[0].Directory, examples)
	assert.Equal(t, "github-actions", examples[0].PackageEcosystem, "PackageEcosystem should be github-actions")

	// second item is the core module
	assert.Equal(t, "/", examples[1].Directory, examples)
	assert.Equal(t, "gomod", examples[1].PackageEcosystem, "PackageEcosystem should be gomod")

	// third item is the pip module
	assert.Equal(t, "/", examples[2].Directory, examples)
	assert.Equal(t, "pip", examples[2].PackageEcosystem, "PackageEcosystem should be pip")
}

// assert content example file in the docs
func assertExampleDocContent(t *testing.T, example Example, exampleDocFile string) {
	content, err := os.ReadFile(exampleDocFile)
	assert.Nil(t, err)

	lower := example.Lower()
	title := example.Title()

	data := sanitiseContent(content)
	assert.Equal(t, data[0], "# "+title)
	assert.Equal(t, data[2], `Not available until the next release of testcontainers-go <a href="https://github.com/testcontainers/testcontainers-go"><span class="tc-version">:material-tag: main</span></a>`)
	assert.Equal(t, data[4], "## Introduction")
	assert.Equal(t, data[6], "The Testcontainers module for "+title+".")
	assert.Equal(t, data[8], "## Adding this module to your project dependencies")
	assert.Equal(t, data[10], "Please run the following command to add the "+title+" module to your Go dependencies:")
	assert.Equal(t, data[13], "go get github.com/testcontainers/testcontainers-go/"+example.ParentDir()+"/"+lower)
	assert.Equal(t, data[18], "<!--codeinclude-->")
	assert.Equal(t, data[19], "[Creating a "+title+" container](../../"+example.ParentDir()+"/"+lower+"/"+lower+".go)")
	assert.Equal(t, data[20], "<!--/codeinclude-->")
	assert.Equal(t, data[22], "<!--codeinclude-->")
	assert.Equal(t, data[23], "[Test for a "+title+" container](../../"+example.ParentDir()+"/"+lower+"/"+lower+"_test.go)")
	assert.Equal(t, data[24], "<!--/codeinclude-->")
	assert.Equal(t, data[28], "The "+title+" module exposes one entrypoint function to create the "+title+" container, and this function receives two parameters:")
	assert.True(t, strings.HasSuffix(data[31], "(*"+title+"Container, error)"))
	assert.Equal(t, "for "+title+". E.g. `testcontainers.WithImage(\""+example.Image+"\")`.", data[44])
}

// assert content example test
func assertExampleTestContent(t *testing.T, example Example, exampleTestFile string) {
	content, err := os.ReadFile(exampleTestFile)
	assert.Nil(t, err)

	data := sanitiseContent(content)
	assert.Equal(t, data[0], "package "+example.Lower())
	assert.Equal(t, data[7], "func Test"+example.Title()+"(t *testing.T) {")
	assert.Equal(t, data[10], "\tcontainer, err := "+example.Entrypoint()+"(ctx)")
}

// assert content example
func assertExampleContent(t *testing.T, example Example, exampleFile string) {
	content, err := os.ReadFile(exampleFile)
	assert.Nil(t, err)

	lower := example.Lower()
	containerName := example.ContainerName()
	exampleName := example.Title()
	entrypoint := example.Entrypoint()

	data := sanitiseContent(content)
	assert.Equal(t, data[0], "package "+lower)
	assert.Equal(t, data[8], "// "+containerName+" represents the "+exampleName+" container type used in the module")
	assert.Equal(t, data[9], "type "+containerName+" struct {")
	assert.Equal(t, data[13], "// "+entrypoint+" creates an instance of the "+exampleName+" container type")
	assert.Equal(t, data[14], "func "+entrypoint+"(ctx context.Context, opts ...testcontainers.ContainerCustomizer) (*"+containerName+", error) {")
	assert.Equal(t, data[16], "\t\tImage: \""+example.Image+"\",")
	assert.Equal(t, data[33], "\treturn &"+containerName+"{Container: container}, nil")
}

// assert content GitHub workflow for the example
func assertExampleGithubWorkflowContent(t *testing.T, example Example, exampleWorkflowFile string) {
	content, err := os.ReadFile(exampleWorkflowFile)
	assert.Nil(t, err)

	data := sanitiseContent(content)
	ctx := getTestRootContext(t)

	modulesList, err := ctx.GetModules()
	assert.Nil(t, err)
	assert.Equal(t, "        module: ["+strings.Join(modulesList, ", ")+"]", data[94])

	examplesList, err := ctx.GetExamples()
	assert.Nil(t, err)
	assert.Equal(t, "        module: ["+strings.Join(examplesList, ", ")+"]", data[110])
}

// assert content go.mod
func assertGoModContent(t *testing.T, example Example, tcVersion string, goModFile string) {
	content, err := os.ReadFile(goModFile)
	assert.Nil(t, err)

	data := sanitiseContent(content)
	assert.Equal(t, "module github.com/testcontainers/testcontainers-go/"+example.ParentDir()+"/"+example.Lower(), data[0])
	assert.Equal(t, "require github.com/testcontainers/testcontainers-go "+tcVersion, data[4])
	assert.Equal(t, "replace github.com/testcontainers/testcontainers-go => ../..", data[6])
}

// assert content Makefile
func assertMakefileContent(t *testing.T, example Example, makefile string) {
	content, err := os.ReadFile(makefile)
	assert.Nil(t, err)

	data := sanitiseContent(content)
	assert.Equal(t, data[4], "\t$(MAKE) test-"+example.Lower())
}

// assert content in the Examples nav from mkdocs.yml
func assertMkdocsExamplesNav(t *testing.T, example Example, originalConfig *mkdocs.Config, tmpCtx *Context) {
	config, err := mkdocs.ReadConfig(tmpCtx.MkdocsConfigFile())
	assert.Nil(t, err)

	parentDir := example.ParentDir()

	examples := config.Nav[4].Examples
	expectedEntries := originalConfig.Nav[4].Examples
	if example.IsModule {
		examples = config.Nav[3].Modules
		expectedEntries = originalConfig.Nav[3].Modules
	}

	assert.Equal(t, len(expectedEntries)+1, len(examples))

	// the example should be in the nav
	found := false
	for _, ex := range examples {
		markdownExample := example.ParentDir() + "/" + example.Lower() + ".md"
		if markdownExample == ex {
			found = true
		}
	}

	assert.True(t, found)

	// first item is the index
	assert.Equal(t, parentDir+"/index.md", examples[0], examples)
}

func sanitiseContent(bytes []byte) []string {
	content := string(bytes)

	// Windows uses \r\n for newlines, but we want to use \n
	content = strings.ReplaceAll(content, "\r\n", "\n")

	data := strings.Split(content, "\n")

	return data
}

func copyInitialDependabotConfig(t *testing.T, tmpCtx *Context) error {
	ctx := getTestRootContext(t)
	return dependabot.CopyConfig(ctx.DependabotConfigFile(), tmpCtx.DependabotConfigFile())
}
