# this test compares current exporter output to the etalon one
# (master at the moment of the test introduction) to quickly
# find out if there's any changes in metrics after upstream update
# exporter is started using default percona flags (taken from pmm-managed)
# TODO try to figure out testing of metrics type change

trap "exit" INT TERM ERR
trap "kill 0" EXIT

NFLAGS=""

for line in $(cat test.exporter-flags.txt)
    do
        NFLAGS+=" $line"
    done

NFLAGS+=" --web.listen-address=127.0.0.1:20001"

../../../node_exporter $NFLAGS &

cat /dev/null > metrics.upstream.txt
sleep 1s # wait for exporter process to fully start
curl http://localhost:20001/metrics --output metrics.upstream.txt

# cleanup comments
#sed '/^#/d' metrics.upstream.txt > metrics.upstream.no-comments.txt
#sed '/^#/d' metrics.percona.txt > metrics.percona.no-comments.txt

splitNames()
{
    cat /dev/null > "$2"

    while IFS="" read -r p || [ -n "$p" ]
    do
        if [[ "$p" =~ ^\#.* ]];
        then
            echo "$p" >> "$2"
        else
            IFS='{ ' read -r -a array <<< "$p"
            echo "${array[0]}" >> "$2"
        fi
    done < "$1"
}

splitNames metrics.upstream.txt metrics.upstream.no-values.txt
splitNames metrics.percona.txt metrics.percona.no-values.txt

git diff --exit-code --no-index -- metrics.percona.no-values.txt metrics.upstream.no-values.txt
