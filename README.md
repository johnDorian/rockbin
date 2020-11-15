# rockbin

[![johnDorian](https://circleci.com/gh/johnDorian/rockbin.svg?style=shield)](https://circleci.com/gh/johnDorian/rockbin) [![codecov](https://codecov.io/gh/johnDorian/rockbin/branch/master/graph/badge.svg)](https://codecov.io/gh/johnDorian/rockbin)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FjohnDorian%2Frockbin.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FjohnDorian%2Frockbin?ref=badge_shield)

This repo contains the code for a simple go based mqtt client which will send the bin status to a mqtt server. 

I'm using home assistant, hence the home assistant auto discovery stuff.

The general idea is that the home assistant config is sent every minute, and the bin value is sent when the `/mnt/data/rockrobo/RoboController.cfg` is modified. The file is only modified after cleaning or returning to the dock. Bin value in RoboController.cfg is named: bin_in_time. This value is expressed in seconds from last cleaning of vacuum bin. After cleaning value is set to 0.

I decided to use a percentage value, in order to use a gauge in home assistant. I'm using 40 minutes as the time that the vacuum needs to be emptied. This can be changed as an input value to the mqtt client. 
It's also possible to create sensor expressed in minutes - please refer to parameters table.


## Contributing

Feel free to make some changes (even write some unit tests) and create a PR. The main aim and scope of this project is to get the robot to report the bin status, so please keep this in mind when creating a PR. 

## Building for the vacuum

***Pre-built binaries are available in the releases.***

You can build it on your computer rather than the vacuum using: 

```bash 
GOARM=7 GOARCH=arm GOOS=linux go build 
```

## Config

You will need to change the address of the mqtt broker and the amount of seconds which is considered to indicate the full capacity of the bin in the rockbin.conf file.

## Install

The client can be started/tested using the following command. 

If you'd like to have sensor expressed in % where 40mins (40*60 = 2400) is 100% then use this configuration:
```bash
./rockbin -mqtt_server mqtt://192.168.0.144:1883 -full_time 2400
```
If you'd like to have sensor expressed in minutes then use this configuration:
```bash
./rockbin -mqtt_server mqtt://192.168.0.144:1883 -measurement_unit min
```
You can pass usernames and passwords via the environment variables: 
```bash
MQTT_USERNAME=mrmqtt MQTT_PASSWORD=coolpass ./rockbin -mqtt_server mqtt://192.168.0.144:1883 -measurement_unit min
```
or via config:
```bash
./rockbin -mqtt_server mqtt://192.168.0.144:1883 -mqtt_user mrmqtt -mqtt_password coolpass -measurement_unit min
```
The state topic could be changed by the mqtt_state_topic config (defaults to 'homeassistant/sensor/%v/state'. The %v is replaced with the sensor_name value):
```bash
./rockbin -mqtt_server mqtt://192.168.0.144:1883 -mqtt_user mrmqtt -mqtt_password coolpass -mqtt_state_topic 'rockbin/%v/state' -measurement_unit min
```
You can increase the logging level to help setting up connections using 
```bash
./rockbin -mqtt_server mqtt://192.168.0.144:1883 -measurement_unit min -log_level debug
```



|parameter|default|description|
|---------|:-----:|:----------|
|-mqtt_server     |mqtt://localhost:1883|MQTT server address|
|-sensor_name|vacuumbin|Name of sensor created in Home Assistant.|
|-full_time|0|When 0 then sensor is expressed in minutes. When greater than 0 then sensor is expressed in % where full_time is number of seconds in 100%.| 
|-measurement_unit|%|In what unit should the measurement be sent (%, sec, min)|
|-file_path|/mnt/data/rockrobo/RoboController.cfg|file path of RoboController.cfg|
|-log_level|Fatal|Level of logging (trace, debug, info, warn, error, fatal, panic).|

        
If your mqtt broker requires authentication, you can set the environment variables (MQTT_USERNAME and MQTT_PASSWORD) in the `rockbin.conf` file. 

### Setting it up as an upstart service


```bash
# put the binary in the correct folder
cp .rockbin /usr/local/bin/
# edit rockbin.conf and set proper parameters to rockbin command
vi rockbin.conf
# put the upstart config file into the correct file
cp .rockbin.conf /etc/init/rockbin.conf
# reload the upstart configs
initctl reload-configuration
# start the service
service rockbin start
# On firmwares >2008 you will have to restart the vacuum
sudo reboot now
```

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

## Troubleshooting connection issues


If you're having difficulty connecting to the mutt server, it is possible to run the rock bin app manually with debug login enabled. To do this you can do the following assuming you have copied the binary to the `/usr/local/bin/` (`cp .rockbin /usr/local/bin/`) folder. 

1. Make sure the app is callable (installed correctly). 
```bash
rockbin --help
```

2. Test the connection to the mqtt server (change the mqtt_server address accordingly). 
```bash
rockbin -mqtt_server mqtt://192.168.0.144:1883 -log_level debug
```

If you require a username and password to connect to the mqtt_server, you can pass these using: 
```bash
MQTT_USERNAME=mqttuser MQTT_PASSWORD='Some%!Strong$Pass' rockbin -mqtt_server mqtt://192.168.0.144:1883 -log_level debug
```

3. If the above works correctly, you will need to update [rockbin.conf](rockbin.conf) accordingly. If you're using authentication a sample of the required [rockbin.conf](rockbin.conf) would look like 

```text
description     "rockbin mqtt publisher for the bin"
start on filesystem and net-device-up IFACE=wlan0
stop on runlevel [!2345]
respawn
umask 022
setuid root
setgid root
console log
env MQTT_USERNAME=mqttuser
env MQTT_PASSWORD='Some%!Strong$Pass'
script 
    exec /usr/local/bin/rockbin -mqtt_server mqtt://192.168.0.144:1883 -full_time 2400
end script
```

If you don't require a username or password, then leave these commented out.  
```text
#env MQTT_USERNAME=mqttuser
#env MQTT_PASSWORD='Some%!Strong$Pass'
```

If the upstart script is not working, but step 2 was working correctly, add the `-log_level debug` flag to the upstart script for more logging information. 

4. After adding the configuration script. Reboot the vacuum: 
```bash
sudo reboot now
```

5. Once the system is back up you should see the service in the list provided by: 
```bash
ps aux
```

6. If you had enabled debug mode earlier. You can check the logs for the output

```bash
cat /var/log/upstart/rockbin.log
```

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FjohnDorian%2Frockbin.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FjohnDorian%2Frockbin?ref=badge_large)

