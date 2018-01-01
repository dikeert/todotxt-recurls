# todotxt-recurls
Recurring tasks for Todo.TXT

No CRON/Systemd/Whatever required!


## How it works

 - Install it as [add-on](https://github.com/todotxt/todo.txt-cli/wiki/Creating-and-Installing-Add-ons) replacing `ls` command for your todo.txt-cli.
 - Create `recur.txt` file in your `TODO_DIR` directory.
 
 Create some recurring tasks in your `TODO_DIR/recur.txt` file:
 
 ```
 Weekly Sunday Monday:do a thing @home
 ```
Now, whenever you type `t ls` this tool will examine your `recur.txt` file for tasks to be added into your `todo.txt` file.
If there are such tasks, it'll add them first and then it'll proxy the call to `todo.txt-cli` itself so you _always_ is able to see your recurring tasks.

## Installation
You need to have Go environment installed

Once you have it installed, simply `go get` the tool:

```
$ go get github.com/dikeert/todotxt-recurls
$ # and install it
$ go install github.com/dikeert/todotxt-recurls
```
Now accoding to [Creating and Installing Add ons](https://github.com/todotxt/todo.txt-cli/wiki/Creating-and-Installing-Add-ons):
```
$ mkdir -p ~/.todo.actions.d
$ ln -s $GOPATH/bin/todotxt-recurls ~/.todo.actions.d/ls
$ # or if you have $GOBIN defined like I do
$ ln -s $GOBIN/todotxt-recurls ~/.todo.actions.d/ls
```

## Syntax

Recurrent task syntax as follows:

```
Executor,key=value,key=value Argument Argument:Todo
```
 - __Executor__ is a recurrent task type, right now only Weekly recurring tasks are supported
 - __key=value__ is an Executor attribute. You don't need to edit them, they exist for Executors to be able to store some runtime data. For example, Weekly executor stores attribute `last` with date of the last execution for a task to avoid adding it again and again when you hit `t ls`
 - __Argument__ a list of arguments for an executor. Weekly executor accepts weekdays as arguments - Sunday, Monday, etc...
 
 ## Recurring task types and executors
 
 Currently there are only one recurring task type supported - Weekly.
 
 Each recurring task type has it's own executor. Executor accepts arguments and may or may not update attributes for recurring task. Attributes is the way executors store data between executions.
 
 ### Weekly
 
 This reccuring task type accepts multiple arguments. Those arguments are weekdays: Sunday, Monday, etc.
 
 The executor of this task creates/updates attributes as follows:
  - __last__ - the date of the last execution of a recurring task. Executor updates it each time the task being executed (added to todo.txt)
 
 The executor works as follows:
  - Receives a task
  - Calculates current date
  - Check if task arguments list has current weekday present, if not, skips the task
  - Check if task has __last__ attribute, if yes, checks if before current date, otherwise skips the task
  - Adds the task into `todo.txt` file
  - Creates/updates __last__ atrribute for the task with current date
 

