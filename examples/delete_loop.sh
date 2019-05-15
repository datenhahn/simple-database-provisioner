#!/bin/bash

for i in {1..10}; do
  sed "s/COUNTER/$i/g" sdi-template.yaml | kubectl delete -f -
  sed "s/COUNTER/$i/g" sdb-template.yaml | kubectl delete -f -
done
