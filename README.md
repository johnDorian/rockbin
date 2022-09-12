# !!!Archived!!!

**I have purchased a new vacuum and no longer use this project, and have therefore archived this project. 
I will not make any new changes to the code, and due to potential future security issues, 
I would recommend that users find a different alternative. 
Feel free to fork the repo if you want.**


# rockbin



[![building](https://github.com/johnDorian/rockbin/actions/workflows/ci.yml/badge.svg)]((https://github.com/johnDorian/rockbin/actions/workflows/ci.yml/badge.svg))
[![codecov](https://codecov.io/gh/johnDorian/rockbin/branch/master/graph/badge.svg)](https://codecov.io/gh/johnDorian/rockbin)
[![gosec](https://goreportcard.com/badge/github.com/johnDorian/rockbin)]((https://goreportcard.com/badge/github.com/johnDorian/rockbin))
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FjohnDorian%2Frockbin.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FjohnDorian%2Frockbin?ref=badge_shield)

This repo contains the code for a simple go based mqtt client which will send the bin status to a mqtt server. 

I'm using home assistant, hence the home assistant auto discovery stuff.

The general idea is that the home assistant config is sent every minute, and the bin value is sent when the `/mnt/data/rockrobo/RoboController.cfg` is modified. The file is only modified after cleaning or returning to the dock. Bin value in RoboController.cfg is named: bin_in_time. This value is expressed in seconds from last cleaning of vacuum bin. After cleaning value is set to 0.

I decided to use a percentage value, in order to use a gauge in home assistant. I'm using 40 minutes as the time that the vacuum needs to be emptied. This can be changed as an input value to the mqtt client. 
It's also possible to create sensor expressed in minutes - please refer to parameters table.


## Contributing

Please contact me first with any suggestied changes, or before opening a PR. The main aim and scope of this project is to get the robot to report the bin status, so please keep this in mind when creating a PR. 


## Installation and setup.

The client can be started/tested using the following commands.

```bash
# copy the binary to the vacuum
scp rockbin root@...:/root/rockbin
# ssh into the vacuum
ssh root@...
# move the binary into the correct location
mv rockbin /usr/local/bin/rockbin
# make the binary executable
chmod +x rockbin
# setup the config file and install the required service script (e.g. /etc/init/S12rockbin). 
# This will overwrite any existing service scripts - please make a backup beforehand. 
rockbin configure
# test the new version
rockbin serve --log_level debug
# If everything seems to be working finish restart the vacuum
reboot now
```

The above command set will setup the rockbin service and setup the configuration file according to your responses.

## Home assistant 
An example of sending the vacuum to the rubbish bin is below: 

```yaml
  - alias: 'Send vacuum to the bin'
    trigger: 
      platform: state
      entity_id: vacuum.rockrobo
      
      to: "returning"
      from: "cleaning"
      
      for: "00:00:03"
    condition:
      - condition: numeric_state
        entity_id: sensor.vacuumbin
        above: 100
    action: 
    - service: mqtt.publish
      data: 
        topic: valetudo/rockrobo/command
        payload: "pause"
    - wait_template: '{{ is_state(''vacuum.rockrobo'', ''idle'') }}'
      continue_on_timeout: 'true'
      timeout: 00:00:05
    - service: mqtt.publish
      data: 
        topic: valetudo/rockrobo/custom_command
        payload: "{\"command\":\"go_to\", \"spot_id\":\"bin\"}"
    - service: notify.telegram_user
      data:
        message: "Please empty the vacuum"
        title: "Vacuum going to the bin"

  - alias: 'Go home when bin is empty'
    trigger: 
      platform: numeric_state
      entity_id: sensor.vacuumbin
      below: 1
    action: 
      - service: vacuum.return_to_base
        entity_id: vacuum.xiaomi_vacuum_cleaner
```


## Building for the vacuum

***Pre-built binaries are available in the releases.***

You can build it on your computer rather than the vacuum using: 

```bash 
GOARM=7 GOARCH=arm GOOS=linux go build 
# or use the taskfile:
task build_vacuum
```
