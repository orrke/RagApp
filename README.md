# RagApp

This is a Rag server written in Go to search through unorganized plain text documents.
I made it by myself in a language I'm not necessarily an expert in, so if you have any suggestions to modify the server, you can open an issue or a pr and I'll review it.

## Quick start

### Requirements

The one main requirement is to have an ollama server running in the background, at the address http://127.0.0.1:11434. If you click that link on your computer, it should tell you "Ollama is running".

### First Start

If you're running the server for the first time, chances are you didn't set up the config file by yourself. luckily, you don't need to immediately. You can simply run the server, and it will create a default configuration file. However, if you do, you need to specify a path for where the server can find your documents. You can simply do that by specifying the `--docs` flag when running the server. For example:
```shell
rag-server --docs "/path/to/your/docs"
```
Once you ran that command once, you shouldn't have to worry about it again, it's saved to a config file and loaded every time you run the server, meaning you'll always be using this path. You can still replace it by using that exact same parameter once more.

If you did not set that `--docs` flag on the first run, it should exit on startup.
That's expected behavior, because it doesn't know where to look for your files, and it won't start indexing your entire system for you to search through it. 
(Unless you want it to of course.)

P.S. this works too if you find it more intuitive: `--docs="/path/to/your/docs"`. If you're on Windows, no need to put a double backslash, it takes care of that for you.

You can also set the model to use with `--model "model-name:version`, if you don't set one yourself, it will default to gemma4:latest, so you'll need to have it pulled.
I recommend you choose a model more suited. This one has 10GB of memory, so for everybody with an 8GB card, you might want to go for a smaller model. Either the slightly smaller gemma4:e2B, or some qwen 3.5 models; they should work great too.

### Running the server

In order to start the server, you can get the binary (or exe) in the releases. In order to launch the server with the base address, you can simply launch the binary. If you want to launch the server at a custom address and port, you can run the server with the following parameters:
```shell
rag-server --address {address} --port {port}
```
The address shouldn't have the `http://` prefix at all, simply put the address like `localhost` or `127.0.0.1`. The server defaults to `0.0.0.0:5051`.

### Querying the server

In order to get your documents, the endpoint you want is `/search`. This endpoint is in POST, you want to give a JSON with a `query` field containing a string.
The server will return a JSON document with 2 fields:
  - result: either 'success' or 'failure' depending on the success of the search query (note that this is only for the search query, other errors in the handler will return a 500)
  - content: the result of the query operation. Either the model's findings in case of a success, or the Error message if one was encountered.

Some query and response examples are included over at [examples/query.json](examples/query.json) and [examples/response.json](examples/response.json).

## Web UI

### Running the UI

If your server is running, it means your web UI is running. the server doesn't really have any online dependencies, it was one of my requirements. 
However, it does use the [Marked Js Library](https://github.com/markedjs/marked), to whom all credits go for doing the actual Markdown processing, the only thing I did was to basically copy their file over [here](./web/assets/scripts/marked.js). 
I don't really have any particular Js skills, so thanks to them for making this library for displaying my Md.

### Showing the UI

To get the actual UI, you can simply open your favorite web browser, and then go to the IP you set for the server at startup. If you didn't set any, then you can simply try to go here with the server running: http://localhost:5051.
If all goes well, you should land on an HTML page prompting you to query the database. The page should be localized according to your browser settings, defaulting to English if your locale isn't supported.

While talking about the locales, I only made an English one, as well as a French one. If you would like your locale to be added, please do reach out to me. 
The only ones supported will probably be the ones Bleve supports (since it's easier for me that way), but if you want to try to implement something yourself, I'd be happy to review any PR.

## Troubleshooting

### /search endpoint fails silently

This is most likely due to you not having pulled the model. If you didn't specify any model, then go to your terminal and run this:
```shell
ollama pull gemma4:latest
```
If you did specify a model, then pull it from ollama to be able to use it.

### The model doesn't find the info

This can be due to several things, the most likely one being your file is not supported. Right now, the extensions that are supported are the following:
`.txt .csv .bat .doc .docx .xlsx .pdf .odt`.
Another reason might be that your data is simply too large. I didn't necessarily run a lot of tests myself, but I did notice that after too many tokens, the model just isn't able to find the information at all anymore. In that case, I would recommend you try another model with a larger context window. If that still does not work, then you should probably try another RAG application, or split your documents. 

### The cmd pops up on Windows

This is an issue when building normally in go. 
In order to not have this terminal pop up, you can try to use the ones in the release section (which should be built correctly), or build it yourself, by adding these flags to your build command: 
`-ldflags -H=windowsgui`. 
This flags the app as GUI app (and not CLI) to windows, which should prevent it from spawning a cmd prompt.

### The model doesn't use the GPU

That's not a problem with my app, it's something with Ollama. If you have an Nvidia or AMD GPU, it should work out of the box with ollama. If you have an Intel GPU, there are not many options left. I recommend using [ipex-LLM](https://github.com/intel/ipex-llm) even though it's deprecated, because it still works fine for most models. For more advanced Ollama troubleshooting, I'll redirect you to their documentation: https://docs.ollama.com/gpu.

### No Auth ?
No, there's no authentication. I built this with the intention of it being for private use. If you really do want to have a form of authentication for the user, then I can try something, but for now I don't plan on adding it to the server at all.

### How can I debug myself?

You also have some server logs, they should be located at `os.UserConfigDir()/RagApp/logs`.  
If you're on Windows, that's `C:\Users\{Your username}\AppData\Roaming\RagApp\logs`.
In PS, you should be able to go there using this command:
```shell
cd $env:AppData\RagApp\logs
```
If you're on Linux, then you should look at `$XDG_CONFIG_HOME/RagApp/logs`. 
In any Linux shell, you should simply be able to go there with these commands: 
```bash 
cd $XDG_CONFIG_HOME/RagApp/logs
```

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

### /search/raw (POST)

This endpoint is for getting a direct response from the bleve index. I used it for debugging, but you might find a better use for it. It simply searches through an index with your query as the search parameters. It doesn't rely on Ollama at all, meaning you don't even need it installed to be able to user this endpoint. This can be useful for lower end hardware, if you don't have the necessary system memory to be able to run a model comfortably.

This endpoint works the exact same as /search. You send it the same request, with this form:
```json
{
  "query": "{your query}"
}
```

And it will return a response in this form:
```json
{
  "result": "{result}",
  "content": "{message}"
}
```

If the query is successful, it will return a string made up of the raw documents as JSON in a string, they should be separated by a comma `,`.
The document should look something like that:
```json
{
  "Title": "{Document title}",
  "Content": "{Entire file content as string}",
  "Path": "{Full doc path}"
}
```
The document gets passed as this exact JSON string to the model after the tool call, meaning the only difference between this and the /search endpoint is that the model might do custom, more directed queries to bleve. There is also of course always a risk of hallucination, which is why you shouldn't trust the answer entirely, and maybe verify with this endpoint.

### /reset (GET)

This endpoint is used to reset the database of the server. It will use the in-memory configuration in order to get the index path and the document path. It will delete the entire existing Bleve index and recreate a new one with the specified documents.

### /reconfig (GET)

This endpoint is used to resynchronize the in-memory configuration of the server with the on-disk configuration. it will fetch the JSON file and set the configuration to the contents of the file. If you do any modification to the file, make sure to call this endpoint to make sure the server gets updated correctly.

## Configuration

You have some configuration options with this server, it can be configured with the JSON file located at `os.ConfigDir()/RagApp/server_config.json`. In there, you can specify a custom location for the Bleve index, as well as another directory for the documents, a different model to use in the backend, or another language. please note that, if you change the language, you will have to reload the index (with the /reset endpoint). You can have an example of this config file at [examples/server_config.json](examples/server_config.json). 

In this file, you have a `last_update` field, it is entirely handled by the server. It's meant to represent the last time the database was updated, and serves to not have to re-index every single file every time we start the server again, which can be annoying if you have a lot of documents.

## Build instructions

It's pretty straight forward, you can either do it yourself and just run `go mod tidy & go build -o rag-server ./cmd/server` (replace `&` with `;` in PowerShell). As an extra tip, if you're on Windows, you can add this: `-ldflags -H=windowsgui` so that the cmd doesn't appear when you launch the app.

You can also use `make {your os}` to build the app for your operating system directly. you can choose between `windows`, `linux` and `macos`, the produced file will go in `dist/rag-server-{os}-{arch}(.exe)`. It will build for arm64 if you're targeting macOS.

If you want to be able to use the Web UI, then you will also need to have the `web/assets` folder available (with the web folder too). 
I didn't want to bundle it in an embed fs, that way it's a bit easier to do modifications on the ui. 
So, in the end, you should end up with this structure:
```
├── rag-server-{os}-{arch}(.exe)
└── web
    └── assets
            ├── locales
            │       └── index
            │           ├── en.json
            │           └── ...
            ├── scripts
            │       ├── main.js
            │       └── ...
            └── style
                    ├── index.css
                    └── ...
```

## Dependencies

This project is directly dependent on a few libraries from other projects.

Bleve:
For searching through the documents in full text in order to not have to tokenize the documents, since they're not organized in any particular pattern that would save context.
https://github.com/blevesearch/bleve

Text extraction:
 - go-docx: for extracting text from docx files - https://github.com/fumiama/go-docx
 - pdf: for extracting text from PDF files - https://github.com/ledongthuc/pdf
 - odf: for extracting text from open document format (odf) files - https://pkg.go.dev/sbinet.org/x/odf

Marked.js:
For displaying the Markdown correctly as HTML tags I did copy the file, but it's not mine at all, all credits go to the people behind Marked.js, their doc is here: https://marked.js.org.
