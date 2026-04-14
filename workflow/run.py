import sys

from seq import loader

args = sys.argv[1:]
# fmt is set per-keyword via Alfred variable.
# seq fmt passes empty fmt; user provides the format string as first arg.
fmt, narr = (args[1], args[2:]) if args[0] == "" else (args[0], args[1:])
sys.stdout.write("\r\n".join(list(loader(fmt, narr))))
