# Credit: https://www.geeksforgeeks.org/check-if-email-address-valid-or-not-in-python/
# Python program to validate an Email 

# import re module 

# re module provides support 
# for regular expressions 
import re

# Make a regular expression 
# for validating an Email 
regex = '^[a-z0-9]+[\._]?[a-z0-9]+[@]\w+[.]\w+$'


# for custom mails use: '^[a-z0-9]+[\._]?[a-z0-9]+[@]\w+[.]\w+$'

# Define a function for 
# for validating an Email 
def check(email):
    # pass the regular expression 
    # and the string in search() method 
    if (re.search(regex, email)):
        return True

    else:
        return False
