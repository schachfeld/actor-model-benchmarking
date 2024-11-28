while true; do
ps -C main -o pid= -o vsz= -o rss= >> mem.log
# gnuplot show_mem.plt
sleep 0.5
done &