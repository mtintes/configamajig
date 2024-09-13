#!/bin/sh

for i in `seq 1 100`
do
    ./bin/configamajig replace --config ./example-files/config.json --input ./example-files/input1.txt --output ./many-outs/output$i.txt
done