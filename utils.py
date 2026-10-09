import re


def convert_count_to_number(text):
    match = re.search(r"(\d+(?:,\d{3})*(?:\.\d+)?)\s*([KMB]?)", text)
    return int(float(match[1].replace(",", "")) * {"": 1, "K": 1000, "M": 1000000, "B": 1000000000}[match[2]]) if match else None