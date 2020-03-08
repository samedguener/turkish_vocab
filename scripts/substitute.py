#!/usr/bin/python

import argparse

parser = argparse.ArgumentParser(description='Substitute variables in a yaml-file.')
parser.add_argument('input', help='Name of the input filename')
parser.add_argument('--output', help='Name of the output files (Default: Overwrite input file)')
parser.add_argument('--values', nargs='+', help='Values to be replaced and replaced with (e.g environment="staging")')

args = parser.parse_args()

with open(args.input) as f:
    newText = f.read()

for value in args.values:
    key, value = value.split("=")        
    newText=newText.replace("$" + key, value)

if args.output is None:
    with open(args.input, "w") as f:
        f.write(newText)
else:
    with open(args.output, "w") as f:
        f.write(newText)