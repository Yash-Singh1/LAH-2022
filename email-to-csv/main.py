import pandas as pd
import numpy as np
import os
import email.parser
import email.policy
from bs4 import BeautifulSoup
import csv

ham_filenames = [name for name in sorted(os.listdir('./hamnspam/ham')) if len(name) > 20]
spam_filenames = [name for name in sorted(os.listdir('./hamnspam/spam')) if len(name) > 20]

def load_email(is_spam, filename):
    directory = "./hamnspam/spam" if is_spam else "./hamnspam/ham"
    with open(os.path.join(directory, filename), "rb") as f:
        return email.parser.BytesParser(policy=email.policy.default).parse(f)
    
ham_emails = [load_email(is_spam=False, filename=name) for name in ham_filenames]
spam_emails = [load_email(is_spam=True, filename=name) for name in spam_filenames]

def html_to_plain(email):
    try:
        soup = BeautifulSoup(email.get_content(), 'html.parser')
        return soup.text.replace('\n\n','')
    except:
        return None

with open('main.csv', 'w') as f:
    # create the csv writer
    writer = csv.writer(f)

    for hamEmail in ham_emails:
        if hamEmail.get_content_type() == 'text/plain' and hamEmail.get_content_subtype() == 'plain':
            writer.writerow([hamEmail.get_content(), 0])
        elif hamEmail.get_content_type() == 'text/html' and hamEmail.get_content_subtype() == 'html':
            plain = html_to_plain(hamEmail)
            if not(plain):
              continue
            writer.writerow([plain, 0])

    for hamEmail in spam_emails:
        if hamEmail.get_content_charset() != 'us-ascii':
          continue
        if hamEmail.get_content_type() == 'text/plain' and hamEmail.get_content_subtype() == 'plain':
            writer.writerow([hamEmail.get_content(), 1])
        elif hamEmail.get_content_type() == 'text/html' and hamEmail.get_content_subtype() == 'html':
            plain = html_to_plain(hamEmail)
            if not(plain):
              continue
            writer.writerow([plain, 1])
