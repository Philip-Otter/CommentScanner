#!/bin/bash

while getopts "hu:ecp" flag; do
 case $flag in
   h)
   echo CommentScanner is a short script for finding information within web pages.
   echo By default it will search for HTML, CSS, and Javascript comments.
   echo Line numbers are listed in front of each match.
   echo
   echo Flags:
   echo '-h    Show help information'
   echo '-u    URL information'
   echo '-e    Search for emails'
   echo 
   echo Usage
   echo 'CommentScanner -u http://{site}        HTML CSS JS comments'
   echo 'CommentScanner -u http://{site} -ec     HTML CSS JS + emails + Credentials'

   ;;
   u)
   url=$OPTARG
   echo 'HTML COMMENTS:  '
   curl -s -L --user-agent firefox $url | grep -n '<!--'
   echo; echo 'JAVASCRIPT COMMENTS:  '
   curl -s -L --user-agent firefox $url | grep '//' | grep -v 'http://' | grep -v -n 'https://'
   echo; echo 'CSS COMMENTS:  '
   curl -s -L --user-agent firefox $url | grep -n '/\*'
   ;;
   e)
   echo; echo 'POSSIBLE EMAILS:  '
   curl -s -L --user-agent firefox $url | grep -P -i -o -n '[-A-Za-z0-9!#$%&'"'"'*+/=?^_`{|}~]+(?:\.[-A-Za-z0-9!#$%&'"'"'*+/=?^_`{|}~]+)*@(?:[A-Za-z0-9](?:[-A-Za-z0-9]*[A-Za-z0-9])?\.)+[A-Za-z0-9](?:[-A-Za-z0-9]*[A-Za-z0-9])?'
   ;;
   c)
   echo; echo 'POSSIBLE CREDENTIALS:  '
   curl -s -L --user-agent firefox $url | grep -i -E -n -C 2 --color 'username|password|logon|login|credentials|admin'
   ;;
   p)
   echo; echo 'POSSIBLE PHONE NUMBERS:  '
   curl -s -L --user-agent firefox $url | grep -P -i -n --color '\d\d\d-\d\d\d-\d\d\d\d|\(\d\d\d\)\s\d\d\d-\d\d\d\d' 
   ;;
   \?)
   echo 'Invalid flag usage! Please use -h to learn more!'
   ;;
 esac
done

