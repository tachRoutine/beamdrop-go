package main

func Help() string {
	return `ekiliBeam - A simple file sharing tool

NOTE: YOU NEED TO BE IN THE SAME NETWORK AS THE RECEIVER

Usage:
  beam [options]

Options:
  -dir string
		Directory to share files from (default ".")
  -h, --help	`
}

func PrintHelp() {
	println(Help())
}