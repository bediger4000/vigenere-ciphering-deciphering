#!/usr/bin/env php
<?php

# PHP implementation of Vigenere encoding
# Usage: vig.php 'keystring' cleartextfilename
# Ciphertext written to stdout.

$cleartext = file_get_contents($argv[2]);

$keystring = $argv[1]
$keylength = strlen($keystring);
$key = array($keylength);

for ($i = 0; $i < $keylength; ++$i) {
	$key[$i] = ord($keystring[$i]);
}

$ciphertext = '';

for ($i = 0; $i < strlen($cleartext); ++$i) {
	$offset = $key[$i%$keylength];
	$ciphertext .= chr((ord($cleartext[$i]) + $offset));
}

print($ciphertext);
