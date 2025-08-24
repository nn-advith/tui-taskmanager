#!/bin/bash

DEFAULT_BINDING=${1:-"N"}
CUSTOM_SHORTCUTS=$(gsettings get org.gnome.settings-daemon.plugins.media-keys custom-keybindings)
# last=$(echo ${CUSTOM_SHORTCUTS} | grep -o "'[^']*'" | tail -n1 | tr -d "'" | grep -oE "custom[0-9]+" | grep -oE "[0-9]+")
# next=$(( ${last} + 1))
if $1; then
    echo -e "Using default binding \033[32m<Ctrl><Alt>N\033[0m"
else
    echo -e "Using custom binding \033[32m<Ctrl><Alt>$1\033[0m"
fi

CUSTOM_SHORTCUTS=$(echo $CUSTOM_SHORTCUTS | tr -d "[] ")
declare -a carray=()
IFS=','

for i in $CUSTOM_SHORTCUTS; do 
    carray+=($(echo $i | tr -d "'"))
done

for i in "${carray[@]}"; do
    binding=$(gsettings get org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${i} binding | tr -d "'")
    if [ "$binding" == "<Ctrl><Alt>$DEFAULT_BINDING" ]; then
        echo -e "binding \033[31m<Ctrl><Alt>N\033[0m is linked to :\t $i.\n Select another binding key."
        exit 1
    fi
done

next=$(( $(echo ${carray[${#carray[@]}-1]} | grep -oE "custom[0-9]+" | grep -oE "[0-9]+") +1 ))

ns="/org/gnome/settings-daemon/plugins/media-keys/custom-keybindings/custom${next}/" 
# nsu=$(echo $ns | sed 's/\//\\\//g')
# updated=$(echo "$CUSTOM_SHORTCUTS" | sed "s/]$/, '${nsu}']/")
carray=("${carray[@]}" "$ns")
updated=$(printf "'%s', " "${carray[@]}")
updated_string="[${updated%, }]"
echo -e "\nAdding new custom shortcut for Task Manager v 9000 ..."
gsettings set org.gnome.settings-daemon.plugins.media-keys custom-keybindings "${updated_string}"

gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${ns} name 'Task  v9000'
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${ns} command "gnome-terminal --geometry=67x30 -- $HOME/go/bin/tui-taskmanager"
gsettings set org.gnome.settings-daemon.plugins.media-keys.custom-keybinding:${ns} binding "<Ctrl><Alt>${DEFAULT_BINDING}"

echo -e "\033[32mDone!\033[0m"