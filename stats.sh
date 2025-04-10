#! /bin/bash

NoColor='\033[0;0m'
Black='\033[0;30m'
Red='\033[0;31m'
Green='\033[0;32m'
Browne='\033[0;33m'
Orange='\033[0;33m'
Blue='\033[0;34m'
Purple='\033[0;35m'
Cyan='\033[0;36m'
LightGray='\033[0;37m'
DarkGray='\033[1;30m'
LightRed='\033[1;31m'
LightGreen='\033[1;32m'
Yellow='\033[1;33m'
LightBlue='\033[1;34m'
LightPurple='\033[1;35m'
LightCyan='\033[1;36m'
White='\033[1;37m'

get_locs() {
    find . -type f -name "$1" | xargs wc -l
}

print_stats() {

    echo -e $Yellow":::: $2 ::::"$DarkGray
    get_locs "$1" | head -n -1
    echo -e $LightGreen$(get_locs $1 | tail -n 1)$NoColor"\n"
}

print_stats "*.go" "GO"
print_stats "*.sql" "SQL"
print_stats "*.sh" "SCRIPTS"
