#!/bin/sh

# Script per inizializzare la directory di sviluppo
# UTILIZZO: ./misc/setup.sh

install_hooks() {
  rm -rf .git/hooks
  cd .git
  ln -s ../misc/git-hooks hooks
  cd ..
}

install_hooks
