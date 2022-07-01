#!/bin/bash
DATE=`date +"%s"`

CURRENT=$DATE.bench
BENCHTIME=3s
BENCHCOUNT=3

if [ $# -eq 1 ]; then 
    BENCHCOUNT=$1
elif [ $# -eq 2 ]; then
    BENCHTIME=$2
fi

if [ -d ".bench" ]
then
    echo "Bench directory exists, skipping..."
else
    echo "Creating benches directory..."
    mkdir .bench
    echo "✔ Done"
fi

echo "Running benchmarks. This may take a while..."
go test -bench=. -count=$BENCHCOUNT -run=^@ -benchmem  -cpuprofile cpu.prof  -memprofile mem.prof -benchtime=$BENCHTIME > .bench/$CURRENT
echo "✔ Done"

cd .bench/

if [ -f "old.bench" ]
then
    echo "Unlinking previous benches..."
    rm old.bench
    echo "✔ Done"
fi

echo "Linking new benches..."
# The previous latest bench becomes the old one
if [ -f "latest.bench" ]
then
    ln latest.bench old.bench
    rm latest.bench
fi

# The newly created bench file becomes the latest one
ln $CURRENT latest.bench

# If no old.bench, latest and old are the same
if [ -f "old.bench" ]
then
    # do nothing
    echo "Existing old.bench file found, continuing..."
else
    echo "No old.bench file found. Linking latest.bench to old.bench..."
    ln latest.bench old.bench
fi
echo "✔ All Done. You can now compare the new and previous benches by running 'make benchcmp'. "
TIME=$((`date +"%s"` - $DATE))
echo "Successfull after $TIME seconds"

