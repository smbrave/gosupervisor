#!/bin/sh

LAST_VERSION=2.1.17
SUB_VERSION=`echo $LAST_VERSION|awk -F"." '{print $3}'`
MID_VERSION=`echo $LAST_VERSION|awk -F"." '{print $2}'`
BIG_VERSION=`echo $LAST_VERSION|awk -F"." '{print $1}'`

((SUB_VERSION+=1))
if ((SUB_VERSION >= 100));then
    SUB_VERSION=0
    ((MID_VERSION+=1))
fi


if ((MID_VERSION>=10));then
    MID_VERSION=0
    ((BIG_VERSION+=1))
fi


NEW_VERSION="${BIG_VERSION}.${MID_VERSION}.${SUB_VERSION}"

echo ${NEW_VERSION}
sed --in-place "s%^LAST_VERSION=.*%LAST_VERSION=${NEW_VERSION}%" ${0}

git push
