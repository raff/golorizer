golorizer
=========

Log colorizer in Go

    usage:
        tail -f server.log | golorizer [options]

### Options:

+ --levels=true/false : enabled/disable automatic coloring by log level (default true)
+ --color=background:foreground : set current color for custom pattern
+ --custom=regexp-pattern : add custom pattern to match and color


    example:
        tail -f server.log | golorizer --custom=EXTERNAL --color purple:black --custom=elapsed


colorizes log entries containing the word "EXTERNAL" with default custom color (white on black)
and entries containing the word "elapsed" with specified color (purple on black)

also colorize log entries according to their log level (green for INFO, cyan for DEBUG, yellow for WARN, red for ERROR, etc.)
