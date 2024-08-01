package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const QTD = 5
const INTERVAL = 5

func main() {
	var (
		comando int
	)

	for {
		comando = ExibirMenu()

		switch comando {
		case 1:
			IniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo Logs...")
			ImprimirRegistroDeLog()
		case 3:
			fmt.Println("Saindo do programa.")
			os.Exit(0)
		default:
			fmt.Println("Não entendo o comando")
			os.Exit(-1)
		}
	}
}

func ExibirMenu() int {
	var comando int
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("3 - Sair")
	fmt.Scan(&comando)
	return comando
}

func IniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := LerSitesDoArquivo()

	for i := 0; i < QTD; i++ {
		for i, site := range sites {
			fmt.Println("Estou passando na posição", i, ":", site)
			TestarSite(site)
		}
		time.Sleep(INTERVAL * time.Second)
		fmt.Println("")
	}
	
}

func TestarSite(site string) {
	resp, err := http.Get(site)

	if err != nil{
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		RegistarLog(site, true)
	} else {
		fmt.Println("Site:", site, "esta com problema. Status code:", resp.StatusCode)
		RegistarLog(site, false)
	}
}

func LerSitesDoArquivo() []string{
	var sites []string

	arquivo, err := os.Open("sites.txt")
	if err != nil{
		fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(arquivo)
	for{
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF{
			break
		}
	}	

	arquivo.Close()
	return sites
}

func RegistarLog(site string, status bool){
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil{
		fmt.Println(err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func ImprimirRegistroDeLog(){
	arquivo, err := os.ReadFile("log.txt")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(string(arquivo))
}