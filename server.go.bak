package main

import (
    "os/exec"
    "log"
    "fmt"
    "net/http"
    "strings"
    "encoding/json"
    "strconv"
    
    "github.com/gorilla/mux"
    "github.com/codegangsta/cli"
)

type Server struct {
    port int
    serverName string
}

type FileRecord struct {
    Permissions string  `json:"permissions"`
    Inodes int          `json:"inodes"`
    Owner string        `json:"owner"`
    Group string        `json:"group"`
    Size int            `json:"size"`
    Month string        `json:"month"`
    Date int            `json:"date"`
    Time string         `json:"time"`
    Name string         `json:"name"`
}

func NewServer(port int, serverName string) (server *Server) {
    server = &Server{
        port: port,
        serverName: serverName}
    
    return
}

func SetupCliForServer(app *cli.App) {    
    
    commands := []cli.Command {
        {
            Name:    "server",
            Aliases: []string{"s"},
            Usage:   "starts an api server to relay audit reports",
            Action: func(c *cli.Context) {
                port, _ := strconv.Atoi(c.String("port"))
                address := c.String("server-address")
                
                server := NewServer(port, address)
                server.run()
            },
            Flags: []cli.Flag {
                cli.StringFlag{
                    Name: "port",
                    Value: "80",
                    Usage: "bind-to port",
                },
                cli.StringFlag{
                    Name: "server-address",
                    Value: "0.0.0.0",
                    Usage: "bind-to address",
                },
            },
        },
    }
    
    app.Commands = append(app.Commands, commands...)
}

func (s *Server) run() {
    
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/ls", Index)
    
    log.Println("Starting server on port", s.port)
    portString := fmt.Sprintf(":%d", s.port)
    log.Fatal(http.ListenAndServe(portString, router))
}

func delimiter(r rune) (bool){
    return r == ' '
}

func Index(w http.ResponseWriter, r *http.Request) {
    
    get := r.URL.Query()
    paths, ok := get["path"]
    var path string
    if ok {
        if len(paths) > 0 {
            path = paths[0]
        }
    }
    
    command := exec.Command("ls", "-l", path)
    if bytes, err := command.Output(); err != nil {
        fmt.Fprintln(w, err)
    } else {
        output := string(bytes)
        
        lines := strings.Split(output, "\n") 
        records := make([]FileRecord, len(lines))
                
        for idx,line := range lines {
            elements := strings.FieldsFunc(line, delimiter)
            
            if len(elements) >= 9 {
                // Join filename slice
                elements[8] = strings.Join(elements[8:], " ")
                inodes, _ := strconv.Atoi(elements[1])
                size, _ := strconv.Atoi(elements[4])
                date, _ := strconv.Atoi(elements[6])
            
                file := FileRecord{
                    Permissions: elements[0],
                    Inodes: inodes,
                    Owner: elements[2],
                    Group: elements[3],
                    Size: size,
                    Month: elements[5],
                    Date: date,
                    Time: elements[7],
                    Name: elements[8]}
                records[idx] = file
            }
        }
        
        log.Println("[Server] serving request for", path)
        json.NewEncoder(w).Encode(records)
    }
}