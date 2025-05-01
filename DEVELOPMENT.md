# Development

This project is designed to be developed and built using the provided Dev
Container.

## Building

To build the project, simply run the provided `build.sh` script.

To run:

`./build.sh`

## Debugging

To debug the project, create a `.env` file in the root directory of the workspace. The `.env` file should contain (at a minimum) the following:

```
CHAT_API_KEY=<your_openai_api_key>
```

The `.env` file is ignored, so there is no risk of accidentally committing your API key (unless you are really trying...)

After the `.env` file is created, you can use the provided launch configurations for debugging by pressing `F5` or using the "Run and Debug" tab in VSCode.