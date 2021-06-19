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

* Press Button to initiaze workflow
* YELLOW LED lights up (READY state)
* Camera starts taking pictures for USER_IDENT and DRINK_IDENT (in this order)
* RED / GREEN LED lights up, depending on outcome (ERROR or SUCCESS)
