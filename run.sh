#!/bin/bash
docker run --rm --name mcrunner -v $PWD/build/bin/mockserver:/mockserver -v $PWD/minecraft:/minecraft registry.mineviet.com/mcrunner:web -c "/mockserver"
