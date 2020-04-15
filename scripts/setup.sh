#!/bin/sh

# Script per inizializzare la directory di sviluppo
# UTILIZZO: ./scripts/setup.sh

install_hooks() {
  rm -rf .git/hooks
  cd .git
  ln -s ../githooks hooks
  cd ..
}

install_hooks
