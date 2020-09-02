#!/bin/sh

tidy -qimc -w 90 --output-bom no --tidy-mark no --clean yes --logical-emphasis yes --indent-attributes yes --vertical-space yes $1/*.html

tidy -qimc -w 90 -xml --output-bom no --clean yes --logical-emphasis yes --indent-attributes yes --vertical-space yes $xmlfiles $1/*.xml
