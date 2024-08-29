# GoAsk - AI-powered Linux Command Assistant

GoAsk is a command-line tool that uses AI to suggest and execute Linux commands based on natural language descriptions.

## Installation

1. Ensure you have Go installed on your system. If not, download and install it from [golang.org](https://golang.org/).

2. Clone this repository or download the `goask.go` file.

3. Install the required dependencies:
   ```
   go get github.com/joho/godotenv
   ```

4. Compile the program:
   ```
   go build -o goask goask.go
   ```

5. Move the compiled binary to a directory in your system's PATH. For example:
   ```
   sudo mv goask /usr/local/bin/
   ```

## Configuration

1. Create a `.goask.env` file in your home directory:
   ```
   touch ~/.goask.env
   ```

2. Open the `.goask.env` file in a text editor and add your Gemini API key and endpoint:
   ```
   GEMINI_API_KEY=your_actual_api_key_here
   GEMINI_API_ENDPOINT=https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent
   ```

   Replace `your_actual_api_key_here` with your real Gemini API key.

3. Save and close the file.

## Usage

Once installed and configured, you can use GoAsk from anywhere in your terminal:

```
goask how to do something in linux
```

For example:
```
goask how to list all files including hidden ones
```

The AI will suggest a command, and you'll be prompted whether you want to execute it.

## Safety Note

Always review the suggested command before executing it, especially for commands that might modify your system.

## Troubleshooting

If you encounter any issues:

1. Ensure your `.goask.env` file is correctly set up in your home directory.
2. Check that you have the necessary permissions to execute the `goask` command.
3. Verify that your Gemini API key is valid and has the necessary permissions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Specify your license here]
