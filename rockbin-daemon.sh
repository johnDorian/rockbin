#!/bin/sh
# from https://github.com/zvldz/vacuum/commits/master/custom-script/files/valetudo-daemon.sh
# created by:  zvldz aka vtel
# This file goes to /usr/local/bin/rockbin-daemon.sh

mkdir -p /var/log/upstart

while :; do
    sleep 5
    if [ `cut -d. -f1 /proc/uptime` -lt 300 ]; then
        echo -n "Waiting for 20 sec after boot..."
        sleep 20
        echo " done."
    fi

    pidof SysUpdate > /dev/null  2>&1
    if [ $? -ne 0 ]; then
    echo "Running Valetudo"
    # Use "quotes" if user or password contains special chars.
    MQTT_USERNAME=xxxx MQTT_PASSWORD=xxxx /usr/local/bin/rockbin -mqtt_server mqtt://xxx.xxx.xxx.xxx:1883 -sensor_name vacuum_name_trash_box -full_time 2400 
    else
    echo "Waiting for SysUpdate to finish..."
    fi
done