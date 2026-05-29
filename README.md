# RagApp

This is a Rag server written in Go to search through unorganized plain text documents.
I made it by myself in a language I'm not necessarily an expert in, so if you have any suggestions to modify the server, you can open an issue or a pr and I'll review it.

## Quick start

### Requirements

You should have a few things setup already for the server to run correctly. Firstly, have a config file. Compared to the directory where the server is, it should be at `./config/server_config.json`. inside of that JSON file, you should have a few fields:
 - bleve_index_path: the path to the bleve index
 - docs_path: the path to your documents
 - model: the model you want to use in ollama
 - language: the language your documents are in, and you want the model to answer in. If null, defaults to English.
 - last_update: this field is managed by the server, you should set it to null initially and let the server handle it.

If this isn't very clear, an example is provided at [examples/server_config.json](examples/server_config.json)

Additionally, the folders you specify inside of that JSON should be created, especially the bleve index path.

You should also have an ollama server running at http://localhost:11434, you should also choose a model that supports tool calling (otherwise the server won't work), and have it already installed through `ollama pull {model name}`.

### Running the server

In order to start the server, you can get the binary (or exe) in the releases. In order to launch the server with the base address, you can simply launch the binary. If you want to launch the server at a custom address and port, you can run the server with the following parameters:
```
rag-server --address {address} --port {port}
```
The address shouldn't have the `http://` prefix at all, simply put the address like `localhost` or `127.0.0.1`. The server defaults to `localhost:5051`.

### Querying the server

In order to get your documents, the endpoint you want is `/search`. This endpoint is in POST, you want to give a JSON with a `query` field containing a string.
The server will return a JSON document with 2 fields:
  - result: either 'success' or 'failure' depending on the success of the search query (note that this is only for the search query, other errors in the handler will return a 500)
  - content: the result of the query operation. Either the model's findings in case of a success, or the Error message if one was encountered.

Some query and response examples are included over at [examples/query.json](examples/query.json) and [examples/response.json](examples/response.json).

## API specification

Here's the documentation for the different endpoints of the server

### / (GET)

This is the root of the server. It can be used to check if the server is alive. Returns OK if the server is indeed alive.

### /search (POST)

This is the main search endpoint. It's used to search through the index for information based on teh given query.
It receives parameters through a JSON format:

```json
{
  "query": "{your query}"
}
```

It then searches through the database. This can take a while since the model might call the tool a few times.
Once the model is done, the server will return its response in a JSON format:

```json
{
  "result": "success",
  "content": "{model answer}"
}
```

If an error is encountered while the search is going on, the server will still return a result:

```json
{
  "result": "failure",
  "content": "{error message}"
}
```

Note that this response is returned only if the server encounters an error while searching. if it encounters an error while parsing the user query, or while returning the result, it will return a code 500.

### /reset (GET)

This endpoint is used to reset the database of the server. It will use the in-memory configuration in order to get the index path and the document path. It will delete the entire existing Bleve index and recreate a new one with the specified documents.

### /reconfig (GET)

This endpoint is used to resynchronize the in-memory configuration of the server with the on-disk configuration. it will fetch the JSON file and set the configuration to the contents of the file. If you do any modification to the file, make sure to call this endpoint to make sure the server gets updated correctly.

## Dependencies

This project is directly dependent on a few libraries from other projects.

Bleve:
For searching through the documents in full text in order to not have to tokenize the documents, since they're not organized in any particular pattern that would save context.
https://github.com/blevesearch/bleve

Text extraction:
 - go-docx: for extracting text from docx files - https://github.com/fumiama/go-docx
 - pdf: for extracting text from pdf files - https://github.com/ledongthuc/pdf
 - odf: for extracting text from open document format (odf) files - https://pkg.go.dev/sbinet.org/x/odf
