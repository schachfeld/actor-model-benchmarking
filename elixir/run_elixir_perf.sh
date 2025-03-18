#!/bin/bash
export PATH=/home/valentin/.asdf/installs/elixir/main-otp-27/bin:/home/valentin/.asdf/installs/erlang/27.1.1/bin:$PATH
exec sudo -E chrt -f 99 perf stat -ddd elixir "$@"
