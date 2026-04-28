# DCRUN - DON'T CARE, RUN

dcrun is a cli to run projects in dev mode in the most popular programming languages, without caring!

# Installing

## Direct bin

```bash
go install github.com/davitostes/dcrun@latest
```

## From source

```bash
git clone https://github.com/davitostes/dcrun
cd dcrun
go build -o dcrun .
mv dcrun ~/.local/bin/
```

Now you have dcrun working!

## How it works

Just enter your project directory and type: 

```bash
dcrun
```

And your project will start to run in dev mode.

## Supported languages and frameworks

- Go
- NodeJS (npm)
- Bun
- Yarn
- pnpm
- Deno
- Rust
- Python (uv, poetry)
- Django
- Dotnet
- Lua
- Make (`dev` target fallback)

## Contributing

Feel free to contribute with PRs or Issues

## License

[MIT LICENSE](LICENSE)
