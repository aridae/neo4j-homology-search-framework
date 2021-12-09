import httplib2
from bs4 import BeautifulSoup, SoupStrainer
import argparse
import subprocess
import logging

LOGFILE = "./wget-logs.txt"

parser = argparse.ArgumentParser(description='Scrap links from 1001genomes page.')
parser.add_argument('-dir', '--directory', help="directory to save files")
parser.add_argument('-url', '--url', help="page url to scrap from")
parser.add_argument('-prefix', '--prefix', help="filename prefix to scrap all files from the page")
parser.add_argument('-ext', '--extension', help="filename extension to scrap all files from the page")
args = parser.parse_args()

http = httplib2.Http()
status, response = http.request(args.url)

fastaurls = list()
for link in BeautifulSoup(response, parse_only=SoupStrainer('a')):
    if link.has_attr('href') and link.text.startswith(args.prefix) and link.text.endswith(args.extension):
        fastaname = link['href']
        fastaurls.append(args.url + link['href']) 
        retcode = subprocess.call(["wget", "-P", args.directory, "-a", LOGFILE, args.url + link['href']])
        if retcode == 0:
            logging.log(logging.DEBUG, f"[exited with {retcode}] sucessfully downloaded {fastaname} to {args.directory}")
        else:
            logging.log(logging.ERROR, f"[exited with {retcode}] **ERROR** failed to download {fastaname} to {args.directory}")