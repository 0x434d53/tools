# change all config in .git/:q
# thrnage the url from 
# [remote "origin"]
#	url = git@github.com:
# to 
#	url = https://github.com/

from pathlib import Path
import fileinput
import re
import sys


def change_config(p):
	print(str(p))
	f = fileinput.FileInput(str(p / 'config'), inplace=True)

	for line in f: 
		print(line.replace('https://github.com/', 'git@github.com:'),end='')

	fileinput.close()

def traverse(p):
	for x in p.iterdir():
		if x.is_dir():
			if x.name == ".git":
				change_config(x)
			else: 
				traverse(x)

args = sys.argv

if len(args) == 2:
	traverse(Path(args[1]))
else: 
	print('Please provide a root path as argument')