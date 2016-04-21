#!/bin/bash

amount_of_rows=2349987
number_of_rows=10000
offset=0
limit=number_of_rows


while [ $offset -lt $amount_of_rows ]; do
  ./sponge import --offset $offset --limit $limit >> sponge.log &
  let offset=$limit
  let limit=limit+number_of_rows
done
