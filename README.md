# sthlmlib

A tool to dump Stockholm library loans as an iCal file. Uses the same graphql
API as the web site. Expect breakage from time to time. Use at your own risk.
Vibe coded with Gemini 2.5. This project is not associated with the Stockholm
libraries.

## Usage

```bash
go install github.com/tomyl/sthlmlib@latest
sthlmlib -card-number <social security number or card number> -pin <pid> -ical -group
```
