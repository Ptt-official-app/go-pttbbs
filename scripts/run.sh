#!/bin/bash

ini_filename=02-config.run.ini

go build -tags dev && ./go-pttbbs -ini ${ini_filename}
