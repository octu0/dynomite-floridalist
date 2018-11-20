#!/bin/bash

curl -s localhost:2101/cluster_describe | python -mjson.tool
