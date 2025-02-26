#!/bin/bash

function listDirInDir() {
    [[ $1 == "-h" || -z $1 ]] && {
        echo "Usage: sudo $FUNCNAME <path/to/dir>"
        return 0
    }
    dirPath=${1}
    dirs=()
    for element in ${dirPath}/*; do
        [[ -d ${element} ]] && {
            dirs+="${element} "
        }
    done
    echo "${dirs[@]}"
    return 0
}

dirs=$(listDirInDir .)
#echo "$dirs"
read -ra dirArray  <<< "$dirs"
current_dir=$(pwd)
#echo "${dirArray[@]}"

for dir in "${dirArray[@]}"; do
	echo  "update modules: $dir"
	cd "${current_dir}" && {
	  cd "$dir" && go get -u && go mod tidy #&
	}
done <<< "$dirArray"

#wait
