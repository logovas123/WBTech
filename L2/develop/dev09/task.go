package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// savePage скачивает содержимое страницы и сохраняет её локально.
func savePage(body []byte, pageURL, saveDir string) error {
	parseURL, err := url.Parse(pageURL)
	if err != nil {
		return fmt.Errorf("error parse url: %v", err)
	}

	// проверка главной страницы
	var path string
	if parseURL.Path == "/" || parseURL.Path == "" {
		path = filepath.Join(saveDir, "index.html")
	} else {
		dir := filepath.Join(saveDir, filepath.Dir(parseURL.Path)) // возвращаем последний сегмент
		fileName := filepath.Base(parseURL.Path)
		if !strings.HasSuffix(fileName, ".html") {
			fileName += ".html" // добавляем расширение .html
		}
		path = filepath.Join(dir, fileName)
	}

	// создаём директории и файл для сохранения
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return fmt.Errorf("error create dir: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error createe file %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		return fmt.Errorf("error write in file: %w", err)
	}

	fmt.Printf("Success save to %s\n", path)
	return nil
}

// extractLinks извлекает ссылки из HTML
func extractLinks(baseURL string, body io.Reader) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(body)
	// проходимся по документу
	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return links, nil
			}
			return nil, tokenizer.Err()

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val
						// игнорируем якорные ссылки
						if strings.HasPrefix(link, "#") {
							break
						}
						// преобразуем относительные URL в абсолютные.
						absoluteLink, err := changeURL(baseURL, link)
						if err != nil {
							fmt.Println("error change url: ", err)
							break
						}
						links = append(links, absoluteLink)
					}
				}
			}
		default:
			continue
		}
	}
}

// changeURL преобразует относительный код в абсолютный
func changeURL(base, relative string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	relativeURL, err := url.Parse(relative)
	if err != nil {
		return "", err
	}
	resolved := baseURL.ResolveReference(relativeURL) // метод для преобразования
	return resolved.String(), nil
}

// wget скачивает сайт рекурсивно
func wget(client *http.Client, baseURL, saveDir string, listURL map[string]bool) error {
	// если url true, то он уже скачан
	if listURL[baseURL] {
		return nil
	}

	listURL[baseURL] = true

	fmt.Printf("downloading: %s\n", baseURL)

	resp, err := client.Get(baseURL)
	if err != nil {
		return fmt.Errorf("error get request to %v: %v", baseURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read response body: %v", err)
	}

	// сохраняем страницу
	if err := savePage(body, baseURL, saveDir); err != nil {
		return fmt.Errorf("error save page: %v", err)
	}

	// извлекаем ссылки со страницы
	links, err := extractLinks(baseURL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	// фильтруем и рекурсивно скачиваем ссылки, принадлежащие тому же хосту
	base, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	for _, link := range links {
		parsedLink, err := url.Parse(link)
		if err != nil {
			continue
		}
		if parsedLink.Host == base.Host {
			if err := wget(client, link, saveDir, listURL); err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	// проверка аргументов.
	if len(os.Args) < 2 {
		fmt.Println("Error args < 2")
		return
	}

	baseURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Println("Error parse url: ", err)
		os.Exit(1)
	}

	// если нет названия директории для сохранения, то папка для сохрания будет иметь название хоста
	var saveDir string
	if len(os.Args) > 2 {
		saveDir = os.Args[2]
	} else {
		saveDir = baseURL.Host
	}

	client := &http.Client{}
	// избегаем повторного скачивания одного и того же сайта (так как скачивание рекурсивное)
	listURL := make(map[string]bool)

	err = wget(client, baseURL.String(), saveDir, listURL)
	if err != nil {
		fmt.Println("Error wget: ", err)
		os.Exit(1)
	}

	fmt.Println("site download success")
}
