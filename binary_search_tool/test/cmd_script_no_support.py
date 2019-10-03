"""Command script without compiler support for pass level bisection.

This script generates a pseudo log when certain bisecting flags are not
supported by compiler.
"""

from __future__ import print_function

import os
import sys


def Main():
  if not os.path.exists('./is_setup'):
    return 1
  print('No support for -opt-bisect-limit or -print-debug-counter.',
        file=sys.stderr)
  return 0


if __name__ == '__main__':
  retval = Main()
  sys.exit(retval)
