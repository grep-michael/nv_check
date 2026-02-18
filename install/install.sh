#!/bin/bash

sudo cp nv_check /usr/local/bin/nv_check
sudo chmod +x /usr/local/bin/nv_check

mkdir -p ~/.config/autostart
cat > ~/.config/autostart/nv_check.desktop <<EOL
[Desktop Entry]
Type=Application
Name=GPU Monitor
Exec=nv_check
X-LXDE-Autostart-Phase=Application
NoDisplay=true
EOL

chown guest:guest ~/.config/autostart/nv_check.desktop