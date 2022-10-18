# Пример Go concurrency

Данный репозиторий нужен, чтобы продемонстрировать на нем, как реализуется концепт concurrency в go. Темы затронутые в статье:

- С какой целью создавать горутины?
- Как создавать горутины?
- Как следить за их исполнением?

## Структура репозитория

Пакет `bruteforce` состоит из двух файлов:

- `bruteforce.go`
- `bruteforce_test.go`

### bruteforce.go

Код содержит функцию `CombinationsBruteForce`

```go
func CombinationsBruteForce(alphabet string, n int) []string {
	if n <= 0 {
		return nil
	}

	// Copy alphabet into initial product set -- a set of
	// one character sets
	prod := make([]string, len(alphabet))
	for i, char := range alphabet {
		prod[i] = string(char)
	}

	for i := 1; i < n; i++ {
		// The bigger product should be the size of the alphabet times the size of
		// the n-1 size product
		next := make([]string, 0, len(alphabet)*len(prod))

		// Add each char to each word and add it to the new set
		for _, word := range prod {
			for _, char := range alphabet {
				next = append(next, word+string(char))
			}
		}

		prod = next
	}

	return prod
}
```

Его основная задача -- генерировать все возможные строки размера `n`, состоящие из букв в `alphabet`

> Данная функция будет некорректно работать с unicode символами, из-за того, что те занимают больше 1 байта. 
> В этом примере можем этим пренебречь

### main.go

Главная задача этого файла: генерировать случайный `PIN` размера `size`. Затем у `PIN` кода генерируется `SHA-256`[^sha256].
Затем программа пытается найти коллизию исходного `PIN` кода[^collision].  

Программа имитирует поведения "подбора" пароля или `PIN` кода. Этот пример был выбран как относительно интересный[^interesting] 
и действительно долго исполняемый

[^sha256]: https://www.simplilearn.com/tutorials/cyber-security-tutorial/sha-256-algorithm
[^collision]: То есть значение, у которого был бы такой же хэш. В данном примере, это почти (99%+) наверняка -- тот же самый исходный `PIN`
[^interesting]: ну не снова же `time.Sleep` писать в функции для показа примера

## Одна горутина

Классическая реализация в рамках одной `main` горутины (без concurrency):

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"go-concurrency-example/bruteforce"
	"log"
	"math/rand"
	"time"
)

const MaxPinSize = 10
const allowedPinCharacters string = "0123456789"

func randomPinCode(size int) string {
	return randomPinCodeWithRand(size, *rand.New(rand.NewSource(time.Now().UnixNano())))
}

func randomPinCodeWithRand(size int, r rand.Rand) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = allowedPinCharacters[r.Intn(len(allowedPinCharacters))]
	}
	return string(b)
}

func hexSha256(input string) string {
	hashedPin := sha256.Sum256([]byte(input))
	hexHashedPin := hex.EncodeToString(hashedPin[:])
	return hexHashedPin
}

func findCollision(hash string, maxPinSize int) (string, error) {
	for i := 0; i < maxPinSize; i++ {
		log.Printf("Iterating %d-sized pins", i)
		combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, i)
		for _, c := range combinations {
			bfHash := hexSha256(c)
			if bfHash == hash {
				return bfHash, nil
			}
		}
	}
	return "", errors.New("not found")
}

func main() {
	size := 8
	pin := randomPinCode(size)
	hash := hexSha256(pin)
	log.Printf("Calculated hash: %s\n", hash)

	start := time.Now()
	collision, err := findCollision(hash, MaxPinSize)
	if err != nil {
		log.Printf("couldn't find a collision")
	} else {
		log.Printf("found collision! %s produces hash %s\n", collision, hash)
	}
	end := time.Now().Sub(start)
	log.Printf("Finished in %d ns / %d ms / %ds", end, end/time.Millisecond, end/time.Second)
}

```

### Примеры запусков

В данном разделе я привожу примеры с логами запуска с моего компьютера. Если вы запустите этот код на своем компьютере, 
то можете получить другие числа, зависящие от вашего hardware

#### size = 8

```text
2022/10/19 00:15:51 Calculated hash: d06e1e0495d3b8d7c5a935385b5c2a61b3703d2ab07d099ec474acf12b22d8d1
2022/10/19 00:15:51 Iterating 0-sized pins
2022/10/19 00:15:51 Iterating 1-sized pins
2022/10/19 00:15:51 Iterating 2-sized pins
2022/10/19 00:15:51 Iterating 3-sized pins
2022/10/19 00:15:51 Iterating 4-sized pins
2022/10/19 00:15:51 Iterating 5-sized pins
2022/10/19 00:15:51 Iterating 6-sized pins
2022/10/19 00:15:52 Iterating 7-sized pins
2022/10/19 00:15:57 Iterating 8-sized pins
2022/10/19 00:16:06 found collision! d06e1e0495d3b8d7c5a935385b5c2a61b3703d2ab07d099ec474acf12b22d8d1 produces hash d06e1e0495d3b8d7c5a935385b5c2a61b3703d2ab07d099ec474acf12b22d8d1
2022/10/19 00:16:06 Finished in 14624338701 ns / 14624 ms / 14s
```

#### size = 7

```text
2022/10/19 00:19:22 Calculated hash: 6461d31330f821a2a2f1c044156238b8524e49f160d6554934a76006a7a466b6
2022/10/19 00:19:22 Iterating 0-sized pins
2022/10/19 00:19:22 Iterating 1-sized pins
2022/10/19 00:19:22 Iterating 2-sized pins
2022/10/19 00:19:22 Iterating 3-sized pins
2022/10/19 00:19:22 Iterating 4-sized pins
2022/10/19 00:19:22 Iterating 5-sized pins
2022/10/19 00:19:22 Iterating 6-sized pins
2022/10/19 00:19:23 Iterating 7-sized pins
2022/10/19 00:19:24 found collision! 6461d31330f821a2a2f1c044156238b8524e49f160d6554934a76006a7a466b6 produces hash 6461d31330f821a2a2f1c044156238b8524e49f160d6554934a76006a7a466b6
2022/10/19 00:19:24 Finished in 1621668676 ns / 1621 ms / 1s
```

#### size = 6

```text
2022/10/19 00:20:03 Calculated hash: 85aade9f5d82d8f4ad810ba5be6833bc8abf9fb175922186330d9303bd1df632
2022/10/19 00:20:03 Iterating 0-sized pins
2022/10/19 00:20:03 Iterating 1-sized pins
2022/10/19 00:20:03 Iterating 2-sized pins
2022/10/19 00:20:03 Iterating 3-sized pins
2022/10/19 00:20:03 Iterating 4-sized pins
2022/10/19 00:20:03 Iterating 5-sized pins
2022/10/19 00:20:03 Iterating 6-sized pins
2022/10/19 00:20:04 found collision! 85aade9f5d82d8f4ad810ba5be6833bc8abf9fb175922186330d9303bd1df632 produces hash 85aade9f5d82d8f4ad810ba5be6833bc8abf9fb175922186330d9303bd1df632
2022/10/19 00:20:04 Finished in 388911605 ns / 388 ms / 0s
```

---

Итоги:

| Количество символов | Время (мс) |
|:-------------------:|:----------:|
|          8          |   14624    |
|          7          |    1621    |
|          6          |    388     |

Таким образом мы видим, что каждый добавляемый символ очень ощутимо замедляет исполнение программы[^password-chars]

[^password-chars]: Подчеркивает важность размера пароля

## Introducing goroutines

Давайте перепишем нашу функцию так, чтобы использовать горутины. Горутины запускаются с помощью ключевого слова `go`.
Так как сама функция `main` достаточно простая, смысла добавлять туда горутины нет. В таком случае, самая долгая часть 
исполнения нашей программы содержится в коде.

Давайте попробуем

```go
func findCollision(hash string, maxPinSize int) (string, error) {
	for i := 0; i < maxPinSize; i++ {
		log.Printf("Iterating %d-sized pins", i)
		combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, i)
		for _, c := range combinations {
			bfHash := hexSha256(c)
			if bfHash == hash {
				return bfHash, nil
			}
		}
	}
	return "", errors.New("not found")
}
```

Самой долгой частью является перебор на поиск коллизий всех _строк одного размера_
Сначала, для удобства, предлагаю вынести логику для каждой итерации главного цикла в отдельную функцию:

```go
func searchForCollision(hash string, pinSize int) string {
	log.Printf("Iterating %d-sized pins", pinSize)
	combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, pinSize)
	for _, c := range combinations {
		bfHash := hexSha256(c)
		if bfHash == hash {
			return bfHash
		}
	}
	return ""
}

func findCollision(hash string, maxPinSize int) (string, error) {
	for i := 0; i < maxPinSize; i++ {
		collision := searchForCollision(hash, i)
		if collision != "" {
			return collision, nil
		}
	}
	return "", errors.New("not found")
}
```

Теперь, если я попытаюсь запустить отдельную горутину, компилятор выдаст мне ошибку:
```go
func findCollision(hash string, maxPinSize int) (string, error) {
	for i := 0; i < maxPinSize; i++ {
		collision := go searchForCollision(hash, i) //Compilation Error!
		if collision != "" {
			return collision, nil
		}
	}
	return "", errors.New("not found")
}
```

В целом это логично, ведь в этот момент у меня будет 2 горутины, исполняемые _независимо друг от друга_. Более того, вызов
go _неблокирующий_, то есть программа запускает горутину с `searchForCollision`, но основная горутина **не будет ждать результата**
исполнения новой горутины. Таким образом, неудивительно, что go запретил мне приравнивать что-либо по левую сторону от создания
новой горутины

### Каналы

Чтобы все-таки использовать горутины, я перепишу участок, следующим образом:

```go
func searchForCollision(hash string, pinSize int, resultChannel chan string) {
	log.Printf("Iterating %d-sized pins", pinSize)
	combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, pinSize)
	for _, c := range combinations {
		bfHash := hexSha256(c)
		if bfHash == hash {
			resultChannel <- c
		}
	}
}

func findCollision(hash string, maxPinSize int) (string, error) {
	collisions := make(chan string)
	for i := 0; i <= maxPinSize; i++ {
		go searchForCollision(hash, i, collisions)
	}
	select {
	case collision := <-collisions:
		return collision, nil
	}
}
```

Здесь я использую [канал](https://go.dev/tour/concurrency/2) -- структуру, в которую можно положить результат. Я делаю это,
чтобы иметь место, где будет лежать валидный результат. Для этого, я изменил функцию `searchForCollision` --- теперь она принимает
канал и не возвращает результат вообще.

После этого, в основной горутине я ожидаю результата с помощью `select` и реагирую на новое "сообщение" в канале

Что же теперь выведет программа:

```text
2022/10/19 01:18:09 Calculated hash: d09f823f96c0f12ecdc3810c369523f21b72ee5d132f42f887694ef110bf74af
2022/10/19 01:18:09 Iterating 8-sized pins
2022/10/19 01:18:09 Iterating 3-sized pins
2022/10/19 01:18:09 Iterating 1-sized pins
2022/10/19 01:18:09 Iterating 0-sized pins
2022/10/19 01:18:09 Iterating 2-sized pins
2022/10/19 01:18:09 Iterating 4-sized pins
2022/10/19 01:18:09 Iterating 5-sized pins
2022/10/19 01:18:09 Iterating 6-sized pins
2022/10/19 01:18:09 Iterating 7-sized pins
2022/10/19 01:18:45 found collision! 93939966 produces hash d09f823f96c0f12ecdc3810c369523f21b72ee5d132f42f887694ef110bf74af
2022/10/19 01:18:45 Finished in 35738401316 ns / 35738 ms / 35s
```

Здесь есть ряд интересных замечаний:

- Порядок вызовов итераций теперь случаен и непоследователен
  - Он будет разным каждый раз, когда вы будете запускать программу
- Программа теперь исполняется **дольше**?!

Итак, почему порядок вызовов теперь случаен -- относительно понятно[^random-order]
[^random-order]: Вся эта история со скедулингом ОС. К тому же в рантайме Go еще и свой скедулер есть, в общем приоритеты были такими, какими получилось

Но почему же программа исполняется **дольше**? Ну, в первую очередь потому, что ничего не бесплатно -- все наши попытки
сделать так, чтобы наши горутины не "разъехались" в разные стороны, нам нужна синхронизация. Каждый раз, ожидание
"сообщения" в канале, блокирует нас и это занимает много времени. Чтобы программа вела себя предсказуемо, программа
делает много неявных синхронизаций внутри себя (например внутри канала), чтобы не допустить различного рода аномалий, в 
ходе записи над общим участком памяти

### WaitGroup

Можем попробовать снизить количество блокировок и занимаемого времени, сменив канал на массив:

```go
func searchForCollision(hash string, pinSize int, resultChannel []string) {
	log.Printf("Iterating %d-sized pins", pinSize)
	combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, pinSize)
	for _, c := range combinations {
		bfHash := hexSha256(c)
		if bfHash == hash {
			resultChannel[pinSize] = c
		}
	}
	resultChannel[pinSize] = ""
}

func findCollision(hash string, maxPinSize int) (string, error) {
	var collisions []string = make([]string, maxPinSize+1)
	for i := 0; i <= maxPinSize; i++ {
		go searchForCollision(hash, i, collisions)
	}
	for i := 0; i <= maxPinSize; i++ {
		if collisions[i] != "" {
			return collisions[i], nil
		}
	}
	return "", errors.New("not found")
}
```

Однако, запустив эту программу, мы видим странную вещь:

```text
2022/10/19 01:48:42 Calculated hash: 4cde341b873e788ee0ca794fe82a69131858d5477d3ebc6d621f55db6f9f1997
2022/10/19 01:48:42 couldn't find a collision
2022/10/19 01:48:42 Finished in 30116 ns / 0 ms / 0s
```

Теперь программа работает _быстро, но неправильно_. Дело в том, что основная горутина не ждет завершения исполнения остальных
горутин, потому сразу пытается найти в массиве результат, но вторая горутина **не успела** записать в массив результат, потому 
мы видим такую аномалию

Можем исправить это поведение, используя [WaitGroup](https://gobyexample.com/waitgroups).

Он позволяет дождаться исполнения горутины до конца, прежде чем что-то делать дальше с помощью `Wait` функции

```go
func searchForCollision(hash string, pinSize int, resultChannel []string, wg *sync.WaitGroup) {
	log.Printf("Iterating %d-sized pins", pinSize)
	combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, pinSize)
	for _, c := range combinations {
		bfHash := hexSha256(c)
		if bfHash == hash {
			resultChannel[pinSize] = c
			wg.Done()
		}
	}
	resultChannel[pinSize] = ""
	wg.Done()
}

func findCollision(hash string, maxPinSize int) (string, error) {
	var collisions []string = make([]string, maxPinSize+1)
	var wg sync.WaitGroup
	for i := 0; i <= maxPinSize; i++ {
		wg.Add(1)
		go searchForCollision(hash, i, collisions, &wg)
	}
	wg.Wait()
	for i := 0; i <= maxPinSize; i++ {
		if collisions[i] != "" {
			return collisions[i], nil
		}
	}
	return "", errors.New("not found")
}
```

```text
2022/10/19 02:39:14 Calculated hash: 98e345f293d2ee4d3103961151535182f4106b42d88ec57a1c632ed6e719cb22
2022/10/19 02:39:14 Iterating 8-sized pins
2022/10/19 02:39:14 Iterating 0-sized pins
2022/10/19 02:39:14 Iterating 6-sized pins
2022/10/19 02:39:14 Iterating 4-sized pins
2022/10/19 02:39:14 Iterating 7-sized pins
2022/10/19 02:39:14 Iterating 5-sized pins
2022/10/19 02:39:14 Iterating 2-sized pins
2022/10/19 02:39:14 Iterating 1-sized pins
2022/10/19 02:39:14 Iterating 3-sized pins
2022/10/19 02:39:48 found collision! 87273267 produces hash 98e345f293d2ee4d3103961151535182f4106b42d88ec57a1c632ed6e719cb22
2022/10/19 02:39:48 Finished in 33738902029 ns / 33738 ms / 33s
```

### Все таки каналы

Тем не менее, использование каналов -- более идиоматический способ работы с ожиданием значения

```go
func searchForCollision(hash string, pinSize int, collisionChan chan string) {
	log.Printf("Iterating %d-sized pins", pinSize)
	combinations := bruteforce.CombinationsBruteForce(allowedPinCharacters, pinSize)
	processPart := func(ccs []string, cc chan string) {
		for _, comb := range combinations {
			bfHash := hexSha256(comb)
			//log.Printf("computing hash for %s:%s", ccs, bfHash)
			if bfHash == hash {
				cc <- comb
			}
		}
	}
	//Есть соблазн создать горутину не для половины списка, а для каждого элемента
	//Тем не менее, это вызовет лишь замедление работы в несколько раз --- слишком большие расходы на синхронизацию
	go processPart(combinations[0:len(combinations)/2], collisionChan)
	go processPart(combinations[len(combinations)/2:], collisionChan)
}

func findCollision(hash string, maxPinSize int) string {
	var collisionChan = make(chan string)
	for i := 0; i <= maxPinSize; i++ {
		go searchForCollision(hash, i, collisionChan)
	}
	select {
	case c := <-collisionChan:
		return c
	}
}
```

Несколько запусков для показа разброса значений:

```text
2022/10/19 02:41:05 Calculated hash: 969f0326cbb86a59e4534950398ea60d2a76a33d3e81516888e7424708d4ba26
2022/10/19 02:41:05 Iterating 8-sized pins
2022/10/19 02:41:05 Iterating 3-sized pins
2022/10/19 02:41:05 Iterating 5-sized pins
2022/10/19 02:41:05 Iterating 1-sized pins
2022/10/19 02:41:05 Iterating 4-sized pins
2022/10/19 02:41:05 Iterating 2-sized pins
2022/10/19 02:41:05 Iterating 6-sized pins
2022/10/19 02:41:05 Iterating 7-sized pins
2022/10/19 02:41:05 Iterating 0-sized pins
2022/10/19 02:41:17 found collision! 17714601 produces hash 969f0326cbb86a59e4534950398ea60d2a76a33d3e81516888e7424708d4ba26
2022/10/19 02:41:17 Finished in 12734207710 ns / 12734 ms / 12s
```

```text
2022/10/19 02:47:18 Calculated hash: a442263ab77e96ac7492112007489a045f912e2b0ed3a0d49c96e619564216f3
2022/10/19 02:47:18 Iterating 8-sized pins
2022/10/19 02:47:18 Iterating 0-sized pins
2022/10/19 02:47:18 Iterating 2-sized pins
2022/10/19 02:47:18 Iterating 1-sized pins
2022/10/19 02:47:18 Iterating 5-sized pins
2022/10/19 02:47:18 Iterating 4-sized pins
2022/10/19 02:47:18 Iterating 6-sized pins
2022/10/19 02:47:18 Iterating 3-sized pins
2022/10/19 02:47:18 Iterating 7-sized pins
2022/10/19 02:47:43 found collision! 49656660 produces hash a442263ab77e96ac7492112007489a045f912e2b0ed3a0d49c96e619564216f3
2022/10/19 02:47:43 Finished in 24865791391 ns / 24865 ms / 24s
```

```text
2022/10/19 02:48:08 Calculated hash: cb198061f1c41b37cefa1706aa053811359794b03303110de8e504620ebb7c9a
2022/10/19 02:48:08 Iterating 8-sized pins
2022/10/19 02:48:08 Iterating 3-sized pins
2022/10/19 02:48:08 Iterating 1-sized pins
2022/10/19 02:48:08 Iterating 0-sized pins
2022/10/19 02:48:08 Iterating 5-sized pins
2022/10/19 02:48:08 Iterating 6-sized pins
2022/10/19 02:48:08 Iterating 4-sized pins
2022/10/19 02:48:08 Iterating 7-sized pins
2022/10/19 02:48:08 Iterating 2-sized pins
2022/10/19 02:48:19 found collision! 10804177 produces hash cb198061f1c41b37cefa1706aa053811359794b03303110de8e504620ebb7c9a
2022/10/19 02:48:19 Finished in 11554575189 ns / 11554 ms / 11s
```

```text
2022/10/19 02:49:03 Calculated hash: 0b205be5d05fcbb83ded3c634c683ee08967b6a65a343f129318e4e786a9011a
2022/10/19 02:49:03 Iterating 8-sized pins
2022/10/19 02:49:03 Iterating 3-sized pins
2022/10/19 02:49:03 Iterating 1-sized pins
2022/10/19 02:49:03 Iterating 0-sized pins
2022/10/19 02:49:03 Iterating 2-sized pins
2022/10/19 02:49:03 Iterating 5-sized pins
2022/10/19 02:49:03 Iterating 6-sized pins
2022/10/19 02:49:03 Iterating 7-sized pins
2022/10/19 02:49:03 Iterating 4-sized pins
2022/10/19 02:49:43 found collision! 81307591 produces hash 0b205be5d05fcbb83ded3c634c683ee08967b6a65a343f129318e4e786a9011a
2022/10/19 02:49:43 Finished in 39849001629 ns / 39849 ms / 39s

```