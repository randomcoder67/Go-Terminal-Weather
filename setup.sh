#!/usr/bin/env sh

currentDir=$(pwd)

if [[ "$1" == "install" ]]; then
	[ -d ~/.local/bin ] || echo "~/.local/bin does not exist"
	cp "${currentDir}/weather" ~/.local/bin/weather
	echo "Install complete"
else
	if [ -z ${XDG_DATA_HOME+x} ]; then
		echo '$XDG_DATA_HOME is not set, would you like to set it to default (~/.local/share)? (y/N)'
		read yesOrNoSetXDGDATAHOME
		if [ "$yesOrNoSetXDGDATAHOME" == "y" ]; then
			XDG_DATA_HOME="$HOME/.local/share"
			echo 'export XDG_DATA_HOME='"\"$HOME/.local/share\"" >> "$HOME/.profile"
		else
			exit
		fi
	fi
	echo "Enter location for locations.csv file (leave blank for $XDG_DATA_HOME):" 
	read locLocation
	echo "package main" >> userSetValues.go
	echo "" >> userSetValues.go
	if [ "$locLocation" == "" ]; then
		echo "const LOCATIONS_FILE_LOCATION = $XDG_DATA_HOME" >> userSetValues.go
	else
		echo "const LOCATIONS_FILE_LOCATION = $locLocation" >> userSetValues.go
	fi
	echo "$locLocation"
	go mod tidy
	if go build -o weatherGo main.go bbc.go metoffice.go definitions.go userSetValues.go; then
		echo "Compile complete"
		rm userSetValues.go
	else
		echo "Compile failed"
	fi
fi
