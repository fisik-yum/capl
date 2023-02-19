#!/bin/bash

file="badapple.mp4"

ffmpeg -i $file -vf format=gray -filter:v scale=512:-1 -c:a copy output.mp4
framecount=$(ffprobe -v error -select_streams v:0 -count_frames -show_entries stream=nb_read_frames -print_format default=nokey=1:noprint_wrappers=1 output.mp4)
echo "$framecount"
if [ -d "frames" ] 
then
    rm -rf frames/*
else
	mkdir frames/
fi

#ffmpeg -i output.mp4 -r 1/1 frames/%03d.bmp
ffmpeg -r 1 -i output.mp4 -r 1 frames/%03d.png
