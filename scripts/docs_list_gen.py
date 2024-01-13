import os

def find_md_files(directory):
    md_files = []
    # Traverse the directory and its subdirectories
    for root, dirs, files in os.walk(directory):
        for file in files:
            # Check if the file is a Markdown file
            if file.endswith('.md'):
                # Construct the full file path
                file_path = os.path.join(root, file)
                # Add the file path to the list
                md_files.append(file_path)
    return md_files

def write_to_file(file_list, output_file):
    # Sort the list of files
    file_list.sort()
    # Write the sorted list to the output file
    with open(output_file, 'w') as f:
        for file in file_list:
            f.write(f"{file}\n")

def main():
    directory = '../../docs'  # Replace with your directory path or if you change the location of script and output remember to update this path
    output_file = 'list.txt'
    md_files = find_md_files(directory)
    write_to_file(md_files, output_file)
    print(f"Markdown files list written to {output_file}")

if __name__ == "__main__":
    main()
