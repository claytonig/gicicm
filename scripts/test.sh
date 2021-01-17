#!/usr/bin/env bash

function tag {
  gh auth login
  gh release create $RELEASE
}

tag
