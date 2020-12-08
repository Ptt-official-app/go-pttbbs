#!/usr/bin/env python
# -*- coding: utf-8 -*-

import yaml
import json
import sys

src = sys.argv[1]
dst = sys.argv[2]

with open(src, 'r') as f:
    the_struct = yaml.load(f)

with open(dst, 'w') as f:
    json.dump(the_struct, f)
