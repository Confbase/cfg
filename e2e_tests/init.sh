#!/bin/bash

init_minimal_no_git() {
    if [ -d "cfg_tests" ]; then
        printf "FAIL. Could not setup test. ./cfg_tests/ already exists.\n"
        exit 1
    fi

    mkdir cfg_tests
    pushd cfg_tests >/dev/null

    output=`cfg init --no-git`
    status="$?"

    expect_status='0'
    expect="Initialized empty base in $(pwd)"

    popd >/dev/null
    rm -r cfg_tests
}

init_minimal_inside_git() {
    if [ -d "cfg_tests" ]; then
        printf "FAIL. Could not setup test. ./cfg_tests/ already exists.\n"
        exit 1
    fi

    mkdir cfg_tests
    pushd cfg_tests >/dev/null

    output=`cfg init 2>&1`
    status="$?"

    expect_status='1'
    expect='error: the current directory is inside a git repository
if you wish to use cfg in conjuction with an existing git repository,
consider running '"'"'cfg init --no-git'"'"'
rolling back changes'

    contents=`ls -a | xargs`
    if [ ! "$contents" = '. ..' ]; then
        printf 'FAIL. expected cfg_tests/ to be empty; got this:\n%s' "$contents"
        exit 1
    fi

    popd >/dev/null
    rm -r cfg_tests
}

tests=(
    "init_minimal_no_git"
    "init_minimal_inside_git"
)
