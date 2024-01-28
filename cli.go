package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"todo/add"
	"todo/db"
	list_actionable "todo/list-actionable"
	list_table "todo/list-table"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "todo is a simple cli utility to manage task in progress",
	Args:  cobra.NoArgs,
}

var addCmd = &cobra.Command{
	Use:   `add "my new task"`,
	Short: "register a new task",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if len(args) == 0 {
			p := tea.NewProgram(add.AddInputModel())
			m, err := p.Run()
			if task, ok := m.(add.Model); ok {
				if task.Value == "" {
					return errors.New("Cannot add empty task")
				}

				err = todoDB.CreateTodo(task.Value, "")

				if err != nil {
					return err
				}

				fmt.Printf("\n new task %q created correctly.\n", task.Value)
			}
			return err
		}

		if err != nil {
			return errors.New("Database initialization failed")
		}

		task := args[0]

		if task == "" {
			return errors.New("Cannot add empty task")
		}

		err = todoDB.CreateTodo(task, "")

		if err != nil {
			return err
		}

		fmt.Printf("new task %q created correctly.\n", task)

		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list [command]",
	Short: "list your tasks, it will list only your pending tasks",
	Long: `By default it list your pending tasks, the same as "todo list pending", but if you want to see your completed task you can use "todo list done", or see all your task with "todo list all".
	`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		dateString, err := cmd.Flags().GetString("date")

		if err != nil {
			return errors.New("Not valid date")
		}

		var todos []db.Todo

		if dateString == "" {
			todos, err = todoDB.GetFilteredTasksByState(db.Pending)
		} else {
			if dateString == "today" {
				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Pending, time.Now())
			} else if dateString == "yesterday" {
				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Pending, time.Now().Add(-24*time.Hour))
			} else {
				date, err := time.Parse("2006-01-02", dateString)

				if err != nil {
					return errors.New("Not valid date")
				}

				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Pending, date)
			}
		}

		if err != nil {
			return err
		}

		m := list_table.NewTodoTable(todos)
		p := tea.NewProgram(m)
		_, err = p.Run()

		if err != nil {
			return err
		}

		return nil
	},
}

var listAllCmd = &cobra.Command{
	Use:   "all",
	Short: "list all your tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		dateString, err := cmd.Flags().GetString("date")

		if err != nil {
			return errors.New("Not valid date")
		}

		var todos []db.Todo

		if dateString == "" {
			todos, err = todoDB.GetTasks()
		} else {
			if dateString == "today" {
				todos, err = todoDB.GetFilteredTasksByCreationDate(time.Now())
			} else if dateString == "yesterday" {
				todos, err = todoDB.GetFilteredTasksByCreationDate(time.Now().Add(-24 * time.Hour))
			} else {
				date, err := time.Parse("2006-01-02", dateString)

				if err != nil {
					return errors.New("Not valid date")
				}

				todos, err = todoDB.GetFilteredTasksByCreationDate(date)
			}
		}

		if err != nil {
			return err
		}

		m := list_table.NewTodoTable(todos)
		p := tea.NewProgram(m)
		_, err = p.Run()

		if err != nil {
			return err
		}

		return nil
	},
}

var listDoneTasksCmd = &cobra.Command{
	Use:   "done",
	Short: "list done tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		dateString, err := cmd.Flags().GetString("date")

		if err != nil {
			return errors.New("Not valid date")
		}

		var todos []db.Todo

		if dateString == "" {
			todos, err = todoDB.GetFilteredTasksByState(db.Done)
		} else {
			if dateString == "today" {
				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Done, time.Now())
			} else if dateString == "yesterday" {
				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Done, time.Now().Add(-24*time.Hour))
			} else {
				date, err := time.Parse("2006-01-02", dateString)

				if err != nil {
					return errors.New("Not valid date")
				}

				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Done, date)
			}
		}

		if err != nil {
			return err
		}

		m := list_table.NewTodoTable(todos)
		p := tea.NewProgram(m)
		_, err = p.Run()

		if err != nil {
			return err
		}

		return nil
	},
}

var listPendingTasksCmd = &cobra.Command{
	Use:   "pending",
	Short: "list pending tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		dateString, err := cmd.Flags().GetString("date")

		if err != nil {
			return errors.New("Not valid date")
		}

		var todos []db.Todo

		if dateString == "" {
			todos, err = todoDB.GetFilteredTasksByState(db.Pending)
		} else {
			if dateString == "today" {
				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Pending, time.Now())
			} else if dateString == "yesterday" {
				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Pending, time.Now().Add(-24*time.Hour))
			} else {
				date, err := time.Parse("2006-01-02", dateString)

				if err != nil {
					return errors.New("Not valid date")
				}

				todos, err = todoDB.GetFilteredTasksByStateAndDate(db.Pending, date)
			}
		}

		if err != nil {
			return err
		}

		m := list_table.NewTodoTable(todos)
		p := tea.NewProgram(m)
		_, err = p.Run()

		if err != nil {
			return err
		}

		return nil
	},
}

var markAsDoneCmd = &cobra.Command{
	Use:   "done",
	Short: "mark the task with the id passed as done",
	Long:  `mark the task as done, "todo done 1" will mark the task with the id 1 as done`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		if len(args) == 0 {
			todos, err := todoDB.GetFilteredTasksByState(db.Pending)

			if err != nil {
				return nil
			}

			m := list_actionable.NewTodoTable(todos)
			p := tea.NewProgram(m)
			model, err := p.Run()

			if err != nil {
				return err
			}

			if task, ok := model.(list_actionable.Model); ok {
				id, err := strconv.Atoi(task.SelectedId)

				if err != nil {
					return err
				}

				err = todoDB.CompleteTodo(id)

				if err != nil {
					return errors.New(fmt.Sprintf("todo with id %d couldn't be marked as done", id))
				}

				fmt.Printf("task with the id %d marked as done.\n", id)

				return nil
			}
		}

		id, err := strconv.Atoi(args[0])

		if err != nil {
			return errors.New("Not a valid task id")
		}

		err = todoDB.CompleteTodo(id)

		if err != nil {
			return errors.New(fmt.Sprintf("todo with id %d couldn't be marked as done", id))
		}

		fmt.Printf("task with the id %d marked as done.\n", id)

		return nil
	},
}

var markAsNotDoneCmd = &cobra.Command{
	Use:   "pending",
	Short: "mark the task with the id passed as pending",
	Long:  `mark the task as pending, "todo pending 1" will mark the task with the id 1 as pending`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		if len(args) == 0 {
			todos, err := todoDB.GetFilteredTasksByState(db.Done)

			if err != nil {
				return err
			}

			m := list_actionable.NewTodoTable(todos)
			p := tea.NewProgram(m)
			model, err := p.Run()

			if err != nil {
				return err
			}

			if task, ok := model.(list_actionable.Model); ok {
				id, err := strconv.Atoi(task.SelectedId)

				if err != nil {
					return nil
				}

				err = todoDB.UncompleteTodo(id)

				if err != nil {
					return errors.New(fmt.Sprintf("todo with id %d couldn't be marked as pending", id))
				}

				fmt.Printf("task with the id %d marked as pending.\n", id)

				return nil
			}
		}

		id, err := strconv.Atoi(args[0])

		if err != nil {
			return errors.New("Not a valid task id")
		}

		err = todoDB.UncompleteTodo(id)

		if err != nil {
			return errors.New(fmt.Sprintf("todo with id %d couldn't be marked as pending", id))
		}

		fmt.Printf("task with the id %d marked as pending.\n", id)

		return nil
	},
}

var deleteTodoCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete the task with the id passed",
	Long:  `delete the task, "todo delete 1" will delete the task with the id 1`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		todoDB, err := db.NewTodoDB()
		defer todoDB.Close()

		if err != nil {
			return err
		}

		if len(args) == 0 {
			todos, err := todoDB.GetTasks()

			if err != nil {
				return err
			}

			m := list_actionable.NewTodoTable(todos)
			p := tea.NewProgram(m)
			model, err := p.Run()

			if err != nil {
				return err
			}

			if task, ok := model.(list_actionable.Model); ok {
				id, err := strconv.Atoi(task.SelectedId)

				if err != nil {
					return nil
				}

				err = todoDB.DeleteTodo(id)

				if err != nil {
					return errors.New(fmt.Sprintf("todo with id %d couldn't be deleted", id))
				}

				fmt.Printf("task with the id %d deleted.\n", id)

				return nil
			}
		}

		id, err := strconv.Atoi(args[0])

		if err != nil {
			return errors.New("Not a valid task id")
		}

		err = todoDB.DeleteTodo(id)

		if err != nil {
			return errors.New(fmt.Sprintf("todo with id %d couldn't be deleted", id))
		}

		fmt.Printf("task with the id %d deleted.\n", id)

		return nil
	},
}

// Flag --date -d today, yesterday, 2024-02-01

func init() {
	listCmd.PersistentFlags().StringP(
		"date",
		"d",
		"",
		"date with format YYYY-MM-DD used to filter by the creation date, some special dates are available: today and yesterday",
	)

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(markAsDoneCmd)
	rootCmd.AddCommand(markAsNotDoneCmd)
	rootCmd.AddCommand(deleteTodoCmd)

	listCmd.AddCommand(listAllCmd)
	listCmd.AddCommand(listPendingTasksCmd)
	listCmd.AddCommand(listDoneTasksCmd)
}
