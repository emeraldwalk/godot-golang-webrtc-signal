# Golang WebRTC Signaling Server

## Systemd Config

1. Open a new config file for the service:
   ```sh
   sudo nano /etc/systemd/system/gowebrtcsignal.service`
   ```
1. Copy / paste the following template filling in <user> and <appfolder> placeholders

   ```ini
   [Unit]
   Description = Go WebRTC Signaling Server

   [Service]
   Type           = simple
   Restart        = always
   RestartSec     = 5s
   StandardOutput = append:%h/<appfolder>/stdout.log
   StandardError  = append:%h/<appfolder>/stderr.log
   ExecStart      = ./signalserver
   WorkingDirectory = %h/<appfolder>
   User = <user>

   [Install]
   WantedBy = multi-user.target
   ```

1. Enable and start the service

```sh
sudo systemctl enable gowebrtcsignal.service
sudo systemctl start gowebrtcsignal.service
sudo systemctl status gowebrtcsignal.service
```
