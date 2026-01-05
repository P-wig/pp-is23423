#!/bin/python3

import math
import os
import random
import re
import sys
import requests
import json



#
# Complete the 'maximumTransfer' function below.
#
# The function is expected to return a STRING_ARRAY.
# The function accepts following parameters:
#  1. STRING name
#  2. STRING city
#
# Base URL: https://jsonmock.hackerrank.com/api/transactions
#
#

def maximumTransfer(name, city):
    # Write your code here
    Base_URL = "https://jsonmock.hackerrank.com/api/transactions"
    credit_amounts = []
    debit_amounts = []
    page = 1
    res = requests.get( f"{Base_URL}?page={page}" )
    total_pages = res.json()["total_pages"]
    
    while page <= total_pages:
        res = requests.get( f"{Base_URL}?page={page}" )
        data = res.json()
        if not data.get("data"):
            print("could not GET response")
            break
            
        for obj in data["data"]:
            if ( obj.get("userName") == name ) and ( obj.get("location", {}).get("city") == city ):
                # Clean the amount string - remove $ and commas
                amount_str = obj.get("amount", "0")
                amount_clean = amount_str.replace("$", "").replace(",", "")
                amount = float(amount_clean)
                txn_type = obj.get("txnType")
                
                if txn_type == "credit":
                    credit_amounts.append(amount)
                elif txn_type == "debit":
                    debit_amounts.append(amount)
        page += 1
    
    result = []
    if credit_amounts:
        result.append(f"{max(credit_amounts):.2f}")
    if debit_amounts:
        result.append(f"{max(debit_amounts):.2f}")
    
    return result
    
    

if __name__ == '__main__':
    result = maximumTransfer("John Oliver", "Ripley")
    for amount in result:
        print(amount)