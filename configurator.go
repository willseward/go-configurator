package main

import (
    "os"
    "io"
    "io/ioutil"
    "log"
    "fmt"
    "strings"
    "text/template"
    "path/filepath"
        
    "github.com/codegangsta/cli"
    "gopkg.in/yaml.v2"
)

// ConfigurationFile holds the record for a config file that should be parsed into 
// the fs
type ConfigurationFile struct {
    templateFilePath string
    configDestinationDir string
    configFileName string
}

// Configurator is the object that does the configuring
type Configurator struct {
    buildDir string
    configFiles []ConfigurationFile
    config map[interface{}]interface{}
}

// Config is a type for the data in the templates
type Config struct {
    Master bool
    TcpPort int
    RemoteAddress string
}

// findAndCreateConfigurationFileRecords looks at the template directory and creates the ConfigurationFile
// objects
func findAndCreateConfigurationFileRecords(templateDir string) (files []ConfigurationFile, err error) {
    // Will create the ConfigurationFile array 
    
    files = make([]ConfigurationFile, 0)
    
    templateDirComponents := strings.Split(templateDir, "/")
    
    // Walk through the template file dir
    filepath.Walk(templateDir, func(path string, f os.FileInfo, err error) error {
        if !f.Mode().IsDir() {
            if strings.HasSuffix(path, ".tmpl") {
                
                dir, file := filepath.Split(path)
                file = strings.TrimRight(file, ".tmpl")
                
                // Remove templateDir from front of dir
                dirComponents := strings.Split(filepath.Clean(dir), "/")
                dirComponents = dirComponents[len(templateDirComponents):]
                dirComponents = append([]string{"/"}, dirComponents...)
                dir = filepath.Join(dirComponents...)
                
                configFile := ConfigurationFile{
                    templateFilePath: path,
                    configDestinationDir: dir,
                    configFileName: file,
                }
                                
                files = append(files, configFile)
            }
        }
        
        return nil
    })
    
    return files, nil
}

// NewConfigurator creates and returns a new Configurator
func NewConfigurator(buildDir string, configFiles []ConfigurationFile, config map[interface{}]interface{}) (*Configurator) {
    return &Configurator{
        buildDir: buildDir,
        configFiles: configFiles,
        config: config,
    }
}

// getTextForFile returns the text contents of a file
func getTextForFile(fileName string) (data string, err error) {
    if bytes, err := ioutil.ReadFile(fileName); err != nil {
        return "", err
    } else {
        data = string(bytes)
        return data, nil
    }
}

// rebuildConfig parses and saves a template record to the temp dir
func (c *Configurator) rebuildConfig(file ConfigurationFile) (path string, err error) {
    // Get file text
    configText, _ := getTextForFile(file.templateFilePath)
    
    // Build templates
    configTemplate, err := template.New(file.configFileName).Parse(configText) 
    if err != nil {
        log.Println("[Configurator] Error pasring template:", file.configFileName)
        return "", err
    }

    // Create the file
    dir := filepath.Join(c.buildDir, file.configDestinationDir)
    if err = os.MkdirAll(dir, 0755); err != nil {
        log.Println("[Configurator] temp directory not writable")
        return "", err
    }
    
    path = filepath.Join(dir, file.configFileName)
    f, err := os.Create(path)
    if err != nil {
        // Problem, move on the next file and report the error
        log.Println("[Configurator] Could not create config file destination:", file.configFileName)
        return "", err
    }
    
    // Closes the file when the function exits
    defer f.Close()
    
    // Process the templates
    err = configTemplate.Execute(f, c.config)
    if err != nil {
        log.Println("[Configurator] Error while executing template:", file.configFileName)
        return "", err
    }
        
    log.Println("[Configurator] Built config:", file.configFileName, "to", path) 
    
    return path, err
}

// createFileIfNotExistsAndIsRegularOrError checks and makes sure a file exists and is regular. If not,
// it will be created. This function also checks the dir structure.
func createFileIfNotExistsAndIsRegularOrError(filePath string) error {
    dir, _ := filepath.Split(filePath)
    stat, err := os.Stat(filePath)

    if os.IsNotExist(err) { 
        
        if err = os.MkdirAll(dir, 0755); err != nil {                    
            log.Println("[Configurator] Could not create directory for configuration files")
            return err
        }
        
        if _, err = os.Create(filePath); err != nil {
            log.Println("[Configurator] Could not create file for configuration files")
            return err
        }
    } 
    
    // Try one more time 
    stat, err = os.Stat(filePath)

    if !stat.Mode().IsRegular(){
        return fmt.Errorf("file is not regular: %s", filePath)
    }
    
    return nil
}

// replaceConfig updates the active config files on the system.
func (c *Configurator) replaceConfig(file ConfigurationFile, newConfigFilePath string) error {
    
    // Check both files exist and are regular
    if err := createFileIfNotExistsAndIsRegularOrError(newConfigFilePath); err != nil {
        log.Println("[Configurator] Source file error")
        return err
    }
        
    destFilePath := filepath.Join(file.configDestinationDir, file.configFileName)
    if err := createFileIfNotExistsAndIsRegularOrError(destFilePath); err != nil {
        log.Println("[Configurator] Destination file error")
        return err
    }
    
    src, err := os.Open(newConfigFilePath)
    if err != nil {
        log.Println("[Configurator] Error opening source file to replacement")
        return err
    }
    
    // Closes the file when the function exits
    defer src.Close()
    
    dest, err := os.OpenFile(destFilePath, os.O_WRONLY, os.ModeAppend)
    if err != nil {
        log.Println("[Configurator] Error opening destination file to replacement")
        return err
    }
    
    defer dest.Close()
    
    if _, err := io.Copy(dest, src); err != nil {
        log.Println("[Configurator] Error copying data")
        return err
    }
    
    // Ensures the contents end up on disk 
    if err := dest.Sync(); err != nil {
        log.Println("[Configurator] Error writing file")
        return err
    }

    log.Println("[Configurator] Successfully replace config file:", destFilePath)
    
    return nil
}

func yamlToMap(filename string) (config map[interface{}]interface{}, err error) {
    fileContents, err := getTextForFile(filename) 
    if err != nil {
        return nil, err
    }
    
    // The config object
    m := make(map[interface{}]interface{})
    
    err = yaml.Unmarshal([]byte(fileContents), &m)
    if err != nil {
        return nil, err
    }
    
    return m, nil
}

// SetupCliForConfigurator sets up the cli to respond to hooks for the Configurator
func SetupCliForConfigurator(app *cli.App) {    
    
    commands := []cli.Command {
        {
            Name:    "update",
            Aliases: []string{"u"},
            Usage:   "updates the auditd configuration",
            Action: func(c *cli.Context) {
                
                dist := c.String("temp")
                templates := c.String("templates")
                test := c.Bool("test")
                configYml := c.String("config")
                // beforeFile := c.String("before")
                // afterFile := c.String("after")
                
                files, err := findAndCreateConfigurationFileRecords(templates)
                if err != nil {
                    log.Println(err)
                    log.Fatal("[Configurator] Problem getting templates")
                }
                
                config, err := yamlToMap(configYml)
                if err != nil {
                    log.Println(err)
                    log.Fatal("[Configurator] Problem getting config from yml file")
                }
                                
                configurator := NewConfigurator(dist, files, config)
                
                if err := configurator.BuildAndUpdateConfig(test); err != nil {
                    log.Println(err)
                    log.Fatal("[Configurator] Problem updating config")
                } else {
                    log.Println("[Configurator] Success! Built", len(files), "config files")
                }
                
            },
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name: "templates",
                    Value: "/var/go-configurator/templates",
                    Usage: "the template file for auditd rules",
                },
                cli.StringFlag{
                    Name: "temp",
                    Value: "/tmp/go-configurator/dist",
                    Usage: "temp directory for resolved configuration files",
                },
                cli.StringFlag{
                    Name: "config",
                    Value: "/var/go-configurator/config.yml",
                    Usage: "the YAML config file to be used for configuration",  
                },
                cli.BoolFlag{
                    Name: "test",
                    Usage: "do not export config files to final location, just build them",
                },
                cli.StringFlag{
                    Name: "before",
                    Value: "",
                    Usage: "the script to run before the configurator starts",
                },
                cli.StringFlag{
                    Name: "after",
                    Value: "",
                    Usage: "the script to run after the configurator finishes",
                },
            },
        },
    }
    
    app.Commands = append(app.Commands, commands...)
}
 
// BuildAndUpdateConfig runs the configurator
func (c *Configurator) BuildAndUpdateConfig(test bool) (err error) {
    // Check that dir exists
    if _, err := os.Stat(c.buildDir); os.IsNotExist(err) {
        // It doesn't, so create it now
        if err = os.MkdirAll(c.buildDir, 0755); err != nil {                    
            log.Println("[Configurator] Could not create directory for temp files")
            return err
        }
    }
    
    for _, file := range c.configFiles {
        path, err := c.rebuildConfig(file);
        if err != nil {
            return err
        }  

        if !test {
            err = c.replaceConfig(file, path); 
            if err != nil {
                return err
            }
        }
    }
    
    return nil
}