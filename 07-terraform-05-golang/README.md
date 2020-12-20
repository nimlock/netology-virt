# Домашнее задание к занятию "7.5. Основы golang"

## Модуль 7. Облачная инфраструктура. Terraform

### Студент: Иван Жиляев

>С `golang` в рамках курса, мы будем работать не много, поэтому можно использовать любой IDE. 
>Но рекомендуем ознакомиться с [GoLand](https://www.jetbrains.com/ru-ru/go/).  

## Задача 1. Установите golang.
>1. Воспользуйтесь инструкций с официального сайта: [https://golang.org/](https://golang.org/).
>2. Так же для тестирования кода можно использовать песочницу: [https://play.golang.org/](https://play.golang.org/).

Вопросов с установкой не возникло.

## Задача 2. Знакомство с gotour.
>У Golang есть обучающая интерактивная консоль [https://tour.golang.org/](https://tour.golang.org/). 
>Рекомендуется изучить максимальное количество примеров. В консоли уже написан необходимый код, 
>осталось только с ним ознакомиться и поэкспериментировать как написано в инструкции в левой части экрана.  

С gotour ознакомился.

## Задача 3. Написание кода. 
>Цель этого задания закрепить знания о базовом синтаксисе языка. Можно использовать редактор кода 
>на своем компьютере, либо использовать песочницу: [https://play.golang.org/](https://play.golang.org/).
>
>1. Напишите программу для перевода метров в футы (1 фут = 0.3048 метр). Можно запросить исходные данные 
>у пользователя, а можно статически задать в коде.
>    Для взаимодействия с пользователем можно использовать функцию `Scanf`:
>    ```
>    package main
>    
>    import "fmt"
>    
>    func main() {
>        fmt.Print("Enter a number: ")
>        var input float64
>        fmt.Scanf("%f", &input)
>    
>        output := input * 2
>    
>        fmt.Println(output)    
>    }
>    ```
> 
>1. Напишите программу, которая найдет наименьший элемент в любом заданном списке, например:
>    ```
>    x := []int{48,96,86,68,57,82,63,70,37,34,83,27,19,97,9,17,}
>    ```
>1. Напишите программу, которая выводит числа от 1 до 100, которые делятся на 3. То есть `(3, 6, 9, …)`.
>
>В виде решения ссылку на код или сам код. 

1. Программа находится в файле [task3-1.go](task3-1/task3-1.go), вот её код:

    ```
    package main

    import (
      "fmt"
      "math"
      "os"
    )

    func main() {
      fmt.Println("Enter a number of meters to convert into ft: ")
      var input float64
      fmt.Fscan(os.Stdin, &input)

      fmt.Printf("Ok, it equal to %.4f ft\n", MetersToFt(input))
    }

    func MetersToFt(meters float64) float64 {
      coeff := math.Pow(0.3048, -1)
      result := math.Round(meters*coeff*10000) / 10000
      return result
    }
    ```

1. Программа находится в файле [task3-2.go](task3-2/task3-2.go), её код:

    ```
    package main

    import (
      "fmt"
    )

    const (
      UintSize = 32 << (^uint(0) >> 32 & 1)
      MaxInt   = 1<<(UintSize-1) - 1 // Определяем наибольшее значение в int
    )

    func main() {
      x := []int{48, 96, 86, 68, 57, 82, 63, 70, 37, 34, 83, 27, 19, 97, 9, 17, 555}
      fmt.Printf("Max value in array/split is %d\n", FindMinInList(x))
    }

    func FindMinInList(list []int) int {
      minValue := MaxInt
      for _, value := range list {
        if value < minValue {
          minValue = value
        }
      }
      return minValue
    }
    ```

1. Программа находится в файле [task3-3.go](task3-3/task3-3.go), её код:

    ```
    package main

    import (
      "fmt"
    )

    func main() {
      start, end := 0, 100
      fmt.Printf("Let`s find all numbers from %d and %d that clear division by 3\n", start, end)
      fmt.Println(ClearDivisionBy3(start, end))
    }

    func ClearDivisionBy3(start, end int) []int {
      result := make([]int, 0, 0)
      for i := start; i < end; i++ {
        if i%3 == 0 {
          result = append(result, i)
        }
      }
      return result
    }
    ```

## Задача 4. Протестировать код (не обязательно).

>Создайте тесты для функций из предыдущего задания. 

1. Тест для функции в файле [task3-1_test.go](task3-1/task3-1_test.go), код тестов:

    ```
    package main

    import "testing"

    func TestMetersToFt(t *testing.T) {
      var v float64
      v = MetersToFt(10)
      if v != 32.8084 {
        t.Error("Expected 32.8084, got ", v)
      }
    }
    ```

    Выполнение тестов:

    ```
    ivan@kubang:~/study/netology-virt/07-terraform-05-golang/task3-1$ go test
    PASS
    ok      _/home/ivan/study/netology-virt/07-terraform-05-golang/task3-1  0.002s
    ```

1. Тест для функции в файле [task3-2_test.go](task3-2/task3-2_test.go), код тестов:

    ```
    package main

    import "testing"

    func TestFindMinInList(t *testing.T) {
      list := []int{1563, 212, 15, 88, 657, 387, 692, 1294}
      var v int
      v = FindMinInList(list)
      if v != 15 {
        t.Error("Expected 15, got ", v)
      }
    }
    ```

    Выполнение тестов:

    ```
    ivan@kubang:~/study/netology-virt/07-terraform-05-golang/task3-2$ go test
    PASS
    ok      _/home/ivan/study/netology-virt/07-terraform-05-golang/task3-2  0.003s
    ```

1. Тест для функции в файле [task3-3_test.go](task3-3/task3-3_test.go), код тестов:

    ```
    package main

    import (
      "reflect"
      "testing"
    )

    func TestClearDivisionBy3(t *testing.T) {
      // var v float64
      v := ClearDivisionBy3(10, 20)
      rightAnswer := []int{12, 15, 18}
      if !(reflect.DeepEqual(v, rightAnswer)) {
        t.Error("Expected ", rightAnswer, " got ", v)
      }
    }
    ```

    Выполнение тестов:

    ```
    ivan@kubang:~/study/netology-virt/07-terraform-05-golang/task3-3$ go test
    PASS
    ok      _/home/ivan/study/netology-virt/07-terraform-05-golang/task3-3  0.002s
    ```
