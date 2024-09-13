# Content

- Modify the service file
- Copy the Service File to the Systemd Directory
- Reload Systemd Daemon
- Enable the Service
- Start the Service
- Check the Service Status
- Reload Configuration on SIGHUP
- See logs


# Modify the service file

``` 
[Service]
User=icalgit
Group=icalgit

Type=simple
WorkingDirectory=/home/icalgit/icalgit
ExecStart=/home/icalgit/bin/ical-git
Restart=on-failure
TimeoutSec=10
```

Change the `user` and `group` to your own user.
Make sure `WorkingDirectory` and `ExecStart` points to existing own paths.

# Copy the Service File to the Systemd Directory

You need to copy your .service file to the `/etc/systemd/system`
directory. This is where systemd looks for service files. You'll need superuser
privileges to do this. 

```console
sudo cp ical-git.service /etc/systemd/system/
```

# Reload Systemd Daemon: 

After adding or modifying any service file, you need to reload the systemd daemon so it can recognize the new or changed service file.
```console
sudo systemctl daemon-reload
```

# Enable the Service

If you want your service to start automatically on boot, you should enable it.
This step is optional but recommended for most services that you want to run
continuously.

```console
sudo systemctl enable ical-git.service 
```

# Start the Service

Now, you can start your service. If you've enabled it, this step is technically
optional since it will start on the next boot, but you'll likely want to start
it immediately for testing purposes.

```console
sudo systemctl start ical-git.service 
```

# Check the Service Status 
To verify that your service is running as expected, you can check its status.

```console
sudo systemctl status ical-git.service 
```

# Reload Configuration on SIGHUP
`ical-git` can reload its configuration file when it receives the SIGHUP signal. To reload the configuration without restarting the service, use the following command:

```console
sudo systemctl reload ical-git.service
```

# See logs

```console
sudo journalctl -u ical-git.service 
```
