#!/usr/bin/env bash
set -e

WORK_DIR=.Adv360
BRANCH=${1:-V3.0}
INFO=${INFO:-true}

logInfo() {
  if ${INFO}; then
    echo -n -e '\e[1;33m'
    echo -n "$(date +%Y-%m-%dT%H:%M:%S) INFO "
    echo -n -e '\e[0m'
    echo "- $*"
  fi
}

logError() {
  if ${INFO}; then
    echo -n -e '\e[1;31m'
    echo -n "$(date +%Y-%m-%dT%H:%M:%S) ERROR "
    echo -n -e '\e[0m'
    echo "- $*"
  fi
}

## Check Pre-Requisites
if ! command -v jq &> /dev/null; then
  logError "jq could not be found, please install it:  brew install jq"
  exit
fi
if ! command -v perl &> /dev/null; then
  logError "perl could not be found, please install it"
  exit
fi
if ! command -v go &> /dev/null; then
  logError "golang could not be found, please install it:  brew install go"
  exit
fi
if ! command -v git &> /dev/null; then
  logError "git could not be found, please install it:  brew install git"
  exit
fi
if [[ ! -f macros.dtsi ]]; then
  cp macros-template.dtsi macros.dtsi
  logError "macros.dtsi was not found, created"
fi
if [[ ! -f keys.json ]]; then
  cp keys-template.json keys.json
  logError "keys.json was not found, created"
fi

## Get Adv360-Pro-ZMK repo
logInfo "Using BRANCH: $BRANCH"
if ! [ -d "${WORK_DIR}" ]; then
  logInfo ${WORK_DIR}/ dir does not exist, creating
  git clone --branch "${BRANCH}" git@github.com:KinesisCorporation/Adv360-Pro-ZMK.git ${WORK_DIR}
else
  logInfo Reverting ${WORK_DIR}/ dir
  cd ${WORK_DIR}
  git fetch origin
  git switch "${BRANCH}"
  git reset --hard "origin/${BRANCH}"
  cd ..
fi

logInfo Copying ${WORK_DIR}/config/macros.dtsi
cp macros.dtsi ${WORK_DIR}/config/macros.dtsi

WROTE_TEMPLATE=false
if ! [ -f "template.keymap" ]; then
  logInfo Creating template.keymap from ${WORK_DIR}/config/adv360.keymap
  ./createTemplateKeymap.pl ${WORK_DIR}/config/adv360.keymap > template.keymap
  WROTE_TEMPLATE=true
else
  logInfo Using an existing custom template.keymap
fi

logInfo Generating ${WORK_DIR}/config/keymap.json and ${WORK_DIR}/config/adv360.keymap
go run keymapGen.go keys.json
jq < keymap.json > ${WORK_DIR}/config/keymap.json
rm keymap.json
mv adv360.keymap ${WORK_DIR}/config/adv360.keymap

if ${WROTE_TEMPLATE}; then
  logInfo Removing temporary template.keymap
  rm template.keymap
fi

cd ${WORK_DIR}
  if [ -f "Makefile" ]; then
    # Force Docker
    sed -e '1d' Makefile > Makefile.tmp
    # shellcheck disable=SC2016 ## Don't want to execute command
    echo 'DOCKER := $(shell { command -v docker; })' > Makefile
    cat Makefile.tmp >> Makefile
    rm Makefile.tmp
  fi

  if [ -x "setup.sh" ]; then
    logInfo Running ${WORK_DIR}/setup.sh
    ./setup.sh
  else
    logInfo Running Make Clean
    make clean || true
  fi
  docker rm zmk || true

  if [ -x "run.sh" ]; then
    logInfo Running ${WORK_DIR}/run.sh
    ./run.sh
  else
    logInfo Running Make All
    make all
  fi
cd ..

cp ${WORK_DIR}/firmware/* firmware/
echo
ls -ltr firmware/

if [[ $(uname -s) == "Darwin" ]]; then
  open firmware
fi
echo
TZ=UTC date +%Y-%m-%dT%H:%M:%SZ
echo To Install Press MOD+① and copy the left file,
echo Then Press MOD+③ and copy the right file.
# shellcheck disable=SC2012 ## Want ls to sort by reverse date
ls -tr firmware/settings-reset.uf2  | sed -e 's/^/cp -X firmware\//;s/$/ \/Volumes\/ADV360PRO/'

# shellcheck disable=SC2012 ## Want ls to sort by reverse date
ls -tr firmware/ | tail -n 2 | sed -e 's/^/cp -X firmware\//;s/$/ \/Volumes\/ADV360PRO/'

logInfo Done