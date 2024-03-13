#! /bin/bash

checkKontext() {
	local ACTION_MSG=${1:-"You're about to apply modifications"}
	KLUSTER=$(kubectl config get-contexts --no-headers=true | awk '/^\*/{print $2}')
	read -n 1 -p "$ACTION_MSG on Kubernetes cluster $KLUSTER. Is that really what you want to do? (y/n): "
	[ $REPLY != "y" ] && exit 7
}
