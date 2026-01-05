#!/bin/bash

export SERIAL=USD3A8000375
echo $BS_NETWORK
echo $BS_USERNAME
export PACKAGENAME=my-simple-setup

# Create setup configuration file from embedded JSON
cat > /tmp/simple-setup.json <<EOF
{
    "networkName": "$BS_NETWORK",
    "username": "$BS_USERNAME",
    "packageName": "$PACKAGENAME",
    "setupType": "lfn",
    "network": {
        "timeServers": [
            "0.pool.ntp.org"
        ],
        "interfaces": [
            {
                "id": "eth0",
                "name": "Ethernet",
                "type": "ethernet",
                "proto": "DHCPv4",
                "contentDownloadEnabled": true,
                "healthReportingEnabled": true
            }
        ]
    }
}
EOF

# add the setup package
./bin/bdeploy-add-setup /tmp/simple-setup.json

# associate the player with the setup package
./bin/bdeploy-associate --setup-name "$PACKAGENAME" --serial "$SERIAL"

./bin/rdws-reformat-storage --serial "$SERIAL" --device sd -y
