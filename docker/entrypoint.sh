#!/bin/sh
set -e

# Ensure timezone is configured
if [ -n "$TZ" ]; then
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime
    echo "$TZ" > /etc/timezone
fi

# Prepare log file
touch /minecraft/cron.log
chmod 666 /minecraft/cron.log

# Install crontab if provided
if [ -f /minecraft/crontab ]; then
    tail -c1 /minecraft/crontab | read -r _ || echo >> /minecraft/crontab
    cp /minecraft/crontab /etc/cron.d/mcrunner
    chmod 0644 /etc/cron.d/mcrunner
fi

# Redirect system cron logs to same file
ln -sf /minecraft/cron.log /var/log/cron.log

# Start cron in background
cron

exec /usr/bin/mcrunner "$@"
