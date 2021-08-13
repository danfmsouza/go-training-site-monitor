package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5

func main() {
	exibirIntroducao()
	for {
		exibirMenu()
		comando := lerComando()
		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			exibirLogs()
		case 0:
			fmt.Println("Saindo do programa: ")
			os.Exit(0)
		default:
			fmt.Println("Instrução inválida!, escolha 0, 1 ou 2!")
			os.Exit(-1)
		}
	}
}

func exibirIntroducao() {
	nome := "Douglas"
	versao := 1.2
	fmt.Println("Olá iMundo! Tudo bem " + nome + "?")
	fmt.Printf("Versão do programa: %.1f \n", versao)
}

func exibirMenu() {
	fmt.Println("")
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair")
}

func lerComando() int {
	comandoLido := 0
	fmt.Scan(&comandoLido)

	return comandoLido
}

func iniciarMonitoramento() {
	sites := lerSitesNoArquivo()

	for i := 0; i <= monitoramentos; i++ {
		for i, site := range sites {
			fmt.Println("Testar o site", i+1, ":", site)
			testaSite(site)
		}
		time.Sleep(delay * time.Second)
	}
}

func testaSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("Opa deu ruim", err)
		arquivofalhas, _ := os.OpenFile("sites.err", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		arquivofalhas.WriteString(time.Now().Format("2006-01-02 15:04:05 MST") + ": " + site + " -error: " + err.Error() + "\n")
		arquivofalhas.Close()
	} else if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado corretamente!")
		registrarLogs(site, true)
	} else if resp.StatusCode != 200 {
		fmt.Println("Site:", site, "não carregado corretamente, código!")
		registrarLogs(site, false)
	} else {
		fmt.Println("Opa deu ruim", err)
	}

}

func lerSitesNoArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.conf")

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF {
			break
		}
	}
	arquivo.Close()
	return sites
}

func registrarLogs(site string, status bool) {
	arquivo, err := os.OpenFile("sites.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Problema ao abrir o arquivo de log: ", err)
	}

	arquivo.WriteString(time.Now().Format("2006-01-02 15:04:05 MST") + " : " + site + " -online: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func exibirLogs() {
	fmt.Println("")
	fmt.Println("1- Exibir Logs de Status: ")
	fmt.Println("2- Exibir Logs de Erro: ")
	fmt.Println("0- Retornar ao Menu Inicial: ")
	comandoLog := lerComando()
	switch comandoLog {
	case 1:
		imprimeLogs(1)
	case 2:
		imprimeLogs(2)
	case 0:
		break
	default:
		fmt.Println("Instrução inválida!, escolha 0, 1 ou 2!")
		break
	}
}

func imprimeLogs(log int) {
	if log == 1 {
		arquivolog, err := ioutil.ReadFile("sites.log")
		if err != nil {
			fmt.Println("Problema ao abrir o arquivo de log: ", err)
		}
		fmt.Println(string(arquivolog))
	} else if log == 2 {
		arquivoerr, err := ioutil.ReadFile("sites.err")
		if err != nil {
			fmt.Println("Problema ao abrir o arquivo de log: ", err)
		}
		fmt.Println(string(arquivoerr))
	} else {
		fmt.Println("Problema ao abrir o arquivo de log")
	}
	return
}
