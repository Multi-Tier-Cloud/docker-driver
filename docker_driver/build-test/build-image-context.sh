#!/bin/bash
gcc -o test-image/test -static test.c
tar -C test-image -cvf test-image.tar Dockerfile test
