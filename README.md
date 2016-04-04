# go-configurator

A simple tool to aid in the battle of configuration file synchronization

## Usage

```
NAME:
   Docker Auditd Exporter - Host Agent - Exports the audit reports from containers to an API

USAGE:
   go-configurator [global options] command [command options] [arguments...]
   
VERSION:
   0.1
   
COMMANDS:
    update	updates the auditd configuration

GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
```

## Building

```
$ go get ./...
$ go install .
$ ./go-configurator update --templates=samples/templates --master --temp=samples/dist --config=samples/config.yml --test 
```

## Building from a Dockerfile

```
$ docker build -t configurator .
$ cd sample
$ docker build -t config-sample .
$ docker run config-sample update --templates=/var/go-configurator/templates --temp=dist/ --test
```

## Future 

- [X] YAML files for variables in templates
- [ ] Intervalic updates
- [ ] Before and After scripts
