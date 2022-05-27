#!/bin/bash

# Run all tests: make test
# Run specific tests: make test entity=<Entity name> 

ENTITY_ARG=${1:-""}
if [ -n "$ENTITY_ARG" ]; then 
    EntityArray=( $ENTITY_ARG )
else 
    EntityArray=("Feed" "Contact" "Follow" "Monitoring" "Auth" "Wallet" "Chat")
fi

# Iterates through entities for which we have a 
# test suite.
# For each entity, drops the test db and recreates it,
# then runs that specific test suite.
# If there is a failure, record it and output the list
# of failed entities in the Summary section.
go clean -testcache

FailArray=()
for entity in ${EntityArray[@]}; do
    echo "Preparing for $entity Tests"
    # convert entity name to lowercase to match directory names
    entityLower=`echo "$entity" | tr '[:upper:]' '[:lower:]'`
    PGPASSWORD=pass dropdb postgres-test --if-exists -U postgres
    PGPASSWORD=pass createdb postgres-test -U postgres
    OUTPUT=$(go test -v -p 1 github.com/VamaSingapore/vama-api/cmd/test/$entityLower/... -run "^Test$entity")
    echo "$OUTPUT"
    # Count instances of FAIL in the output
    FailCount=$(echo $OUTPUT | grep "FAIL"| wc -l)
    if [ $FailCount -gt 0 ]; then
        FailArray+=($entity)
    fi
    echo -------------------------------------------------------------------------------
done

echo "SUMMARY"
if [ ${#FailArray[@]} -gt 0 ]; then
    for entity in ${FailArray[@]}; do
        echo "$entity tests FAILED"
    done

    exit 1
else 
    echo "ALL TESTS PASSED"
fi 