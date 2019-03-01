package server

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Model struct {
	Dsn string
	Db  *sql.DB
}

func NewModel() *Model {
	model := Model{
		Dsn: "root:root@/go_task",
	}
	var err error
	model.Db, err = sql.Open("mysql", model.Dsn)
	if err != nil {
		panic(err)
		return nil
	}
	return &model
}

type Device struct {
	Id      int64
	Uuid    string
	Token   string
	Updated int64
	Created int64
}

func (m *Model) CheckKeyExists(key string) error {
	_, err := m.getDeviceInfo("SELECT * FROM device WHERE uuid = '"+key+"'", Device{})
	return err
}

func (m *Model) GetDeviceInfoByKey(key string) (*Device, error) {
	return m.getDeviceInfo("SELECT * FROM device WHERE uuid = '"+key+"'", Device{})
}

func (m *Model) getDeviceInfo(sql string, d Device) (*Device, error) {
	row := m.Db.QueryRow(sql)
	err := row.Scan(&d.Id, &d.Uuid, &d.Token, &d.Updated, &d.Created)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (m *Model) addDeviceInfo(sql string, d Device) {
	_, err := m.Db.Exec(sql, d.Uuid, d.Token, d.Updated, d.Created)
	if err != nil {
		fmt.Println("insert failed,", err)
	}
}

func (m *Model) AddDeviceInfo(token string, uuid string) {
	now := time.Now().Unix()
	d := Device{Uuid: uuid, Token: token, Created: now, Updated: now}
	m.addDeviceInfo("INSERT INTO device (uuid,token,created,updated) VALUES (?,?,?,?)", d)
}

type Task struct {
	Uuid     string
	Title    string
	Category string
	Body     string
	Url      string
}

func (m *Model) addTaskInfo(sql string, task Task) {
	_, err := m.Db.Exec(sql, task.Uuid, task.Title, task.Category, task.Body, task.Url)
	if err != nil {
		fmt.Println("insert failed,", err)
	}
}

func (m *Model) AddTaskInfo(task Task) {
	m.addTaskInfo("INSERT INTO task (uuid,title,category,body) VALUES (?,?,?,?)", task)
}

func (m *Model) GetTaskInfoByKey(key string) (*Task, error) {
	return m.getTaskInfoByKey("SELECT * FROM task WHERE uuid = '"+key+"'", Task{})
}

func (m *Model) getTaskInfoByKey(sql string, task Task) (*Task, error) {
	var id int64
	row := m.Db.QueryRow(sql)
	err := row.Scan(&id, &task.Uuid, &task.Title, &task.Category, &task.Body, &task.Url)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
