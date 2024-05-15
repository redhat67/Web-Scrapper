package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/common-nighthawk/go-figure"
	"github.com/gocolly/colly/v2"
)

func main() {

	asciiArt :=
		`                                                                                                                                        
                  
	
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⣟⣽⢟⣕⢄⠀⠀⠀⠀⠀⠀⣧⠀⡀⠀⠀⠀⠀⢠⢮⡾⣽⡻⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡿⣮⣟⡼⡌⡆⠀⠀⠀⠀⡰⡵⠰⡃⠀⠀⠀⢠⢃⡻⣜⢯⣻⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⡿⣮⡵⣫⡵⢹⠀⠀⠀⠻⡄⡧⣈⠚⠀⠀⠀⢸⢺⡗⣭⣛⢿⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⡿⡵⢞⣫⡆⡹⠀⠀⠀⢑⠥⡨⡐⠫⣢⠀⠀⢸⢺⣽⢚⡭⣯⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣴⣿⣿⣶⣻⡭⣷⢣⠇⠄⠤⢀⠼⣄⣸⡌⡆⢠⡧⠀⠈⡷⣮⢯⣛⣶⣿⣿⣷⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⠀⠀
		⠀⠀⠈⠢⡀⠠⡒⠀⠂⠐⠀⠂⠀⣀⣴⢾⣿⣿⣿⣷⣚⣷⡿⣳⠛⠒⠒⠒⢇⠏⡚⠓⠈⢀⢬⠓⠒⠒⠚⢞⡯⣟⣶⣛⣿⣿⣿⣿⡦⣖⠀⠀⠀⠀⠂⠐⠀⡄⠀⠐⠁⠀⠀⠀
		⠀⠀⠀⠀⠘⡄⠈⠢⡀⢀⣠⡴⣟⢿⣹⢿⣺⣿⣟⣶⣛⡭⠛⠁⠀⠀⠀⠀⠈⣑⣜⡆⢔⣕⣁⠀⠀⠀⠀⠀⠙⢯⡞⣛⣾⣿⣿⣸⣿⣹⣿⠶⣤⡀⠀⡠⠊⢀⠎⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠈⣦⡴⣾⢯⡽⢾⡹⣾⣝⣾⣟⠫⠞⠋⡁⠤⠀⠒⠀⠉⠁⠀⠀⠀⢈⣷⣁⠀⠀⠀⠀⠉⠁⠐⠂⠤⢈⡙⠛⠯⣻⣾⣷⣏⢾⣻⢫⡽⣻⠶⣤⡂⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⣾⣯⣷⣯⣻⣿⣧⣿⣟⡕⢋⢠⠐⠊⠁⠀⡀⣤⢀⣲⣦⣭⣥⣤⣴⣶⣤⣴⣦⣤⣬⣭⣥⣖⣂⣤⣄⡀⠀⠁⠂⡄⡉⢪⣟⣿⣵⣯⣾⣭⣟⣏⣿⡄⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⢸⢿⣿⣿⣿⣯⣿⣿⢳⣾⠇⠋⢀⢠⠰⠘⢃⣡⣾⣿⣿⢿⢹⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢻⣿⣿⣷⣮⣁⠃⠶⡄⡀⠁⠳⢹⣏⣿⣟⣿⣿⣾⣿⣿⣇⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⢸⣟⢟⢿⣝⣟⣻⢻⢿⠁⢀⠜⠋⠀⠀⣠⣿⣿⠏⣿⣧⡇⣸⣿⣟⣯⡿⣽⣻⡿⣽⣿⣏⢈⣼⣿⡟⣧⣝⢧⡀⠀⠁⠃⡄⠀⡏⡿⣻⣻⣻⡽⣫⢻⣿⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⢸⣿⢼⣱⢬⣙⡻⠭⣯⣦⣀⠑⠄⣀⣼⣿⣿⣮⠃⢿⣻⡿⣿⠟⣾⢷⣻⢯⣷⣟⣿⡻⣿⣿⣟⣿⠃⡾⣿⣷⡵⣄⢀⠔⢁⣠⣷⠿⢝⡚⡭⢮⣹⢿⡇⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⢿⣾⣹⢾⡱⣺⠭⣽⣱⢞⡿⣩⢿⣟⣳⣿⢿⣿⣿⣿⡽⡍⢸⣟⣯⣟⣿⢾⡽⣾⡇⠙⣺⣽⣿⣷⣿⣿⣷⢻⣷⣭⢻⣗⢦⣻⡱⣽⢚⣟⣦⣿⣿⠃⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⡸⢿⣿⣟⣧⣷⣹⢥⠷⣞⣱⢏⡾⣱⣎⠧⣽⡱⠯⣟⢿⡡⠈⣿⣳⢿⡾⣿⡽⣷⠃⢄⡧⣟⡝⡧⢿⣲⣙⣧⢳⡞⣧⠺⡭⣖⢿⣸⣏⣷⣿⣿⠏⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⢀⠜⠀⠠⠛⠿⣿⣿⣿⣿⣿⣵⣯⣾⢾⡳⣌⣿⢣⣝⣟⢣⠗⣝⢷⣿⣟⣯⣟⡷⣿⣿⡶⣫⣚⢖⡻⣵⣋⠷⣞⠼⣯⣳⣽⣯⣿⣿⣿⣿⣿⠿⠛⢅⠀⠑⡀⠀⠀⠀⠀
		⠀⠀⢀⠎⠀⡰⠁⠀⠀⠀⠈⠉⠛⣿⣿⣿⣿⣷⣿⣋⣮⢷⡹⣌⣧⢿⢫⢵⣻⣟⣾⣽⣻⢷⡿⡽⣜⢻⢮⣕⢫⢷⣽⣎⣿⣿⣿⣿⣿⡿⡟⠉⠉⠀⠀⠀⠀⢢⠀⠘⠄⠀⠀⠀
		⠀⢀⠊⠀⡐⠀⠀⢀⣀⠀⠀⠀⢸⣿⣾⣿⣿⣿⣿⣿⣾⡷⢻⡹⣜⣣⣿⢫⢧⣿⣟⣾⣽⣿⢻⣞⣧⣏⡳⢯⣛⢾⣶⣿⣿⣿⣿⣿⣿⡇⢻⠀⠀⠀⢀⣀⠀⠀⠡⠀⠈⡄⠀⠀
		⠀⠌⠀⡰⠀⠀⠀⢸⣿⣷⣤⣀⠘⣿⣿⣿⣿⣿⣿⣽⣷⣿⣿⣿⣭⣷⣼⣟⣾⣻⣿⣯⢿⣏⣿⣿⣶⣾⣽⣻⣟⣿⣻⣿⣽⣿⣿⣿⣿⣇⡏⢀⣤⣴⡿⡝⠀⠀⠀⢣⠀⠘⡀	⠀
		⠐⠀⢠⠁⠀⠀⠀⠀⠙⣷⣿⣻⣷⣿⣿⣿⣻⣟⣿⣯⣿⣿⣽⣯⣟⣭⣷⣿⣾⡿⣯⣿⢿⣿⣷⣿⣿⣼⣹⣿⣻⣿⣿⣽⣿⢿⣻⣽⣿⣿⢳⣿⣟⣟⠟⠀⠀⠀⠀⠀⠆⠀⢡⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣷⣟⣿⣽⣿⣯⣷⣿⡾⣿⣶⣯⣟⠿⣟⣿⣻⣽⢿⣽⣻⣿⣻⣾⣳⢯⣟⣿⣻⣿⣿⣟⣿⣭⣛⣿⣳⣯⣿⣏⣾⣟⣾⣿⠀⠀⠀⠀⠀⠀⠘⡀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⢿⣿⣿⣟⣿⢿⣯⣟⣷⣻⣾⣽⣻⣾⡿⣷⣿⣻⣾⢿⣷⣿⢿⣻⣿⣻⣿⢿⣻⣾⣽⣾⣽⣻⣽⡿⣿⢷⣽⣿⣿⣿⢹⠀⠀⠀⠀⠀⠀⠀⠃⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣟⣿⣿⣿⣯⣿⣟⣿⣿⣿⣷⣿⣻⣾⣽⢿⣽⣻⢾⡿⣯⢿⣻⣽⢯⡿⣽⡻⣫⣵⣾⣾⣿⣟⣿⣿⢿⡟⣾⣿⣟⣿⡞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣻⣽⣿⣯⣟⣿⣿⣿⡿⠟⢻⠟⡿⣾⣭⣽⣿⣽⣿⣿⣿⣿⣿⣿⣿⣿⣿⠟⠟⢛⠿⣿⣿⣿⣿⣿⢸⣿⡿⣿⡞⠀⠀⠀⠀⠀⠀⠀⠀⠘⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⢿⣿⣟⣿⣿⣿⣿⠀⠀⢼⠠⠌⠘⣿⣿⣿⣿⡿⣿⣟⢿⣿⣿⣿⡋⠀⠄⡤⠄⠈⣿⣿⣿⣿⣼⢸⣿⣿⢻⠀⠀⠀⠀⠀⠀⠀⠀⠀⢘⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣻⣿⣿⣿⣿⣶⣶⣶⣶⣿⣿⣿⣿⣿⢿⣻⣯⣸⣷⣿⣿⣿⣷⣶⣷⣶⣾⣿⣿⣿⣿⢿⣾⣿⣿⣾⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⢿⣿⣿⢷⣿⣿⣿⣯⡿⢷⣿⣿⣿⣿⣿⣻⣷⣿⣿⣿⣿⣞⣿⣻⣽⢿⣿⣿⣷⣿⣟⣫⣿⣿⣿⣿⣷⡌⣻⠟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣿⣯⣿⣳⣿⣽⣻⣿⣻⡽⣟⣾⣳⣿⣿⣿⣿⣿⣿⣿⣷⢿⣿⣻⢾⣻⣟⡿⣿⡯⢿⣾⣿⣵⣿⡣⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⢿⣿⣿⣿⣿⣿⣿⣿⣯⣷⡿⣿⣿⣿⣿⣿⣿⣿⣿⡾⣟⣯⣯⣴⣿⣿⣿⣿⣿⣿⣿⠎⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⠀⠀⠀
		⠀⠀⠈⠆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡠⠊⣿⣿⣿⣿⣿⣿⣿⡾⣟⣿⣿⣿⣿⣿⣿⣿⣿⣷⣿⣿⣿⣿⣿⣿⣿⣿⣿⡉⠢⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠃⠀⠄⠀
		⠀⠀⠀⠘⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠊⠠⠈⠈⣿⣷⣿⢻⣷⣿⣿⢿⣻⣿⡿⣿⣿⢿⣿⣻⣷⣿⡻⣿⡿⡏⣿⣷⣿⠃⠑⢄⠐⢄⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠎⠀⡰⠀⠀
		⠀⠀⢢⠀⠘⡀⠀⠀⠀⠀⠀⠀⡠⠐⢀⠔⠁⠀⠀⢸⣿⣿⣼⠂⢻⠃⣿⢻⣷⣿⣿⣾⣿⡟⣿⡟⢿⠋⣿⠁⣸⣿⣿⣿⠀⠀⠀⠑⠄⠑⢄⠀⠀⠀⠀⠀⠀⢀⠎⠀⡰⠁⠀⠀
		⠀⠀⠀⠡⡀⠈⢄⠀⢀⡠⠔⠊⠀⠀⠆⠀⠀⠀⠀⢸⣿⣿⣿⣷⣿⠀⣯⠀⡿⢹⠛⣿⠟⣷⣿⠃⣾⡄⠙⣶⣿⣿⣽⣿⠀⠀⠀⠀⠈⡄⠀⠈⠂⠄⣀⠀⢠⠊⠀⡐⠁⠀⠀⠀
		⠀⠀⠀⠀⠑⡀⠀⠪⡉⠉⠉⠉⠉⠉⠁⠀⠀⠀⠀⠘⣿⣿⢿⣿⢧⡚⣿⣤⡇⡸⠀⢿⠀⢹⣿⣦⣿⡇⣿⣿⣿⣿⣾⠇⠀⠀⠀⡀⠈⠉⠉⠉⠉⠉⠉⡙⠁⢀⠜⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠘⢆⠀⠘⢀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣧⢿⣿⣿⣿⣿⣷⣦⣿⣤⣾⣿⣿⣿⣿⡿⡜⣿⣷⡟⠀⠀⣴⣿⣿⣿⣶⡄⠀⠀⣠⠞⠀⣠⠆⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠈⠳⡄⠈⠃⠀⠀⣤⣤⡄⢠⠀⢠⡄⢨⣿⡿⣿⣯⠹⠻⣿⣿⣿⣿⣿⣿⢿⡿⣿⡿⣿⣽⣾⣿⢻⢧⣤⣼⣿⣵⣦⢨⣿⣿⣤⠞⠁⠀⠞⠁⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠈⠢⡀⠀⠐⢌⡉⠐⠢⡀⠀⠀⢰⠹⢿⣻⢿⣸⣄⢹⡻⢿⣿⣿⣿⡿⣿⣜⢿⣞⣿⣿⣳⠏⠉⠀⠀⠙⣿⠟⠀⣹⡇⡇⢀⠔⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠑⢄⡀⠈⠂⠄⡀⠑⠄⠘⠴⠉⣿⣟⣿⣧⣼⣧⣸⡆⣼⢀⣿⣿⣟⣧⡻⣿⣿⣷⣦⡰⠀⡠⠊⡠⡶⠈⣾⣿⠗⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠒⠄⡀⠀⠁⠂⠆⠀⢌⠈⣿⣯⢿⡽⣯⣿⢿⡿⣾⣟⣿⣿⢿⣻⣮⣿⣾⡿⣦⣘⠂⠁⠀⣠⣾⡿⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⡀⢀⡀⠀⠀⠀⠀⠁⠂⡄⢀⡀⠀⠁⠐⠛⠻⢽⣳⣟⣯⣟⣷⣻⣾⣽⠿⢿⣿⣞⣯⢿⡿⣿⣻⣿⢿⣟⢿⡿⠁⣀⠀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⡠⠊⠀⠀⠀⠈⡆⠀⠀⠀⢀⠎⢠⠂⠈⠑⡀⠠⡀⠤⠤⠀⠀⠀⠀⠀⠀⠀⠠⠤⠄⢘⠙⠻⠯⢿⣟⣷⣯⡿⠾⠋⢀⠊⠀⠀⠀⠈⢢⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⢠⠁⠀⢠⠊⠀⠈⠃⠀⠀⢀⠂⡠⠁⠀⠀⠀⠈⢄⠐⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡠⠁⡰⠁⠀⠀⠀⢢⠈⢄⠀⠀⠈⠊⠀⠁⠂⠀⠀⡇⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠈⡄⠀⠘⢄⠀⠀⠀⢀⠠⠁⡐⠁⠀⠀⠀⠀⠀⠈⢂⠈⢄⠀⠀⠀⠀⠀⠀⠀⠀⡐⠁⠔⠀⠀⠀⠀⠀⠀⠡⡀⠢⡀⠀⠀⠀⢀⠄⠀⠀⠇⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠐⢄⠀⠀⠉⠐⠂⠁⠀⠌⠀⠀⠀⠀⠀⠀⠀⠀⠀⠣⠀⠢⠀⠀⠀⠀⠀⢀⠌⢀⠌⠀⠀⠀⠀⠀⠀⠀⠀⠑⢄⠀⠁⠒⠈⠀⠀⢀⠌⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠂⠠⠄⠐⠊⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠱⡀⠡⡀⠀⠀⢠⠊⢀⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠒⠀⠄⠀⠂⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠄⠐⠄⡠⠁⡠⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢆⠈⠀⡐⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
		⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢢⠌⠀⠀⠀⠀⠀
                                                                                                           
	`
	fmt.Println(asciiArt)

	myFigure := figure.NewFigure("Application Name", "", true)
	myFigure.Print()

	var targetURL string
	var scrapeHTML, scrapeLinks, captureScreenshot bool

	flag.StringVar(&targetURL, "url", "", "Web sitesinin URL'si")
	flag.BoolVar(&scrapeHTML, "html", false, "HTML içeriğini çek")
	flag.BoolVar(&scrapeLinks, "links", false, "Linkleri çek")
	flag.BoolVar(&captureScreenshot, "screenshot", false, "Ekran görüntüsü al")
	flag.Parse()

	if targetURL == "" {
		fmt.Println("Hata: URL belirtilmedi.")
		os.Exit(1)
	}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set(
			"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		)
	})

	if scrapeHTML {
		c.OnHTML("html", func(e *colly.HTMLElement) {
			fileName := "Websitesi.html"
			file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			file.WriteString(e.Text)
		})
	}

	if scrapeLinks {
		c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			link := e.Attr("href")
			fileName := "Linkler.txt"
			file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			for _, v := range strings.SplitAfter(link, "http") {
				file.WriteString(fmt.Sprintln(v))
			}
		})
	}

	if captureScreenshot {
		ctx, cancel := context.WithTimeout(context.Background(), 77*time.Second) // 45 saniye bekleme süresi
		defer cancel()

		err := captureScreen(ctx, targetURL)
		if err != nil {
			log.Println("Ekran görüntüsü alınırken bir hata oluştu:", err)
		}
	}

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Hata:", err)
	})

	err := c.Visit(targetURL)
	if err != nil {
		fmt.Println("URL'yi ziyaret ederken bir hata oluştu:", err)
		os.Exit(1)
	}
}

func captureScreen(ctx context.Context, url string) error {
	ctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Navigate(url), chromedp.CaptureScreenshot(&buf)); err != nil {
		return err
	}

	fileName := "yakaladımseni.png"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(buf)
	fmt.Println("Ekran görüntüsü başarıyla kaydedildi:", fileName)
	return nil
}
