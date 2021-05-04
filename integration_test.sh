#!/usr/bin/env bash

set -e

rm -rf test/
mkdir -p test/tmp/

S3_BUCKET=tt-development-nathanw-test
RETENTION_ARGS="--backend local --dir test"

if [[ $S3 ]]; then
    aws s3 rm s3://$S3_BUCKET --recursive
    RETENTION_ARGS="--backend s3 --bucket $S3_BUCKET"
fi

assert_exists() {
    if [[ $S3 ]]; then
        if ! aws s3 ls s3://$S3_BUCKET/$1; then
            echo "File $1 should exist in $S3_BUCKET but doesn't"
            exit 1
        fi
    else
        if test ! -f "test/$1"; then
            echo "File test/$1 should exist but doesn't"
            exit 1
        fi
    fi
}

assert_not_exists() {
     if [[ $S3 ]]; then
        if aws s3 ls s3://$S3_BUCKET/$1; then
            echo "File $1 shouldn't exist in $S3_BUCKET but does"
            exit 1
        fi
    else
        if test -f "test/$1"; then
            echo "File test/$1 should not exist but does"
            exit 1
        fi
    fi
}

create_file() {
    if [[ $S3 ]]; then
        touch test/tmp/hello
        aws s3 cp test/tmp/hello s3://$S3_BUCKET/$1
    else
        touch test/$1
    fi
}

create_file a
./backup-retention --period daily --num 5 $RETENTION_ARGS
sleep 1
create_file b
./backup-retention --period daily --num 5 $RETENTION_ARGS
sleep 1
create_file c
./backup-retention --period daily --num 5 $RETENTION_ARGS
sleep 1
create_file d
./backup-retention --period daily --num 5 $RETENTION_ARGS
sleep 1
create_file e
./backup-retention --period daily --num 5 $RETENTION_ARGS
sleep 1

assert_exists daily/a
assert_exists daily/b
assert_exists daily/c
assert_exists daily/d
assert_exists daily/e

create_file f
./backup-retention --period daily --num 5 $RETENTION_ARGS
assert_exists daily/f
assert_not_exists daily/a

./backup-retention --num 5 $RETENTION_ARGS
assert_not_exists a
