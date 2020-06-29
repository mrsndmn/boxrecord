#!/bin/bash -ex

rm -rf /tmp/o3

if [ -f /tmp/o3/box.pid ]; then
    kill "$(cat /tmp/o3/box.pid)" || true
fi

mkdir -p /tmp/o3/work
mkdir -p /tmp/o3/snap
mkdir -p /tmp/o3/xlogs

octopus3_box --init-storage -c box.cfg
octopus3_box -c box.cfg
