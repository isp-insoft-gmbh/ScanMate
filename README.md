# ScanMate

## Project stages

1. Simple user identification + drink identification with minimum input and feedback
2. User input via touch and visual feedback plus optional audio feedback
3. Locking / Unlocking a fridge

## Hardware required

* Stage 1
  * Ras-PI (Any generation)
    * SD-Card (4GB+ is recommended)
    * Charger + Cable
  * Camera module
  * 3 different colored LEDs (YELLOW, GREEN and RED)
  * Hardware button
* Stage 2
  * Touchscreen
* Stage 3
  * Locking mechanism
    * Solenoid
    * "Key hole" (maybe)
  * Case

## Tools required on Ras-PI

* raspistill - Take image (Installed by default)
* zbarimg - Scan barcode (Installable via sudo apt install zbar-tools)

## How to setup Ras-PI

1. Install Raspbian Lite image via Raspbian Imager
2. Create /boot/ssh to activate sshd (ssh server)
3. login via `ssh pi@IP`
4. Run `sudo raspi-config` in ssh session and select the following
  * Interface
  * Enable camera
  * Finish
5. `sudo reboot` to make sure `raspi-config` properly applied everything

## Testing in a development environment

The codebase allows switching into different operation modes via build-tags.
Those can be supplied by running the desired go command with the `-tags` flag.
Possible values for tag are `dummy`, `test`.
The `test` tag is required for unit tests, while the `dummy` tag is required
for manual testing.

Examples:

```
go run . # run normally on a raspberry pi
go run -tags dummy .
go test -tags test ./...
```

Attemping to run the unit-tests on the raspberry pi would require user
interaction and is currently not supported.

Running with `dummy` enabled will ask for user input, whenever the program
would interact with any raspberry pi hardware, such as GPIO ports or a camera.
During input querying, the application can't be killed via Ctrl-C. Instead
Ctrl-D can be used. The SIGTERM and SIGINT signals will still be caught if
sent any other way, as these work via process and not the terminal input.

These values can also be set in VS Code, in order to let the go language server
know which tags to use. Simply edit the `go.buildTags` value in
`.vscode/settings.json`. These changes affect the compiler, the linter and
other go tools, even documentation popups.

## User Stories

### Purchasing a drink

1. Initialize purchasing process (via button press for example)
2. Parse Ident-Card to identify user
3. Parse barcode of drink to know price
4. Send both to MeMate in order to confirm the purchase

### Registering a user

1. Manually add `Ident-Card` code to userprofile in MeMate
2. Print `Ident-Card` barcode
3. Glue barcode onto `card` (cardboard)

## MeMate requirements

* Save bardcode for each drink
* Save `Ident-Card` barcode for each user
* Allow purchasing a drink via `Ident-Card` barcode and drink barcode
* Optional
  * Generate `Ident-Card` barcode image (could be external tool as well)

## Stage one purchase workflow

The idea with this workflow is to have different signals in different stages
and allow the user to have some failures inbetween. This makes sure that
user doesn't have to rerun stage one when failing stage two. This concept
does not yet take acconut for any failures during server communication.

0. Wait for button press
1. Wait 2s
2. Set LED BLUE for "waiting for identification barcode" indication
3. Wait 1.5s
4. Take picture
5. Set LED NONE
6. Parse image for barcode
7. If fail { if failcount == 10 { GO TO 1. } else { BLINK RED 500ms; GO TO 2. } }
8. Set LED MAGENTA for "waiting for drink barcode" indication
9. Wait 1.5s
10. Take picture
11. Set LED NONE
12. Parse image for barcode
13. If fail { if failcount == 10 { GO TO 1. } else { BLINK RED 500ms; GO TO 8. } }
14. Set LED GREEN
15. Wait 3s
16. Set LED NONE
17. GO TO 0.
