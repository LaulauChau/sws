# Sowesign Code Generator

## Requirements

- Go 1.23.5 or higher
- [Templ CLI](https://templ.guide/quick-start/installation)
- [Tailwind CSS CLI](https://tailwindcss.com/blog/standalone-cli)

## Configuration

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit the `.env` file with your Sowesign credentials:
```env
SOWESIGN_CODE_ETABLISSEMENT=your_code_etablissement
SOWESIGN_IDENTIFIANT=your_identifiant
SOWESIGN_PIN=your_pin
```

Make sure to keep your `.env` file secure and never commit it to version control.

## Usage

Run the application:
```bash
make run
```

Go to [http://localhost:8080](http://localhost:8080) to see the application.

## License

[MIT License](LICENSE)
