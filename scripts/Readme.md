# Docs Generation:

This directory dailypay/ledger/scripts/docs_generator consists of several python files that have the capability to generate docs for any type of code or documentation files one may provide to it.

# Steps to Follow:

1. Firstly one needs to give the program, <u>relative path</u> of the folder/directory where all the TO BE ANALYZED documents or code files are located, in this case inside the `../../docs` folder. This can be auto-generated using the `docs_list_gen.py` file which will generate a list of relative paths for all the files located in the provided directory. See `list.txt` shown below for further clarity.

   ![Where to set the directory read path](../../../../Backend_master_class/bank_api/simplebank/scripts/images/img.png)

2. So simply run `python3 docs_list_gen.py` from inside the ledger/scripts/docs_generator directory, after this one will observe `list.txt` will be populated with the relevant data.

   ![What list.txt looks like](../../../../Backend_master_class/bank_api/simplebank/scripts/images/img_1.png)

3. Then simply run `python3 paid_description_generator` , which will automatically start generating docs for the directory specified above, and the resulting analysis or summary being generated will automatically be populated in the `result.txt` file.

<u>Points to note</u> - Each request/file takes approximately 50 seconds to send the request/text to openai server and get a response/summary from it , due to rate limit issues this latency should not be reduced , but one can upgrade to a better model (gpt-4 etc.) but then costs will also increase, so please be patient while it generates the required documentation for the provided files.

![img.png](../../../../Backend_master_class/bank_api/simplebank/scripts/images/img_2.png)

4. One can read the necessary documentation from either `result.txt` file or they can use the output like in this case to AUTO GENERATE A TABLE OF CONTENTS FOR ALL THE DOCS INSIDE THE LEDGER REPO, using the `toc_generator.py`, which in this case will use the `result.txt` to get a 1 line description of every document inside the ledger/docs folder.

5. Lastly, to update the Table of contents -

- Just rerun `python3 docs_list_gen.py` file - which update the list.txt.
- Then run `python3 paid_description_generator.py` - which will skip generating summaries for the ones already present inside the result.txt and will generate summaries for the new files added to the selected directory ( like ../../docs ).
- Then simply run `python3 toc_generator.py` file - To update the Table of contents and generate summaries for it.
