package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	testData := &[]string{"слиток", "тяпка", "пятак", "листок", "пятка", "столик"}
	result := FindAnagramm(testData)
	for k, v := range *result {
		fmt.Printf("Key: %s Set: %v\n", k, *v)
	}
}

/*
Проходимся по словарю, каждое слово преобразуем в нижний регистр. Если слово нельзя получитть по ключу, возможно
два вариант:
1) ключ-анаграмма уже создан, и нужно проверить, что ключ действительно это анаграмма
2) ни самого слова ни его анаграммы среди ключей нет, и мы создаём новый ключ.
Если первый вариант после прохода по мапе выполнился, то flag станет true. И второй вариант не выполнится.
И наоборот та же логика.
*/
func FindAnagramm(dict *[]string) *map[string]*[]string {
	result := make(map[string]*[]string, 0)

	for _, word := range *dict {
		word = strings.ToLower(word)
		flag := false
		if _, ok := result[word]; !ok {
			for k := range result {
				if thisIsAnagramm(word, k) {
					*result[k] = append(*result[k], word)
					flag = true
					break
				}
			}
			if !flag {
				result[word] = &[]string{word}
			}

		}
	}

	// удлаяем множества с одним элементом и сортируем остальные множетсва
	for k := range result {
		if len(*result[k]) == 1 {
			delete(result, k)
			continue
		}
		sort.Slice(*result[k], func(i, j int) bool {
			return (*result[k])[i] < (*result[k])[j]
		})

	}

	return &result
}

// Функция для проверки двух слов на анаграмму.
// Сортируем оба слова и сравниваем.
func thisIsAnagramm(word1, word2 string) bool {
	if len(word1) != len(word2) {
		return false
	}

	sortedWord1 := []rune(word1)
	sortedWord2 := []rune(word2)
	sort.Slice(sortedWord1, func(i, j int) bool {
		return sortedWord1[i] < sortedWord1[j]
	})
	sort.Slice(sortedWord2, func(i, j int) bool {
		return sortedWord2[i] < sortedWord2[j]
	})

	return string(sortedWord1) == string(sortedWord2)
}
