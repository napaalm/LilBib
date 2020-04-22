#!/bin/sh

# Script per inizializzare la directory di sviluppo
# UTILIZZO: ./scripts/setup.sh

install_hooks() {
  rm -rf .git/hooks
  cd .git
  ln -s ../githooks hooks
  cd ..
}

[ -d .git ] || { echo 'Execute this script from the project folder'; exit 1; }
install_hooks
