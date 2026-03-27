# sh_test

Convert the line-based input format in `example.txt` into XML.

## Download

Built binaries are published in the repository's release downloads:

- `sh_test-linux-amd64`
- `sh_test-darwin-amd64`
- `sh_test-windows-amd64.exe`

Pick the file that matches your platform, download it, and run it from a terminal.

## Command Usage

```text
sh_test <input-file> [-o <output-file>]
```

### Arguments

- `<input-file>`: required path to the source text file.
- `-o <output-file>` or `--output <output-file>`: optional path to write the XML output.

If `-o` is omitted, the XML is written to standard output.

## Examples

Print XML to the terminal:

```powershell
.\sh_test-windows-amd64.exe .\example.txt
```

```bash
./sh_test-linux-amd64 ./example.txt
```

Write XML to a file:

```powershell
.\sh_test-windows-amd64.exe .\example.txt -o .\output.xml
```

```bash
./sh_test-linux-amd64 ./example.txt --output ./output.xml
```

## Input Format

The input file is a pipe-delimited text file where the first field selects the record type:

- `P|<first-name>|<last-name>`: starts a new person.
- `T|<mobile>|<landline>`: adds phone data to the current person or current family member.
- `A|<street>|<city>|[postal-code]`: adds address data to the current person or current family member.
- `F|<name>|<born>`: adds a family member under the current person.

Order matters. `T`, `A`, and `F` apply to the most recently opened person, and `T` and `A` apply to the most recently opened family member when one exists.

Example input:

```text
P|Carl Gustaf|Bernadotte
T|0768-101801|08-101801
A|Drottningholms slott|Stockholm|10001
F|Victoria|1977
A|Haga Slott|Stockholm|10002
F|Carl Philip|1979
T|0768-101802|08-101802
P|Barack|Obama
A|1600 Pennsylvania Avenue|Washington, D.C
```

## Output

The command generates XML with this top-level structure:

```xml
<people>
  <person>
    ...
  </person>
</people>
```
