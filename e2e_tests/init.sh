#!/bin/bash

init_minimal_no_git() {
    output=`cfg init --no-git`
    status="$?"

    expect_status='0'
    expect="Initialized empty base in $(pwd)"

    contents=`ls -a | xargs`
    expect_contents='. .. .cfg .cfg.json .gitignore'
    if [ ! "$contents" = "$expect_contents" ]; then
        printf 'FAIL. expected cfg_tests/ to have contents:\n%s\ngot:\n%s\n' "$expect_contents" "$contents"
        exit 1
    fi
}

init_minimal_inside_git() {
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
}

tests=(
    "init_minimal_no_git"
    "init_minimal_inside_git"
)
