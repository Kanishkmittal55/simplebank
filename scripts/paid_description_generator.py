import os
import time
import openai
import math
from openai import OpenAI
import tiktoken
import constants

client = OpenAI()
client.api_key = constants.APIKEY

# Constants
MAX_TOKENS = 4096
MAX_TOTAL_TOKENS = 40000

def get_token_count(text):
    """Get the token count for a given text using OpenAI's API"""
    try:
        encoding = tiktoken.encoding_for_model("gpt-3.5-turbo")
        encoded_text = encoding.encode(text)
        return len(encoded_text)
    except openai.OpenAIError as e:
        print(f"Error in getting token count: {e}")
        return None

def remove_less_important_words(text):
    """Remove less important words to reduce token count"""
    # List of words to be removed (common stopwords, etc.)
    stopwords = set(["a", "an", "the", "and", "or", "but", "is", "are", "was", "were", "be", "been", "being", "have", "has", "had", "do", "does", "did", "will", "would", "shall", "should", "can", "could", "may", "might", "must", "ought"])
    words = text.split()
    # Remove stopwords and reconstruct the text
    filtered_words = [word for word in words if word.lower() not in stopwords]
    return ' '.join(filtered_words)

def is_summary_present(md_file_path, result_txt_path):
    """Check if the summary of the file already exists in result.txt"""
    if not os.path.exists(result_txt_path):
        return False
    search_string = f"Summary for {md_file_path}:"
    with open(result_txt_path, 'r', encoding='utf-8') as result_file:
        result_file_content = result_file.read()
        if search_string in result_file_content:
            return True
    return False

def split_text(text, tokenCount):
    words = text.split()
    n = len(words)
    parts = math.ceil(tokenCount / MAX_TOKENS)
    part_size = math.ceil(n / parts)
    print(f"Total words: {n}, Parts: {parts}, Part size: {part_size}")

    split_texts = [' '.join(words[i:i + part_size]) for i in range(0, n, part_size)]
    for i, part in enumerate(split_texts, 1):
        print(f"Part {i} length: {len(part.split())} words")

    return split_texts

def process_in_parts(text, token_count):
    summaries = []

    split_texts = split_text(text, token_count)
    for index, part in enumerate(split_texts):
        print(f"Processing part {index + 1}/{len(split_texts)}")
        system_message = {"role": "system", "content": "You are an assistant who provides 1 line summary for documents."}
        user_message = {"role": "user", "content": part}
        try:
            completion = client.chat.completions.create(
                model="gpt-3.5-turbo",
                messages=[system_message, user_message]
            )
            summary = completion.choices[0].message.content
            summaries.append(summary)
            print(f"Received summary for part {index + 1}: {len(summary.split())} words")

        except openai.BadRequestError as e:
            print(f"Error in processing part {index + 1}: {e}")
            return None

    return " ".join(summaries)

def process_md_file(md_file_path, file_name, result_txt_path):
    """Process the Markdown file and summarize its content"""
    print(f"Processing file: {md_file_path}")  # Log the file being processed

    # Read the content of the Markdown file
    with open(md_file_path, 'r', encoding='utf-8') as md_file:
        md_content = md_file.read()

    # Remove less important words
    md_content_filtered = remove_less_important_words(md_content)
    # token_count = len(md_content_filtered.split())

    # Use OpenAI API to get the accurate token count
    token_count = get_token_count(md_content_filtered)
    if token_count is None:
        print("Failed to get token count from OpenAI.")
        return

    if token_count > MAX_TOTAL_TOKENS:
        print(f"Skipping {file_name}, file is too large (exceeds 40,000 tokens).")
        return

    # Define the system and user messages
    system_message = {"role": "system", "content": "You are an assistant who provides 1 line summary for documents."}
    user_message = {"role": "user", "content": md_content_filtered}

    try:
        completion = client.chat.completions.create(
            model="gpt-3.5-turbo-16k",
            messages=[system_message, user_message]
        )
        summary = completion.choices[0].message.content
        print(f"Summary result for: {file_name}: {summary}\n\n")  # Log the result

    except openai.BadRequestError as e:
        with open(result_txt_path, 'a', encoding='utf-8') as result_file:
            result_file.write(f"Summary for: {md_file_path}:\nIt is still to do because text was too long to analyze\n\n")
        print(f"Un-rounded_parts_needed: {token_count / MAX_TOKENS} ")
        print(f"token count : {token_count} & max_tokens : {MAX_TOKENS} ")
        summary = process_in_parts(md_content_filtered, token_count)
        if summary is None:
            return

    with open(result_txt_path, 'a', encoding='utf-8') as result_file:
        result_file.write(f"Summary for: {md_file_path}:\n{summary}\n\n")

def main():
    # Paths
    docs_path = "../../docs"
    list_txt_path = "list.txt"
    result_txt_path = "result.txt"

    # Read the list of files from list.txt
    with open(list_txt_path, 'r', encoding='utf-8') as list_file:
        md_files = list_file.readlines()

    # Iterate over the sorted list of Markdown files
    for md_file_path in md_files:
        md_file_path = md_file_path.strip()
        file_name = os.path.basename(md_file_path)

        # Skip if summary already exists
        if is_summary_present(md_file_path, result_txt_path):
            print(f"Skipping {md_file_path}, summary already exists.")
            continue

        # Process the Markdown file
        process_md_file(md_file_path, file_name, result_txt_path)
        time.sleep(50)  # Delay to avoid hitting rate limits

    print("Summaries have been written to result.txt")

if __name__ == "__main__":
    main()