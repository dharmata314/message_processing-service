#!/bin/bash

/usr/local/bin/wait-for-it.sh kafka:9092 --timeout=30

/app
