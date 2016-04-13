package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "path/filepath"
    
    "gopkg.in/yaml.v2"
)

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

func ensureDirectoryCreated(directoryPath string) error {
    _, err := os.Stat(directoryPath)
    
    if os.IsNotExist(err) {
        if err = os.MkdirAll(directoryPath, 0755); err != nil {                    
            return fmt.Errorf("Could not create directory for configuration files:", err)
        }
    } else if err != nil {
        return fmt.Errorf("Other file error:", err)
    }
    
    return nil
}

func ensureFileCreated(filePath string) (file *os.File, err error) {
    _, err = os.Stat(filePath)
    
    if os.IsNotExist(err) {
        file, err = os.Create(filePath)
        if err != nil {                 
            return nil, fmt.Errorf("Could not create file for configuration files")
        }
    } else if err != nil {
        return nil, err
    }
    
    file, err = os.Create(filePath) 
    if err != nil {
        return nil, fmt.Errorf("Could not open file", err)
    }
    
    return 
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

// createFileIfNotExistsAndIsRegularOrError checks and makes sure a file exists and is regular. If not,
// it will be created. This function also checks the dir structure.
func createFileIfNotExistsAndIsRegularOrError(filePath string) error {
    dir, _ := filepath.Split(filePath)

    if err := ensureDirectoryCreated(dir); err != nil {
        return err
    }
    
    if _, err := ensureFileCreated(filePath); err != nil {
        return err
    }
    
    // Try one more time 
    stat, _ := os.Stat(filePath)

    if !stat.Mode().IsRegular() {
        return fmt.Errorf("file is not regular: %s", filePath)
    }
    
    return nil
}