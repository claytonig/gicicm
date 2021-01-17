#!/usr/bin/env bash

function tag {
  gh auth login --with-token < $TOKEN
  gh release create $RELEASE
}

tag
