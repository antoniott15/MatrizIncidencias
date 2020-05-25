# Laboratorio 6.2 Base de datos


## Seccion 1
Para esta seccion es necesario generar el mecanismo que filtre los datos de las stopwords y generar un archivo que guarde estos.

### 1.a Filtrar por stop words

Parte de codigo que genera un mapa con las keys correspondientes a cada stopword.
```golang
func scanStopLisT(path string) (map[string]int, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)
	words := make(map[string]int)
	for scanner.Scan() {
		words[scanner.Text()] = 1
	}

	return words, nil
}
```


### 1.b El resultado de la lista de incidencias se encuentra en el archivo result.

```golang
func scanWords(path string, stopWords map[string]int, words map[string][]string) (map[string][]string, error) {
	if path == PATHSTOPLIST {
		return nil, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if _, ok := stopWords[scanner.Text()]; !ok {
			words[scanner.Text()] = append(words[scanner.Text()], path)
		}
	}

	return words, nil
}
```

## Seccion 2

Para esta seccion se implemento un query lenguage que transforma de de una consulta normal a una booleana.


Esta seccion permite discernir en and, or y not y generará una respuesta por si esta que esta o no la palabra en los libros especificados.
```golang
func QueryMachine(query string, toQuery map[string][]string) bool {
	query = strings.ToLower(query)
	str := strings.Split(query, " ")

	var ok, ok2 bool
	if str[0] == "not" {
		_, ok = toQuery[str[1]]
		_, ok2 = toQuery[str[len(str)-1]]
	} else {
		_, ok = toQuery[str[0]]
		_, ok2 = toQuery[str[len(str)-1]]
	}

	if !strings.Contains(query, "not") {
		if strings.Contains(query, "or") {
			return ok || ok2
		} else if strings.Contains(query, "and") {
			return ok && ok2
		}
	} else {
		if str[0] == "not" {
			if strings.Contains(query, "or") {
				return !ok || ok2
			} else if strings.Contains(query, "and") {
				return !ok && ok2
			}
		} else {
			if strings.Contains(query, "or") {
				return ok || !ok2
			} else if strings.Contains(query, "and") {
				return ok && !ok2
			}
		}
	}

	return false
}
```


###  Para las consultas

Para las consultas se tomo encuenta 3 consultas basicas, en el archivo usted podría generar más de ser necesario


```golang
	query1 := "huye or asd"
	fmt.Print(QueryMachine(query1, words))

	query2 := "jardinero and llegan"
	fmt.Print(QueryMachine(query2, words))

	query3 := "mithril and not funeral"
	fmt.Print(QueryMachine(query3, words))
```


###  Para usar
Para poder visualizar el programa podría ejecutar haciendo:

    $ go build

y ejecutar el binario.