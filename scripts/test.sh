#!/usr/bin/env bash

function tag {
  gh release create $RELEASE
}

tag
