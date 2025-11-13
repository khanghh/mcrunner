#!/bin/sh
set -e

# Ensure timezone is configured
if [ -n "$TZ" ]; then
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime
    echo "$TZ" > /etc/timezone
fi

# Prepare log file
mkdir -p /minecraft/logs
touch /minecraft/logs/cron.log
chmod 666 /minecraft/logs/cron.log

# Install crontab if provided
if [ -f /minecraft/crontab ]; then
    tail -c1 /minecraft/crontab | read -r _ || echo >> /minecraft/crontab
    cp /minecraft/crontab /etc/cron.d/mcrunner
    chmod 0644 /etc/cron.d/mcrunner
fi

# Link system cron logs to minecraft/logs/cron.log
ln -sf /minecraft/logs/cron.log /var/log/cron.log

# Start cron in background
crond

echo "Starting mcrunner with args: $@"
exec /usr/bin/mcrunner "$@"
