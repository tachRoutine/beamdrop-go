package main

func Help() string {
	return `ekiliBeam - A simple file sharing tool

Usage:
  beam [options]

Options:
  -dir string
		Directory to share files from (default ".")
  -h, --help
		Show this help message and exit	`
}

func PrintHelp() {
	println(Help())
}