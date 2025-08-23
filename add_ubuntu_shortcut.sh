#!/bin/bash

# TODO: Add check if shortcut isussed already

CUSTOM_SHORTCUTS=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings)
last=$(echo ${CUSTOM_SHORTCUTS} | grep -o "'[^']*'" | tail -n1 | tr -d "'" | grep -oE "custom[0-9]+" | grep -oE "[0-9]+")
next=$(( ${last} + 1))

ns="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom${next}/" 
nsu=$(echo $ns | sed 's/\//\\\//g')
updated=$(echo "$CUSTOM_SHORTCUTS" | sed "s/]$/, '${nsu}']/")
# echo $updated
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "${updated}"

gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${ns} name 'Task  v9000'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${ns} command "gnome-terminal --geometry=67x30 -- $HOME/go/bin/tui-taskmanager"
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${ns} binding '<Ctrl><Alt>N'
