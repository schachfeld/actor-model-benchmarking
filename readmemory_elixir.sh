while true; do
ps -C beam.smp -o pid= -o vsz= -o rss= >> mem_beam.log
# gnuplot show_mem.plt
sleep 0.5
done &