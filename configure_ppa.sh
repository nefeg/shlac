#!/usr/bin/env bash
# project root
#   |_src
#   |_WORK_DIR
#       |_SOURCE_DIR
#       |_TMP_DIR
#       |_BUILD_NAME.tar.gz
#
#

#
# Usage ./this-script <project name> <version> <isTest> <project source>
# Example: ./this-script shlac 0.25-0ubuntu1 0 https://verefkin@bitbucket.org/verefkin/shlac.git
#

# date ; sudo service ntp stop ; sudo ntpdate -s time.nist.gov ; sudo service ntp start ; date
#
# *** DIR: %project root%
#

###################################
#   checking version              #
###################################
if [[ $1 == '' ]];then
    echo "expected version (ex.: 0.25-0ubuntu1)"
    exit 1
fi

PROJECT_NAME=$1
VERSION="$2"
IS_TEST="$3"
PROJECT_SOURCE=$(echo $4|sed 's/\//\\\//g')

PROJECT_ROOT=$(pwd)
BUILD_NAME="${PROJECT_NAME}-${VERSION}"
BUILD_NAME_ALT="${PROJECT_NAME}_${VERSION}"
WORK_DIR=$(pwd)"/build"
TMP_DIR=$WORK_DIR/tmp
SOURCE_DIR=$WORK_DIR/source

export DEBEMAIL="evgeny.nefedkin@umbrella-web.com"
export DEBFULLNAME="Evgeny Nefedkin"


echo "$PROJECT_NAME ($BUILD_NAME_ALT)"

###################################
#   making workspace              #
###################################
#git archive --format=tar.gz -o $WORK_DIR/$BUILD_NAME.tar.gz ppa
## copy sources
mkdir -p $SOURCE_DIR
cp -rf src/ $SOURCE_DIR/src/
cp -rf $SOURCE_DIR/src/vendor/* $SOURCE_DIR/src/

cp Makefile $SOURCE_DIR/Makefile
cp config.json $SOURCE_DIR/config.json
cp readme.md $SOURCE_DIR/readme.md

## remove .git's directories
find $SOURCE_DIR/ -name .git | rm -fr

## make archive
#
# *** DIR: %project root%/$WORK_DIR/$SOURCE_DIR
#
cd $SOURCE_DIR
tar -czf ../$BUILD_NAME.tar.gz ./*


###################################
#   make skeleton                 #
###################################
mkdir $TMP_DIR -p
cd $TMP_DIR;
dh_make -f ../$BUILD_NAME.tar.gz -s -e ${DEBEMAIL} -c gpl2 -y --createorig -p $BUILD_NAME_ALT


###################################
#   edit 'debian/control' file    #
###################################
sed -i "s/#Vcs-Browser: http:\/\/git.debian.org\/?p=collab-maint\/${PROJECT_NAME}.git;a=summary/Vcs-Browser: ${PROJECT_SOURCE}/" debian/control
sed -i "s/#Vcs-Git: git:\/\/git.debian.org\/collab-maint\/${PROJECT_NAME}.git/Vcs-Git: ${PROJECT_SOURCE}/" debian/control
sed -i "s/<insert the upstream URL, if relevant>/${PROJECT_SOURCE}/" debian/control
#sed -ri 's/^Depends: .*/Depends: /' debian/control
sed -i 's/<insert up to 60 chars description>/Distributed and concurrency job manager/' debian/control
sed -i 's/<insert long description, indented with spaces>/Distributed and concurrency job manager with cron syntax/' debian/control
sed -i 's/Section: unknown/Section: utils/' debian/control
sed -i 's/Build-Depends: debhelper (>= 8.0.0)/Build-Depends: debhelper (>= 9.0.0), golang (>= 1.9)/' debian/control
sed -i 's/Standards-Version: 3.9.4/Standards-Version: 3.9.5/' debian/control


###################################
#   edit 'debian/copyright' file  #
###################################
sed -i "s/<url:\/\/example.com>/${PROJECT_SOURCE}/" debian/copyright
sed -i "s/<years> <put author.s name and email here>/${DEBFULLNAME} <${DEBEMAIL}>/" debian/copyright
sed -i 's/<years> <likewise for another author>/\n/' debian/copyright
sed -ri 's/^#.+//' debian/copyright

###################################
#   edit 'debian/changelog' file  #
###################################

if [[ -e ${TMP_DIR}/debian/changelog ]];then

    cp -f ${PROJECT_ROOT}/changelog ${TMP_DIR}/debian/changelog

else
    if [[ "$IS_TEST" != "1" ]];then
        vi debian/changelog
    fi
fi

#sed -i 's/unstable;/trusty;/' debian/changelog
#
#if [[ "$IS_TEST" != "1" ]];then
#    vi debian/changelog
#fi


## Remove useless files
rm debian/*.ex &>/dev/null
rm debian/*.EX &>/dev/null
#
### build package (https://help.launchpad.net/Packaging/PPA/BuildingASourcePackage)
#debuild -S -sa
#
### Upload to PPA
#dput -d ppa:onm/shlanc ${WORK_DIR}/${BUILD_NAME_ALT}-1_source.changes
#
#rm -rf $WORK_DIR

