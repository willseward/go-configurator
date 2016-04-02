#!/bin/bash

# This should be the entrypoint for this container
# Tasks are as follows:
#   1/ Run the configurator and setup the audit configuration
#   2/ Start the audit daemon
#   3/ Go back to step 1 every x minutes (default x=10)

dockeraudit-host update --temp=/tmp/dockeraudit --slave