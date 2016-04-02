package main

import (
    "os"
    
    "github.com/codegangsta/cli"
)

func main() {    
    app := cli.NewApp()
    
    app.Name = "Docker Auditd Exporter - Host Agent"
    app.Usage = "Exports the audit reports from containers to an API"
    app.Version = "0.1"
    
    SetupCliForConfigurator(app)
    // SetupCliForServer(app)
    
    app.Run(os.Args)
}