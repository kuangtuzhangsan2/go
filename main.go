package main

import "fmt"

type stu struct {
	id    int
	name  string
	class string
}

type s1 struct {
	alls []*stu
}


func news(id int, name, class string) *stu {
	return &stu{
		id:    id,
		name:  name,
		class: class,
	}
}


func r1() *s1 {
	return &s1{
		alls: make([]*stu, 0, 100),
	}
}


func (s *s1) add(new *stu) {
	s.alls = append(s.alls, new)
}


func (s *s1) mod(new *stu) {
	for i, v := range s.alls { 
		if new.id == v.id {
			s.alls[i] = new
			return
		}
	}
	fmt.Println("未找到该学号学生")
}

// 显示学生（修复 printf 格式）
func (s *s1) show() {
	for _, i := range s.alls {
		fmt.Printf("学号：%d 姓名：%s 班级：%s\n", i.id, i.name, i.class)
	}
}


func get() *stu {
	var (
		id    int
		name  string
		class string
	)
	fmt.Print("请输入学号：")
	fmt.Scan(&id)
	fmt.Print("请输入姓名：")
	fmt.Scan(&name)
	fmt.Print("请输入班级：")
	fmt.Scan(&class)

	stud := stu{
		id:    id,
		name:  name,
		class: class,
	}
	return &stud
}

func main() {

	sm := r1()

	for {
		fmt.Println("=====菜单=====")
		fmt.Println("1. 添加学生")
		fmt.Println("2. 修改学生")
		fmt.Println("3. 显示所有学生")
		fmt.Println("4. 退出")
		fmt.Print("请输入选项：")

		var input int
		fmt.Scan(&input)

		switch input {
		case 1:
			stud := get()
			sm.add(stud)
		case 2:
			stud := get()
			sm.mod(stud)
		case 3:
			sm.show()
		case 4:
			return 
		}
	}
}