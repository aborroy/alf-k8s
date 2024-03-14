#! /bin/bash

checkKontext() {
	local ACTION_MSG=${1:-"You're about to apply modifications"}
	HAS_ONE_CONTEXT_ONLY=$(kubectl config get-contexts --no-headers=true | wc -l)
	case $HAS_ONE_CONTEXT_ONLY in
		0)
			echo "No Kubernetes context is set. I can't proceed"
			exit 7
			;;
		1)
			return 0
			;;
		*)
			KLUSTER=$(kubectl config get-contexts --no-headers=true | awk '/^\*/{print $2}')
			read -n 1 -p "$ACTION_MSG on Kubernetes cluster $KLUSTER. Is that really what you want to do? (y/n): "
			[ $REPLY != "y" ] && exit 7
			return 0
			;;
	esac
}
