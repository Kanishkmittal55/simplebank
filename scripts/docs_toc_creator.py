

import os

def read_summaries_from_result_file(result_txt_path):
    summaries = {}
    with open(result_txt_path, 'r', encoding='utf-8') as file:
        content = file.read()

    for entry in content.strip().split('\n\n'):
        if entry.startswith('Summary for'):
            # Extract the file path and the summary
            path_end_index = entry.find(':')
            summary_start_index = entry.find('\n', path_end_index)
            file_path = entry[len('Summary for'):path_end_index].strip()
            summary = entry[summary_start_index:].strip()

            # Log the file path and summary
            print(f"Reading summary for: {file_path}")
            print(f"Summary: {summary}\n")

            summaries[file_path] = summary

    return summaries



def list_markdown_files(directory):
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith('.md'):
                yield os.path.join(root, file)

def generate_toc(directory, summaries):
    toc_dict = {}
    for full_path in sorted(list_markdown_files(directory)):
        path_parts = full_path.split(os.sep)

        # Use the path without the filename for structuring the TOC
        directory_structure = os.sep.join(path_parts[:-1])
        file_name = path_parts[-1]

        if directory_structure not in toc_dict:
            toc_dict[directory_structure] = []
        toc_dict[directory_structure].append(file_name)

    toc_lines = ['## Table of Contents\n']
    for directory, files in toc_dict.items():
        # Add the directory as a main point
        base_path = os.path.join("ledger", "docs")
        formatted_directory = os.path.join(base_path, os.path.relpath(directory, start="../../docs"))
        toc_lines.append(f"- {formatted_directory}")

        # Add each file within the directory as a sub-point
        for file in files:
            relative_path = os.path.join(directory, file).replace('\\', '/')
            first_part, dots, relative_file_path = relative_path.partition('../../')
            final_paths = dots+relative_file_path
            summary = summaries.get(final_paths, "No summary available.")
            toc_lines.append(f"  - [{file}](https://github.com/dailypay/ledger/tree/main/{relative_file_path}) - {summary}")

        # Add a newline for space after each section
        toc_lines.append('\n')

    return '\n'.join(toc_lines)

def clear_old_toc(readme_path):
    start_marker = '<!-- TOC -->'
    end_marker = '<!-- END TOC -->'
    with open(readme_path, 'r') as file:
        content = file.readlines()

    # Check if both start and end markers are present
    start_index = None
    end_index = None
    for i, line in enumerate(content):
        if start_marker in line:
            start_index = i
        if end_marker in line:
            end_index = i
            break  # No need to continue once we found the end marker

    # If both markers are found, remove the old TOC
    if start_index is not None and end_index is not None:
        content = content[:start_index + 1] + content[end_index:]

    return content

def update_readme_with_toc(readme_path, toc):
    # Clear the old TOC
    content = clear_old_toc(readme_path)

    # Find the position to insert the new TOC
    toc_index = next((i for i, line in enumerate(content) if '<!-- TOC -->' in line), None)
    if toc_index is not None:
        # Insert the new TOC
        content = content[:toc_index + 1] + [toc + '\n'] + content[toc_index + 1:]

    # Write the updated content back to the README
    with open(readme_path, 'w') as file:
        file.writelines(content)

def main():
    # Paths
    docs_path = os.path.join(os.getcwd(), '../..', 'docs')
    print("pwd:", os.getcwd())
    readme_path = os.path.join(os.getcwd(), '../..', 'docs/table_of_contents.md')
    result_txt_path = os.path.join(os.getcwd(),
                                   '../../../../Backend_master_class/bank_api/simplebank/scripts/result.txt')
    print("Path to docs directory:", docs_path)
    print("Path to README.md:", readme_path)
    print("Path to result.txt:", result_txt_path)

    # Read summaries from the result.txt file
    summaries = read_summaries_from_result_file(result_txt_path)

    # Generate TOC with summaries
    toc = generate_toc(docs_path, summaries)

    # Update README.md with the new TOC
    update_readme_with_toc(readme_path, toc)

if __name__ == "__main__":
    main()
