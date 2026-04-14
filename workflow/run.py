import sys

from seq import loader

sys.stdout.write("\r\n".join(list(loader(sys.argv[1], sys.argv[2:]))))
