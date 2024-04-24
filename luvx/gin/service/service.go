package service

type Runner struct {
    Name    string
    Crontab string
    Fn      func()
}
