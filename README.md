# ScanMate

## Tools required on Ras-PI

* raspistill - Take image (Installed by default)
* zbarimg - Scan barcode (Installable via sudo apt install zbar-tools)

## How to setup Ras-PI

1. Raspbian Lite image via Raspbian Imager installieren
2. /boot/ssh erstellen um sshd zu aktivieren
3. Einloggen via ssh pi@IP
4. sudo raspi-config
  * Interface
  * Enable camera
  * Finish
5. sudo reboot

## User registration

1. Manuelle Eintragung in Memate
2. Drucken des Barcodes
3. Auf Pappe kleben

## User Story

1. Parse Ident-Card
2. Parse barcode
3. Send both to memate in one call

## Idee für Screenless Workflow:

### Setup

* Button
* 3 LEDs (YELLOW, GREEN, RED)

### Workflow

* Press Button
* YELLOW LED lights up (READY state)
* Camera starts taking pictures for USER_IDENT and DRINK_IDENT
* RED / GREEN LED lights up, depending on outcome (ERROR or SUCCESS)

## TODOs MeMate

* Find java barcode parser

