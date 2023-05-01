#!/usr/bin/env bash

# Global vars
myCmd=""

# ----------------------------------------------------------------------------------------------------------------
# Functions
# ----------------------------------------------------------------------------------------------------------------

function base_command() {
    # Check auth key was passed
    if [ -z "$KEY" ]
    then
        echo "Environment Variable \$KEY is empty, exitting!"
        exit 1
    else
        myCmd="tailscale up --authkey=$KEY"
    fi
}

function add_exit_node() {
    # Check auth key was passed
    if [ -z "$EXIT_NODE" ]
    then
        echo "Exit Node was not requested and will not be configured."
        EXIT_NODE=$false
    fi

    # If exit node if requested, add to command
    if [ "$EXIT_NODE" = true ]
    then
        echo "Configuring as an exit node"
        myCmd+=" --advertise-exit-node"
    fi
}

function add_avertise_routes() {
    if [ -z "$ADVERTISE_ROUTES" ]; 
    then 
        echo "No routes were requested to be advertised"; 
    else 
        echo "Will addvertise the provided routes in env var ADVERTISE_ROUTES: $ADVERTISE_ROUTES"
        myCmd+=" --advertise-routes=$ADVERTISE_ROUTES"
    fi
}

function add_machine_name() {
    if [ -z "$MACHINE_NAME" ]; 
    then 
        echo "No routes were requested to be advertised"; 
    else 
        echo "Will use MACHINE_NAME: $MACHINE_NAME"
        myCmd+=" --hostname=$MACHINE_NAME"
    fi
}

# ----------------------------------------------------------------------------------------------------------------
# Main
# ----------------------------------------------------------------------------------------------------------------

# Create base command
base_command

# Copgiure Exit Node
add_exit_node

# Configure advertise_routes
add_avertise_routes

# Configure tailscale machine name
add_machine_name

# Start tailscale daemon with user mode networking
tailscaled --tun=userspace-networking --socks5-server=localhost:1055 &

# Start aproxy
/aproxy &

# Output and Run tailscale command
echo $myCmd
eval $myCmd

# Run Inifinite loop to keep container alive
while true; 
do 
    echo "$(date), Sleeping 30 seconds..."
    sleep 30
done

