#!/bin/bash
# Parameter n - number of files to fetch in the directory

# 1. ssh into remote server and list the files that need to fetch
# 2. scp to local 

# path="/opt3/sedgwick/images/Main_timelapse/2016/11"
path="/opt3/sedgwick/images/zooniverse"
target="/Users/michaelzhang/Downloads/WTB_samples/time_lapse"

files=$(ssh heronalps@128.111.39.240 ls $path | head -n $1)

while IFS= read -r file; 
do

scp "heronalps@128.111.39.240:$path/$file" $target; 

done <<< "$files"
