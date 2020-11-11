#!/bin/sh

###############################################
## This file is copied from https://raw.githubusercontent.com/zvldz/vacuum/46ea4810d10e8b5bc4957c15e23981ebf708bbc9/custom-script/files/valetudo-daemon.sh
## Several modifications (commented below) have been made to the original file. 
###############################################


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
        echo "Running Rockbin"
        echo '|/bin/false' > /proc/sys/kernel/core_pattern
        if [ -f "/root/bin/busybox" ]; then
            (
                # Make valetudo very likely to get killed when out of memory
                echo 1000 > /proc/self/oom_score_adj
                # Also run it with absolutely lowest CPU and I/O priority to not disturb anything critical on robot
                # modified the next lines for rockbin instead of valetudo
                #exec /root/bin/busybox ionice -c3 nice -n19 MQTT_USERNAME=xxxx MQTT_PASSWORD=xxxx  /usr/local/bin/rockbin -mqtt_server mqtt://192.168.0.144:1883 -full_time 2400 >> /var/log/upstart/valetudo.log 2>&1
                exec /root/bin/busybox ionice -c3 nice -n19 /usr/local/bin/rockbin -mqtt_server mqtt://192.168.0.144:1883 -full_time 2400 >> /var/log/upstart/valetudo.log 2>&1
            )
        else
            # modified the next line for rockbin
            #MQTT_USERNAME=xxxx MQTT_PASSWORD=xxxx nice -n 19 /usr/local/bin/rockbin -mqtt_server mqtt://192.168.0.144:1883 -full_time 2400 >> /var/log/upstart/valetudo.log 2>&1
            nice -n 19 /usr/local/bin/rockbin -mqtt_server mqtt://192.168.0.144:1883 -full_time 2400 >> /var/log/upstart/valetudo.log 2>&1            
        fi
    else
        echo "Waiting for SysUpdate to finish..."
    fi
done