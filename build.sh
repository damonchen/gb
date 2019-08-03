#!/bin/bash

set -e

get_osname(){
    echo $(uname -s | awk '{print tolower($0)}')
}

get_arch() {
  local a=$(uname -m)
  case ${a} in
    "x86_64" | "amd64" )
        echo "amd64"
        ;;
    "i386" | "i486" | "i586")
        echo "386"
        ;;
    *)
        echo ${NIL}
        ;;
    esac
}


main() {
  local os=$(get_osname)


}


main