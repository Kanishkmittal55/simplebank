import os
import fnmatch
from reportlab.lib.pagesizes import letter
from reportlab.pdfgen import canvas
from reportlab.lib.utils import ImageReader

# Function to find files with matching patterns
def find_files(directory, file_patterns):
    for root, dirs, files in os.walk(directory):
        for basename in files:
            for pattern in file_patterns:
                if fnmatch.fnmatch(basename, pattern):
                    filename = os.path.join(root, basename)
                    yield filename

# Function to create the PDF document
def create_pdf(directory, file_types, output_file):
    c = canvas.Canvas(output_file, pagesize=letter)
    width, height = letter
    y = height - 40  # Starting Y position

    for file in find_files(directory, file_types):
        if y < 100:  # New page if not enough space
            c.showPage()
            y = height - 40

        # Adding file path and name
        c.drawString(40, y, f"File Path: {os.path.dirname(file)}")
        y -= 20
        c.drawString(40, y, f"File Name: {os.path.basename(file)}")
        y -= 20

        # Check if file is an image or a text file
        if file.lower().endswith(('.png', '.jpg', '.jpeg', '.gif', '.bmp')):
            try:
                image = ImageReader(file)
                c.drawImage(image, 40, y - 100, width=400, preserveAspectRatio=True)
                y -= 120  # Adjust space for the next file
            except Exception as e:
                c.drawString(40, y, "Image detected but cannot be added: " + str(e))
                y -= 20
        else:
            # Adding file contents
            try:
                with open(file, 'r') as f:
                    contents = f.read(500)  # Read first 500 characters
                    c.drawString(40, y, "Contents: " + contents[:500])
                    y -= 40  # Adjust space for the next file
            except Exception as e:
                c.drawString(40, y, "Error reading file: " + str(e))
                y -= 20

    c.save()

# Specify the directory, file types and output file
directory = '.'  # Current directory
file_types = ['*.go', '*.yaml', '*.js', '*.png', '*.jpg', '*.jpeg', '*.gif', '*.bmp']  # Add more file types if needed
output_file = 'output.pdf'

create_pdf(directory, file_types, output_file)
