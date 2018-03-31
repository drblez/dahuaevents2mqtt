# Send Dahua IP Camera events to MQTT

# Configuration file

    {
      "Defaults": {
        "Map": {
          "Start": "ON",
          "Stop": "OFF"
        }
      },
      "Cameras": [
        {
          "Host": "192.168.1.248",
          "Port": "80",
          "Username": "admin",
          "Password": "admin",
          "Events": ["VideoMotion"],
          "Topic": "/openhab/parking_1_motion/state/set",
          "Map": {
            "Start": "ON",
            "Stop": "OFF"
          }
        },
        {
          "Host": "192.168.1.194",
          "Topic": "/openhab/parking_2_motion/state/set"
        }
      ],
      "MQTT": {
        "Host": "192.168.1.31",
        "Port": "1883",
        "Timeout": 60
      }
    }

## Camera parameters

- Host (string, no default value) camera host, ex. *"192.168.1.248"*;
- Port (string, default value *"80"*) camera port;
- Username (string, default value *"admin"*) camera user;
- Password (string, default value *"admin"*) camera password;
- Events (array of string, default value *\["VideoMotion"\]*) list of camera's event for handle by service;
- Topic (string, no default value) MQTT topic for sending messages, ex. *"/openhab/parking_1_motion/state/set"*;
- Map (map with string values, no default value)

### Map

The events receiving from camera looks like this:

    Code=VideoMotion;action=Stop;index=0

Before sending an event to a queue, you may need to convert the *action* to the appropriate value for you.

    "Map": {
        "Start": "ON",
        "Stop": "OFF"
    }

# Parameters

    dahuaevents2mqtt [install|remove|run|start|stop]

## Install

Installing as service

## Remove

Remove service from system

## Run

Run as regular program

## Start

Start previously installed service

## Stop

Stop previously installed service

# Plans

- Map for Code/Topic