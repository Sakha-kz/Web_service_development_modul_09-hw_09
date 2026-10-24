package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Task представляет задачу в системе
type Task struct {
	ID          int
	Title       string
	Description string
	Priority    string // "Высокий", "Средний", "Низкий"
	Status      string // "Новая", "В работе", "Завершена"
}

// TaskStorage определяет интерфейс для управления хранилищем задач
type TaskStorage interface {
	AddTask(title, description, priority string) (*Task, error)
	GetAllTasks() []Task
	GetTaskByID(id int) (*Task, error)
	UpdateTaskStatus(id int, status string) error
	DeleteTask(id int) error
	FindTasks(keyword string) []Task
	GetStatistics() (total, newCount, inProgress, completed int, completionRate float64)
}

// InMemoryTaskStorage реализует интерфейс TaskStorage в оперативной памяти
type InMemoryTaskStorage struct {
	tasks  map[int]Task
	nextID int
}

func NewInMemoryTaskStorage() *InMemoryTaskStorage {
	storage := &InMemoryTaskStorage{
		tasks:  make(map[int]Task),
		nextID: 1,
	}

	// Инициализация демонстрационными задачами
	storage.tasks[1] = Task{ID: 1, Title: "Изучить основы Go", Description: "Синтаксис, структуры и срезы", Priority: "Высокий", Status: "Завершена"}
	storage.tasks[2] = Task{ID: 2, Title: "Написать Task Manager", Description: "Реализовать консольное меню и интерфейс", Priority: "Высокий", Status: "Завершена"}
	storage.tasks[3] = Task{ID: 3, Title: "Изучить указатели в Go", Description: "Передача по значению и ссылке", Priority: "Средний", Status: "В работе"}
	storage.tasks[4] = Task{ID: 4, Title: "Сдать практические работы", Description: "Загрузить решения в Git", Priority: "Высокий", Status: "Новая"}
	storage.tasks[5] = Task{ID: 5, Title: "Прочитать документацию", Description: "Effective Go и стандартная библиотека", Priority: "Низкий", Status: "Новая"}
	storage.nextID = 6

	return storage
}

func (s *InMemoryTaskStorage) AddTask(title, description, priority string) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("название задачи не может быть пустым")
	}

	priority = strings.Title(strings.ToLower(strings.TrimSpace(priority)))
	if priority != "Высокий" && priority != "Средний" && priority != "Низкий" {
		priority = "Средний"
	}

	task := Task{
		ID:          s.nextID,
		Title:       title,
		Description: strings.TrimSpace(description),
		Priority:    priority,
		Status:      "Новая",
	}

	s.tasks[s.nextID] = task
	s.nextID++
	return &task, nil
}

func (s *InMemoryTaskStorage) GetAllTasks() []Task {
	var list []Task
	for _, t := range s.tasks {
		list = append(list, t)
	}
	// Сортировка по ID
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[i].ID > list[j].ID {
				list[i], list[j] = list[j], list[i]
			}
		}
	}
	return list
}

func (s *InMemoryTaskStorage) GetTaskByID(id int) (*Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("задача с ID %d не найдена", id)
	}
	return &task, nil
}

func (s *InMemoryTaskStorage) UpdateTaskStatus(id int, status string) error {
	task, exists := s.tasks[id]
	if !exists {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	status = strings.TrimSpace(status)
	task.Status = status
	s.tasks[id] = task
	return nil
}

func (s *InMemoryTaskStorage) DeleteTask(id int) error {
	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}
	delete(s.tasks, id)
	return nil
}

func (s *InMemoryTaskStorage) FindTasks(keyword string) []Task {
	var matched []Task
	kw := strings.ToLower(strings.TrimSpace(keyword))
	for _, t := range s.tasks {
		if strings.Contains(strings.ToLower(t.Title), kw) || strings.Contains(strings.ToLower(t.Description), kw) {
			matched = append(matched, t)
		}
	}
	return matched
}

func (s *InMemoryTaskStorage) GetStatistics() (total, newCount, inProgress, completed int, completionRate float64) {
	total = len(s.tasks)
	for _, t := range s.tasks {
		switch t.Status {
		case "Новая":
			newCount++
		case "В работе":
			inProgress++
		case "Завершена":
			completed++
		}
	}
	if total > 0 {
		completionRate = (float64(completed) / float64(total)) * 100
	}
	return
}

func main() {
	var storage TaskStorage = NewInMemoryTaskStorage()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n============ TASK MANAGER ============")
		fmt.Println("1. Создать задачу")
		fmt.Println("2. Список всех задач")
		fmt.Println("3. Просмотр задачи по ID")
		fmt.Println("4. Изменить статус задачи")
		fmt.Println("5. Удалить задачу")
		fmt.Println("6. Найти задачу")
		fmt.Println("7. Статистика")
		fmt.Println("0. Выход")
		fmt.Print("Выберите действие: ")

		if !scanner.Scan() {
			break
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Название: ")
			scanner.Scan()
			title := scanner.Text()

			fmt.Print("Описание: ")
			scanner.Scan()
			desc := scanner.Text()

			fmt.Print("Приоритет (Высокий/Средний/Низкий): ")
			scanner.Scan()
			prio := scanner.Text()

			created, err := storage.AddTask(title, desc, prio)
			if err != nil {
				fmt.Printf("Ошибка: %v\n", err)
			} else {
				fmt.Printf("Задача #%d успешно создана!\n", created.ID)
			}

		case "2":
			tasks := storage.GetAllTasks()
			if len(tasks) == 0 {
				fmt.Println("Список задач пуст.")
			} else {
				fmt.Printf("\n%-4s | %-25s | %-12s | %s\n", "ID", "Название", "Приоритет", "Статус")
				fmt.Println(strings.Repeat("-", 60))
				for _, t := range tasks {
					fmt.Printf("#%-3d | %-25s | %-12s | %s\n", t.ID, t.Title, t.Priority, t.Status)
				}
				fmt.Println(strings.Repeat("-", 60))
			}

		case "3":
			fmt.Print("Введите ID задачи: ")
			scanner.Scan()
			id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("Некорректный ID.")
				continue
			}
			t, err := storage.GetTaskByID(id)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\n--- ЗАДАЧА #%d ---\n", t.ID)
				fmt.Printf("Название:    %s\n", t.Title)
				fmt.Printf("Описание:    %s\n", t.Description)
				fmt.Printf("Приоритет:   %s\n", t.Priority)
				fmt.Printf("Статус:      %s\n", t.Status)
			}

		case "4":
			fmt.Print("Введите ID задачи: ")
			scanner.Scan()
			id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("Некорректный ID.")
				continue
			}
			fmt.Println("Доступные статусы: 1. Новая | 2. В работе | 3. Завершена")
			fmt.Print("Выберите новый статус (1-3): ")
			scanner.Scan()
			sChoice := strings.TrimSpace(scanner.Text())
			var statusStr string
			switch sChoice {
			case "1":
				statusStr = "Новая"
			case "2":
				statusStr = "В работе"
			case "3":
				statusStr = "Завершена"
			default:
				statusStr = sChoice
			}

			if err := storage.UpdateTaskStatus(id, statusStr); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Статус задачи #%d успешно изменен на '%s'!\n", id, statusStr)
			}

		case "5":
			fmt.Print("Введите ID задачи для удаления: ")
			scanner.Scan()
			id, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Println("Некорректный ID.")
				continue
			}
			if err := storage.DeleteTask(id); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Задача #%d успешно удалена!\n", id)
			}

		case "6":
			fmt.Print("Введите ключевое слово для поиска: ")
			scanner.Scan()
			kw := scanner.Text()
			results := storage.FindTasks(kw)
			fmt.Printf("Найдено задач: %d\n", len(results))
			for _, t := range results {
				fmt.Printf("- #%d: %s [%s] (%s)\n", t.ID, t.Title, t.Priority, t.Status)
			}

		case "7":
			tot, n, inp, comp, rate := storage.GetStatistics()
			fmt.Println("\n--- СТАТИСТИКА ЗАДАЧ ---")
			fmt.Printf("Всего задач:        %d\n", tot)
			fmt.Printf("Новых:              %d\n", n)
			fmt.Printf("В работе:           %d\n", inp)
			fmt.Printf("Завершённых:        %d\n", comp)
			fmt.Printf("Процент выполнения: %.1f%%\n", rate)

		case "0":
			fmt.Println("Работа программы Task Manager завершена.")
			return

		default:
			fmt.Println("Неверный пункт меню. Выберите от 0 до 7.")
		}
	}
}
