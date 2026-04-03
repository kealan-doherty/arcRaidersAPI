import json 
from pathlib import Path 
import csv 
import os 
import glob

"""
JSON files were grabbed from the Open Source arcraiders-data repository. 

Link: https://github.com/RaidTheory/arcraiders-data

this sciprt reads a JSON folder of items and converts them to a CSV file. 
It filters the data to include only specific keys: 
"type", "rarity", "weightKg", "value", and "isWeapon". T
he resulting CSV file will have these keys as column headers 
and their corresponding values in the rows.

"""

repo_root = Path(__file__).resolve().parent.parent
folder_path = glob.glob(str(repo_root / "items" / "*.json"))

selected_keys = ["id", "type", "rarity", "weightKg", "value", "isWeapon"]

filtered_data = []

"""
parse the JSON files for the necessary keys and values, to make writing to a CSV file easier.
"""

for file_path in folder_path:
    with open(file_path, "r") as f:
        raw_data = json.load(f)
    filtered_data.append({key: raw_data.get(key) for key in selected_keys})

"""
write the parsed JSON data to a CSV file with ready to be read into the database. 
"""

output_path = Path(__file__).resolve().parent / "data.csv"

with open (output_path, "w", newline="") as csvfile:
    writer = csv.DictWriter(csvfile, fieldnames=selected_keys)
    writer.writeheader()
    writer.writerows(filtered_data)
